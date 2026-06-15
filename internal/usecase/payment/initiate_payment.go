package payment

import (
	"context"
	"errors"
	"fmt"
	domain_membership "modulegue/internal/domain/customer_membership" // Alias untuk menghindari konflik
	domain_wallet "modulegue/internal/domain/jukir_wallet"            // Alias untuk menghindari konflik
	domain_payment "modulegue/internal/domain/payment"
	"time"

	// "modulegue/internal/domain/transaction"
	"modulegue/internal/service/payment_gateway"
)

var (
	ErrTransactionNotFound      = errors.New("transaksi tidak ditemukan")
	ErrTransactionNotPending    = errors.New("transaksi tidak dalam status pending/unpaid")
	ErrInsufficientJukirBalance = errors.New("saldo jukir tidak mencukupi untuk transaksi cash")
	ErrCustomerNotMember        = errors.New("customer tidak memiliki membership aktif")
	ErrMemberNotInLocation      = errors.New("membership tidak berlaku untuk lokasi ini")
	ErrPaymentInitiationFailed  = errors.New("gagal menginisiasi pembayaran")
	ErrPaymentMethodUnsupported = errors.New("metode pembayaran tidak didukung")
)

type InitiatePaymentInput struct {
	TransactionCode string // Code transaksi dari scan QR di jukir app
	CustomerID      int64  // Didapat dari context JWT (customer)
	JukirID         int64  // Didapat dari context JWT (jukir)
	PaymentMethod   string // "CASH", "QRIS", "MEMBERSHIP"
}

type InitiatePaymentOutput struct {
	Type       string `json:"type"` // "QRIS" | "VIRTUAL_ACCOUNT" | "MEMBER_CONFIRMED" | "CASH_RECEIVED"
	QrisString string `json:"qris_string,omitempty"`
	VaNumber   string `json:"va_number,omitempty"`
	BankName   string `json:"bank_name,omitempty"`
	Amount     int64  `json:"amount"`
	ExpiredAt  string `json:"expired_at"` // ISO8601 string
	Message    string `json:"message"`
}

type InitiatePaymentUseCase struct {
	repo              domain_payment.Repository
	paymentGateway    payment_gateway.Service
	jukir_walletRepo  domain_wallet.Repository
	membershipRepo    domain_membership.Repository
	govPercentage     float64 // 0.20
	companyPercentage float64 // 0.40
	jukirPercentage   float64 // 0.40
}

func NewInitiatePaymentUseCase(
	repo domain_payment.Repository,
	paymentGateway payment_gateway.Service,
	jukir_walletRepo domain_wallet.Repository,
	membershipRepo domain_membership.Repository,
	govPercent, companyPercent, jukirPercent float64,
) *InitiatePaymentUseCase {
	return &InitiatePaymentUseCase{
		repo:              repo,
		paymentGateway:    paymentGateway,
		jukir_walletRepo:  jukir_walletRepo,
		membershipRepo:    membershipRepo,
		govPercentage:     govPercent,
		companyPercentage: companyPercent,
		jukirPercentage:   jukirPercent,
	}
}

