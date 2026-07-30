package invoice

import "time"

// ==========================================
// 1. MAIN INVOICE STRUCT
// ==========================================

// UniversalInvoiceResponseDto struct utama untuk detail invoice
type UniversalInvoiceResponseDto struct {
	TrxCode         string    `json:"trx_code"`         // e.g. "INV/20260722/LSP-0091"
	TransactionType string    `json:"transaction_type"` // "PARKING" | "WALLET" | "MEMBERSHIP"
	Title           string    `json:"title"`            // e.g. "Parkir On-Street", "Top Up Saldo", "LineSpot Gold"
	Status          string    `json:"status"`           // "PAID" | "PENDING" | "FAILED" | "REFUNDED"
	StatusText      string    `json:"status_text"`      // e.g. "Lunas", "Menunggu Pembayaran"
	CreatedAt       time.Time `json:"created_at"`       // e.g. 2026-07-22 15:13:00
	TotalAmount     int64     `json:"total_amount"`     // Total tagihan akhir
	InvoiceUrl      string    `json:"invoice_url"`

	// Dynamic Details (Hanya 1 yang terisi sesuai transaction_type)
	ParkingDetails    *ParkingInvoiceDetailDto    `json:"parking_details,omitempty"`
	WalletDetails     *WalletInvoiceDetailDto     `json:"wallet_details,omitempty"`
	MembershipDetails *MembershipInvoiceDetailDto `json:"membership_details,omitempty"`

	// Sub-components
	PriceBreakdown PriceBreakdownDto `json:"price_breakdown"`
	PaymentMethod  PaymentMethodDto  `json:"payment_method"`
	CustomerInfo   CustomerInfoDto   `json:"customer_info"`
}

// ==========================================
// 2. DYNAMIC DETAILS (BERUBAH SESUAI TIPE)
// ==========================================

// ParkingInvoiceDetailDto khusus transaksi parkir
type ParkingInvoiceDetailDto struct {
	LocationName  string     `json:"location_name"`            // e.g. "On-Street Area A (Jl. Sudirman)"
	LicensePlate  string     `json:"license_plate"`            // e.g. "B 1234 ABC"
	VehicleType   string     `json:"vehicle_type"`             // "Motor" | "Mobil"
	CheckInTime   time.Time  `json:"check_in_time"`            // Jam masuk
	CheckOutTime  *time.Time `json:"check_out_time,omitempty"` // Jam keluar
	DurationText  string     `json:"duration_text"`            // e.g. "1 Jam 20 Menit"
	AttendantName string     `json:"attendant_name,omitempty"` // e.g. "Budi Hartono" (Jukir/Petugas)
}

// WalletInvoiceDetailDto khusus transaksi top-up atau transfer
type WalletInvoiceDetailDto struct {
	SenderName       string `json:"sender_name"`       // e.g. "BCA Virtual Account" atau "Ilham M."
	SenderAccount    string `json:"sender_account"`    // e.g. "**** 7432"
	RecipientName    string `json:"recipient_name"`    // e.g. "LineSpot Pay Balance"
	RecipientAccount string `json:"recipient_account"` // e.g. "+6285168107712"
	BankRefNo        string `json:"bank_ref_no,omitempty"`
}

// MembershipInvoiceDetailDto khusus transaksi paket berlangganan
type MembershipInvoiceDetailDto struct {
	PackageName     string     `json:"package_name"`  // e.g. "LineSpot Gold Personal"
	PeriodStart     time.Time  `json:"period_start"`  // Tanggal mulai
	PeriodEnd       time.Time  `json:"period_end"`    // Tanggal berakhir
	MaxVehicles     int        `json:"max_vehicles"`  // e.g. 2 (Plat nomor)
	IsAutoRenew     bool       `json:"is_auto_renew"` // Perpanjangan otomatis
	NextBillingDate *time.Time `json:"next_billing_date,omitempty"`
}

// ==========================================
// 3. COMMON SUB-COMPONENTS
// ==========================================

// PriceBreakdownDto rincian komponen harga & biaya
type PriceBreakdownDto struct {
	BasePrice      int64  `json:"base_price"`                // Subtotal utama / Tarif dasar
	AdminFee       int64  `json:"admin_fee,omitempty"`       // Biaya layanan app / admin
	TaxAmount      int64  `json:"tax_amount,omitempty"`      // PPN (11%)
	DiscountAmount int64  `json:"discount_amount,omitempty"` // Potongan voucher / cashback
	DiscountCode   string `json:"discount_code,omitempty"`   // Kode promo
	FinalTotal     int64  `json:"final_total"`               // Total bersih
}

// PaymentMethodDto rincian sumber pembayaran yang digunakan
type PaymentMethodDto struct {
	ChannelCode string `json:"channel_code"`         // "LINESPOT_PAY" | "BCA_VA" | "GOPAY"
	ChannelName string `json:"channel_name"`         // e.g. "Saldo LineSpot Pay", "BCA Virtual Account"
	AccountNo   string `json:"account_no,omitempty"` // Masked ID e.g. "**** 7432"
	IconURL     string `json:"icon_url,omitempty"`
}

// CustomerInfoDto informasi dasar pembeli/pengguna
type CustomerInfoDto struct {
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"` // e.g. "Ilham Muhamad Bachtiar"
	Email    string `json:"email"`     // e.g. "ilhammb7@gmail.com"
	Phone    string `json:"phone"`     // e.g. "+6285168107712"
}
