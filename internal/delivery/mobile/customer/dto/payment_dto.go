package dto

// import "time"

// Request
type SubmitQrRequest struct {
	ScannedQrString string `json:"scanned_qr_string"`
}

type ExecutePaymentRequest struct {
	CustomerID string `json:"customer_id"`
	Total      int64  `json:"total"`
	SessionID  string `json:"session_id"` // Kirim session ID dari tahap sebelumnya
}

// Response
type ScanDetailDto struct {
	CustomerId   string         `json:"customer_id"`
	CustomerName string         `json:"customer_name"`
	Lokasi       string         `json:"lokasi"`
	Duration     string         `json:"duration"`
	IsMember     bool           `json:"is_member"`
	Total        int64          `json:"total"`
	Breakdown    []PriceItemDto `json:"breakdown"`
}

type PriceItemDto struct {
	Label  string `json:"label"`
	Amount int64  `json:"amount"`
}

type PaymentInfoDto struct {
	Type       string `json:"type"` // "QRIS" | "VIRTUAL_ACCOUNT"
	QrisString string `json:"qris_string,omitempty"`
	VaNumber   string `json:"va_number,omitempty"`
	BankName   string `json:"bank_name,omitempty"`
	Amount     int64  `json:"amount"`
	ExpiredAt  string `json:"expired_at"` // ISO8601 string
}