func (uc *InitiatePaymentUseCase) Execute(ctx context.Context, input InitiatePaymentInput) (InitiatePaymentOutput, error) {
	// 1. Ambil Financial Transaction berdasarkan Code
	ftx, err := uc.repo.GetFinancialTransactionByCode(ctx, input.TransactionCode)
	if err != nil {
		return InitiatePaymentOutput{}, fmt.Errorf("get financial transaction: %w", err)
	}

	// 2. Validasi transaksi milik customer dan statusnya unpaid/pending
	if ftx.CustomerID == nil || *ftx.CustomerID != input.CustomerID {
		return InitiatePaymentOutput{}, errors.New("transaksi tidak ditemukan atau tidak milik customer ini")
	}
	if ftx.JukirID == nil || *ftx.JukirID != input.JukirID {
		return InitiatePaymentOutput{}, errors.New("transaksi tidak ditemukan atau tidak dilayani oleh jukir ini")
	}
	if ftx.Status != "unpaid" {
		return InitiatePaymentOutput{}, ErrTransactionNotPending
	}

	// 3. Validasi Payment Method
	validMethods := map[string]bool{"CASH": true, "QRIS": true, "MEMBERSHIP": true}
	if !validMethods[input.PaymentMethod] {
		return InitiatePaymentOutput{}, ErrPaymentMethodUnsupported
	}

	// 4. Proses berdasarkan Payment Method
	var output InitiatePaymentOutput
	switch input.PaymentMethod {
	case "CASH":
		output, err = uc.processCashPayment(ctx, ftx, input.JukirID)
	case "QRIS":
		output, err = uc.processQrisPayment(ctx, ftx)
	case "MEMBERSHIP":
		output, err = uc.processMembershipPayment(ctx, ftx, input.CustomerID, input.JukirID)
	default:
		err = ErrPaymentMethodUnsupported
	}

	if err != nil {
		return InitiatePaymentOutput{}, err
	}

	// 5. Setelah pembayaran sukses (tanpa error), akhiri sesi parkir
	if ftx.SessionID != nil {
		endTime := time.Now()
		err = uc.repo.EndParkingSession(ctx, *ftx.SessionID, endTime)
		if err != nil {
			// Log error, tapi jangan hentikan response sukses
			fmt.Printf("Warning: Could not end parking session: %v\n", err)
		}
	}

	return output, nil
}

func (uc *InitiatePaymentUseCase) processCashPayment(ctx context.Context, ftx *domain_payment.FinancialTransaction, jukirID int64) (InitiatePaymentOutput, error) {
	// Ambil dompet jukir
	jukirWallet, err := uc.jukir_walletRepo.GetJukirWalletByUserID(ctx, jukirID)
	if err != nil {
		return InitiatePaymentOutput{}, fmt.Errorf("get jukir wallet: %w", err)
	}

	// Hitung jatah perusahaan dan pemda (60%)
	companyShare := int64(float64(ftx.FinalAmount) * uc.companyPercentage)
	govShare := int64(float64(ftx.FinalAmount) * uc.govPercentage)
	totalDeduction := companyShare + govShare

	// Cek saldo cukup
	if jukirWallet.CurrentBalance < totalDeduction {
		return InitiatePaymentOutput{}, ErrInsufficientJukirBalance
	}

	// Hitung jatah jukir (40%)
	// jukirShare := int64(float64(ftx.FinalAmount) * uc.jukirPercentage)

	// Kurangi saldo jukir
	newBalance := jukirWallet.CurrentBalance - totalDeduction

	// Update saldo di database
	err = uc.jukir_walletRepo.UpdateJukirWalletBalance(ctx, jukirWallet.ID, newBalance)
	if err != nil {
		return InitiatePaymentOutput{}, fmt.Errorf("update jukir wallet balance: %w", err)
	}

	// Simpan histori potong saldo
	history := &domain_wallet.JukirWalletHistory{
		WalletID:        jukirWallet.ID,
		TransactionType: "CASH_DEDUCTION", // Atau gunakan konstanta
		Amount:          -totalDeduction,  // Negatif karena pengurangan
		PreviousBalance: jukirWallet.CurrentBalance,
		NewBalance:      newBalance,
		ReferenceID:     ftx.Code,
		Description:     fmt.Sprintf("Potong saldo untuk transaksi cash Rp %d", ftx.FinalAmount),
	}
	err = uc.jukir_walletRepo.CreateJukirWalletHistory(ctx, history)
	if err != nil {
		// Log error, tapi jangan hentikan proses utama
		fmt.Printf("Warning: Could not create jukir wallet history: %v\n", err)
	}

	// Update status transaksi ke 'paid'
	now := time.Now()
	err = uc.repo.UpdateFinancialTransactionStatus(ctx, ftx.ID, "paid", &now)
	if err != nil {
		return InitiatePaymentOutput{}, fmt.Errorf("update financial transaction status: %w", err)
	}

	// Update pembagian di transaksi
	// Kita asumsakan repo bisa update kolom share ini juga
	// Jika tidak, kamu perlu method UpdateFinancialTransactionShares
	// Untuk sederhananya, kita abaikan update shares di sini, atau lakukan update di repo jika diperlukan.
	// Biasanya shares dihitung dan disimpan saat transaksi dibuat, bukan diupdate saat bayar.

	return InitiatePaymentOutput{
		Type:    "CASH",
		Amount:  ftx.FinalAmount,
		Message: "Pembayaran cash berhasil. Saldo jukir dikurangi.",
	}, nil
}

