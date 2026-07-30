package invoice

import "time"

// ==========================================
// 1. MAIN INVOICE STRUCT
// ==========================================

// UniversalInvoiceResponseModel struct utama untuk detail invoice
type UniversalInvoiceResponseModel struct {
	TrxCode         string    `json:"trx_code"`         // e.g. "TRX-PRK-14218471981"
	TransactionType string    `json:"transaction_type"` // "PARKING" | "TOPUP" | "MEMBERSHIP" | "TRANSFER"
	Title           string    `json:"title"`            // e.g. "Parkir On-Street", "Top Up Saldo", "LineSpot Gold", "Transfer"
	Status          string    `json:"status"`           // "PAID" | "PENDING" | "FAILED" | "REFUNDED"
	StatusText      string    `json:"status_text"`      // e.g. "Lunas", "Menunggu Pembayaran"
	CreatedAt       time.Time `json:"created_at"`       // e.g. 2026-07-22 15:13:00
	TotalAmount     int64     `json:"total_amount"`     // Total tagihan akhir
	InvoiceUrl      string

	// Dynamic Details (Hanya 1 yang terisi sesuai transaction_type)
	ParkingDetails    *ParkingInvoiceDetailModel    `json:"parking_details,omitempty"`
	WalletDetails     *WalletInvoiceDetailModel     `json:"wallet_details,omitempty"`
	MembershipDetails *MembershipInvoiceDetailModel `json:"membership_details,omitempty"`

	// Sub-components
	PriceBreakdown PriceBreakdownModel `json:"price_breakdown"`
	PaymentMethod  PaymentMethodModel  `json:"payment_method"`
	CustomerInfo   CustomerInfoModel   `json:"customer_info"`
}

// ==========================================
// 2. DYNAMIC DETAILS (BERUBAH SESUAI TIPE)
// ==========================================

// ParkingInvoiceDetailModel khusus transaksi parkir
type ParkingInvoiceDetailModel struct {
	LocationName  string     `json:"location_name"`            // e.g. "On-Street Area A (Jl. Sudirman)"
	LicensePlate  string     `json:"license_plate"`            // e.g. "B 1234 ABC"
	VehicleType   string     `json:"vehicle_type"`             // "Motor" | "Mobil"
	CheckInTime   time.Time  `json:"check_in_time"`            // Jam masuk
	CheckOutTime  *time.Time `json:"check_out_time,omitempty"` // Jam keluar
	DurationText  string     `json:"duration_text"`            // e.g. "1 Jam 20 Menit"
	AttendantName string     `json:"attendant_name,omitempty"` // e.g. "Budi Hartono" (Jukir/Petugas)
}

// WalletInvoiceDetailModel khusus transaksi top-up atau transfer
type WalletInvoiceDetailModel struct {
	SenderName       string `json:"sender_name"`       // e.g. "BCA Virtual Account" atau "Ilham M."
	SenderAccount    string `json:"sender_account"`    // e.g. "**** 7432"
	RecipientName    string `json:"recipient_name"`    // e.g. "LineSpot Pay Balance"
	RecipientAccount string `json:"recipient_account"` // e.g. "+6285168107712"
	BankRefNo        string `json:"bank_ref_no,omitempty"`
}

// MembershipInvoiceDetailModel khusus transaksi paket berlangganan
type MembershipInvoiceDetailModel struct {
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

// PriceBreakdownModel rincian komponen harga & biaya
type PriceBreakdownModel struct {
	BasePrice      int64  `json:"base_price"`                // Subtotal utama / Tarif dasar
	AdminFee       int64  `json:"admin_fee,omitempty"`       // Biaya layanan app / admin
	TaxAmount      int64  `json:"tax_amount,omitempty"`      // PPN (11%)
	DiscountAmount int64  `json:"discount_amount,omitempty"` // Potongan voucher / cashback
	DiscountCode   string `json:"discount_code,omitempty"`   // Kode promo
	FinalTotal     int64  `json:"final_total"`               // Total bersih
}

// PaymentMethodModel rincian sumber pembayaran yang digunakan
type PaymentMethodModel struct {
	ChannelCode string `json:"channel_code"`         // "LINESPOT_PAY" | "BCA_VA" | "GOPAY"
	ChannelName string `json:"channel_name"`         // e.g. "Saldo LineSpot Pay", "BCA Virtual Account"
	AccountNo   string `json:"account_no,omitempty"` // Masked ID e.g. "**** 7432"
	IconURL     string `json:"icon_url,omitempty"`
}

// CustomerInfoModel informasi dasar pembeli/pengguna
type CustomerInfoModel struct {
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"` // e.g. "Ilham Muhamad Bachtiar"
	Email    string `json:"email"`     // e.g. "ilhammb7@gmail.com"
	Phone    string `json:"phone"`     // e.g. "+6285168107712"
}
