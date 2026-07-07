package payment

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"math"
// 	"modulegue/internal/domain/payment"
// 	"time"
// )

// type ExecutePaymentInput struct {
// 	CustomerID int64 // Didapat dari context JWT
// 	Total      int64
// 	SessionID  int64 // Dari tahap sebelumnya
// }

// type ExecutePaymentOutput struct {
// 	Type       string `json:"type"` // "QRIS" | "VIRTUAL_ACCOUNT"
// 	QrisString string `json:"qris_string,omitempty"`
// 	VaNumber   string `json:"va_number,omitempty"`
// 	BankName   string `json:"bank_name,omitempty"`
// 	Amount     int64  `json:"amount"`
// 	ExpiredAt  string `json:"expired_at"` // ISO8601 string
// }

// type ExecutePaymentUseCase struct {
// 	repo payment.Repository
// }

// func NewExecutePaymentUseCase(repo payment.Repository) *ExecutePaymentUseCase {
// 	return &ExecutePaymentUseCase{repo: repo}
// }

// func (uc *ExecutePaymentUseCase) Execute(ctx context.Context, input ExecutePaymentInput) (ExecutePaymentOutput, error) {
// 	// 1. Ambil detail sesi untuk validasi total
// 	session, err := uc.repo.GetActiveSessionByCode(ctx, fmt.Sprintf("%d", input.SessionID))
// 	if err != nil {
// 		return ExecutePaymentOutput{}, fmt.Errorf("session not found: %w", err)
// 	}

// 	// 2. Hitung ulang total untuk validasi (opsional, bisa di-cache dari tahap sebelumnya)
// 	duration := time.Since(session.StartedAt)
// 	durationHours := int64(math.Ceil(duration.Hours()))
// 	tariff, err := uc.repo.GetTariffForLocationAndVehicle(ctx, session.LocationID, session.VehicleTypeID)
// 	if err != nil {
// 		return ExecutePaymentOutput{}, fmt.Errorf("failed to get tariff: %w", err)
// 	}
// 	calculatedTotal := durationHours * tariff

// 	if calculatedTotal != input.Total {
// 		return ExecutePaymentOutput{}, errors.New("calculated total does not match requested total")
// 	}

// 	// 3. Buat FinancialTransaction (belum dibayar)
// 	transactionCode := fmt.Sprintf("TX_%d_%d", session.ID, time.Now().UnixNano()) // Generate unique code
// 	financialTx := &payment.FinancialTransaction{
// 		Code:              transactionCode,
// 		OperationType:     "parking_fee",
// 		TransactionSource: "qr_scan",
// 		SessionID:         &session.ID,
// 		LocationID:        session.LocationID,
// 		CustomerID:        &session.CustomerID,
// 		SubtotalAmount:    input.Total,
// 		FinalAmount:       input.Total, // Sederhanakan
// 		CurrencyCode:      "IDR",
// 		Status:            "unpaid",
// 		OccurredAt:        time.Now(),
// 		CreatedAt:         time.Now(),
// 	}
// 	err = uc.repo.CreateFinancialTransaction(ctx, financialTx)
// 	if err != nil {
// 		return ExecutePaymentOutput{}, fmt.Errorf("failed to create financial transaction: %w", err)
// 	}

// 	// 4. Buat PaymentEvent (untuk gateway)
// 	expiryTime := time.Now().Add(1 * time.Hour) // Misalnya, QR/VA expire dalam 1 jam
// 	paymentEventCode := fmt.Sprintf("PE_%d_%d", financialTx.ID, time.Now().UnixNano())
// 	paymentEvent := &payment.PaymentEvent{
// 		Code:                paymentEventCode,
// 		ContextType:         "transaction",
// 		ReferenceEntityType: "financial_parking_transaction",
// 		ReferenceEntityID:   financialTx.ID,
// 		GrossAmount:         input.Total,
// 		NetAmount:           input.Total, // Sederhanakan
// 		CurrencyCode:        "IDR",
// 		Status:              "pending",
// 		CreatedAt:           time.Now(),
// 		ExpiredAt:           &expiryTime,
// 		PaymentChannelName:  "qris", // Bisa disesuaikan, misalnya pilih channel
// 		ProviderReference:   "",     // Akan diisi oleh gateway
// 		ChannelCode:         "",     // Akan diisi oleh gateway
// 	}
// 	err = uc.repo.CreatePaymentEvent(ctx, paymentEvent)
// 	if err != nil {
// 		return ExecutePaymentOutput{}, fmt.Errorf("failed to create payment event: %w", err)
// 	}

// 	// 5. Hubungkan PaymentEvent ke FinancialTransaction
// 	err = uc.repo.LinkPaymentEventToTransaction(ctx, paymentEvent.ID, financialTx.ID)
// 	if err != nil {
// 		return ExecutePaymentOutput{}, fmt.Errorf("failed to link payment event to transaction: %w", err)
// 	}

// 	// 6. Simulasikan pembuatan instruksi pembayaran (QRIS/VA)
// 	// Dalam dunia nyata, ini dilakukan oleh service pembayaran eksternal (Midtrans, Doku, dll)
// 	// Kita hanya return dummy data sebagai contoh
// 	// --- PROSES INI HARUS DILAKUKAN OLEH PAYMENT SERVICE ---
// 	// paymentService.CreateQRIS(paymentEvent.Code, paymentEvent.GrossAmount, expiryTime)
// 	// paymentService.CreateVA(paymentEvent.Code, paymentEvent.GrossAmount, expiryTime, bankCode)
// 	// -------------------------------------------

// 	// Contoh output simulasi QRIS
// 	return ExecutePaymentOutput{
// 		Type:       "QRIS",
// 		QrisString: "simulated_qris_string_" + paymentEvent.Code, // Diganti oleh service sebenarnya
// 		Amount:     paymentEvent.GrossAmount,
// 		ExpiredAt:  expiryTime.Format(time.RFC3339),
// 	}, nil
// }