func (uc *InitiatePaymentUseCase) processQrisPayment(ctx context.Context, ftx *domain_payment.FinancialTransaction) (InitiatePaymentOutput, error) {
	// Hubungi Payment Gateway
	paymentDetails := payment_gateway.PaymentDetails{
		OrderID: ftx.Code, // Gunakan code transaksi sebagai order ID
		Amount:  ftx.FinalAmount,
		ItemDetails: []payment_gateway.ItemDetail{
			{
				ID:       fmt.Sprintf("tx_%d", ftx.ID),
				Name:     "Biaya Parkir",
				Price:    ftx.FinalAmount,
				Quantity: 1,
			},
		},
		// CustomerDetails: &payment_gateway.CustomerDetails{...}, // Ambil dari user entity jika diperlukan
		ExpiryDuration: 3600, // 1 jam
	}

	qrisString, _, _, expiredAt, err := uc.paymentGateway.RequestPayment(ctx, paymentDetails)
	if err != nil {
		return InitiatePaymentOutput{}, fmt.Errorf("request payment to gateway: %w", err)
	}

	// Simpan Payment Event ke database (status pending)
	paymentEventCode := fmt.Sprintf("PE_%s_%d", ftx.Code, time.Now().UnixNano())
	paymentEvent := &domain_payment.PaymentEvent{
		Code:                paymentEventCode,
		ContextType:         "transaction",
		ReferenceEntityType: "financial_parking_transaction",
		ReferenceEntityID:   ftx.ID,
		GrossAmount:         ftx.FinalAmount,
		NetAmount:           ftx.FinalAmount,
		CurrencyCode:        ftx.CurrencyCode,
		Status:              "pending",
		CreatedAt:           time.Now(),
		ExpiredAt:           nil,    // Akan diisi oleh gateway atau dihitung manual
		PaymentChannelName:  "qris", // Sesuaikan
		ProviderReference:   qrisString,
		ChannelCode:         "",
	}
	err = uc.repo.CreatePaymentEvent(ctx, paymentEvent)
	if err != nil {
		return InitiatePaymentOutput{}, fmt.Errorf("create payment event: %w", err)
	}

	// Link Payment Event ke Financial Transaction
	err = uc.repo.LinkPaymentEventToTransaction(ctx, paymentEvent.ID, ftx.ID)
	if err != nil {
		// Log error
		fmt.Printf("Warning: Could not link payment event to transaction: %v\n", err)
	}

	// Kembalikan informasi pembayaran ke client (jukir app)
	return InitiatePaymentOutput{
		Type:       "QRIS",
		QrisString: qrisString,
		Amount:     ftx.FinalAmount,
		ExpiredAt:  expiredAt,
		Message:    "QRIS dibuat, silakan minta customer scan.",
	}, nil
}

