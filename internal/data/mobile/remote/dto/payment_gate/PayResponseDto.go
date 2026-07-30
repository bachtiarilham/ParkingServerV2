package paymentgate

type PayResponseDto struct {
	OrderID       string `json:"order_id"`
	GrossAmount   int64  `json:"gross_amount"`
	PaymentMethod string `json:"payment_method"` // misal: "gopay", "bca_va", "indomaret", "cash"

	Status     string `json:"status"`      // "PENDING", "PAID", "WAITING_FOR_CASH"
	ExpiryTime string `json:"expiry_time"` // Format ISO 8601 / RFC 3339

	// SnapToken & RedirectURL digunakan jika Anda memakai Midtrans Mobile SDK / Webview
	SnapToken   string `json:"snap_token,omitempty"`
	RedirectURL string `json:"redirect_url,omitempty"`

	// PaymentAction memuat data spesifik sesuai dari 6 grup pembayaran yang dipilih
	PaymentAction *PaymentActionDto `json:"payment_action,omitempty"`
}

// PaymentAction memuat field dinamis. Gunakan pointer (*) dan omitempty
// agar field yang kosong tidak muncul di JSON response mobile.
type PaymentActionDto struct {
	// 1. Khusus Grup E-Wallet (GoPay, ShopeePay, DANA)
	DeepLinkURL string `json:"deep_link_url,omitempty"` // URL untuk buka aplikasi (gojek://...)

	// 3. Khusus Grup Indomaret / Alfamart (Over-The-Counter / OTC)
	StoreName   string `json:"store_name,omitempty"`   // "INDOMARET" atau "ALFAMART"
	PaymentCode string `json:"payment_code,omitempty"` // Kode bayar yang ditunjukkan ke kasir

	// 4. Khusus Grup QRIS
	QRCodeString string `json:"qr_code_string,omitempty"` // String asli QR untuk digenerate lokal di HP
	QRCodeURL    string `json:"qr_code_url,omitempty"`    // URL gambar QR Code siap pakai dari Midtrans

	// 5. Khusus Grup Cash (Tunai ke Driver / Non-Midtrans)
	Instruction string `json:"instruction,omitempty"` // Misal: "Siapkan uang tunai pas untuk driver"
}