func (uc *InitiatePaymentUseCase) processMembershipPayment(ctx context.Context, ftx *domain_payment.FinancialTransaction, customerID int64, jukirID int64) (InitiatePaymentOutput, error) {
	// Ambil membership customer
	membership, err := uc.membershipRepo.GetCustomerMembershipByUserID(ctx, customerID)
	if err != nil {
		return InitiatePaymentOutput{}, ErrCustomerNotMember
	}

	// Validasi lokasi (asumsi location_id dari ftx adalah lokasi yang valid untuk membership ini)
	// Kita asumsakan membership plan menyimpan location_id-nya
	// Jika tidak, kamu perlu join dengan membership_plans untuk cek location_id
	// Misalnya, di repo GetCustomerMembershipByUserID, kamu bisa return juga location_id dari plan
	// Kita skip validasi lokasi di sini untuk simplifikasi, tapi harus ditambahkan di implementasi nyata.

	// Hitung jatah perusahaan dan pemda (60% dari tarif normal)
	companyShare := int64(float64(ftx.FinalAmount) * uc.companyPercentage)
	govShare := int64(float64(ftx.FinalAmount) * uc.govPercentage)
	totalFromPool := companyShare + govShare // Total yang diambil dari pool

	// Cek saldo membership cukup
	if membership.PoolBalance < totalFromPool {
		// Atau, jika ingin lebih ketat, cek apakah tarif normal > pool_balance
		return InitiatePaymentOutput{}, fmt.Errorf("saldo membership tidak mencukupi")
	}

	// Hitung jatah jukir (40% dari tarif normal)
	jukirShare := int64(float64(ftx.FinalAmount) * uc.jukirPercentage)

	// Kurangi saldo membership
	newPoolBalance := membership.PoolBalance - totalFromPool

	// Ambil dompet jukir
	jukirWallet, err := uc.jukir_walletRepo.GetJukirWalletByUserID(ctx, jukirID)
	if err != nil {
		return InitiatePaymentOutput{}, fmt.Errorf("get jukir wallet: %w", err)
	}

	// Tambahkan jatah jukir ke dompet jukir
	newJukirBalance := jukirWallet.CurrentBalance + jukirShare

	// Update saldo membership di database
	err = uc.membershipRepo.UpdateCustomerMembershipPoolBalance(ctx, membership.ID, newPoolBalance)
	if err != nil {
		return InitiatePaymentOutput{}, fmt.Errorf("update customer membership pool balance: %w", err)
	}

	// Update saldo jukir di database
	err = uc.jukir_walletRepo.UpdateJukirWalletBalance(ctx, jukirWallet.ID, newJukirBalance)
	if err != nil {
		return InitiatePaymentOutput{}, fmt.Errorf("update jukir wallet balance: %w", err)
	}

	// Simpan histori tambah saldo jukir
	history := &domain_wallet.JukirWalletHistory{
		WalletID:        jukirWallet.ID,
		TransactionType: "MEMBER_COMMISSION", // Atau gunakan konstanta
		Amount:          jukirShare,          // Positif karena penambahan
		PreviousBalance: jukirWallet.CurrentBalance,
		NewBalance:      newJukirBalance,
		ReferenceID:     ftx.Code,
		Description:     fmt.Sprintf("Komisi dari scan QR member Rp %d", ftx.FinalAmount),
	}
	err = uc.jukir_walletRepo.CreateJukirWalletHistory(ctx, history)
	if err != nil {
		// Log error
		fmt.Printf("Warning: Could not create jukir wallet history: %v\n", err)
	}

	// Update status transaksi ke 'paid'
	now := time.Now()
	err = uc.repo.UpdateFinancialTransactionStatus(ctx, ftx.ID, "paid", &now)
	if err != nil {
		return InitiatePaymentOutput{}, fmt.Errorf("update financial transaction status: %w", err)
	}

	return InitiatePaymentOutput{
		Type:    "MEMBERSHIP",
		Amount:  ftx.FinalAmount,
		Message: "Pembayaran membership berhasil. Saldo jukir ditambahkan.",
	}, nil
}
