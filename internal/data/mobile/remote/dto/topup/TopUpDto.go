package topup

import "time"

// ==============================================================================
// 2. POST /topup/create
// ==============================================================================

type TopupCreateRequestDto struct {
	Amount            int64  `json:"amount"`
	PaymentMethodCode string `json:"paymentMethodCode"`
}

type TopupCreateResponseDto struct {
	TopupTransactionID int64     `json:"topupTransactionId"`
	TopupCode          string    `json:"topupCode"`
	Amount             int64     `json:"amount"`
	AdminFee           int64     `json:"adminFee"`
	TotalAmount        int64     `json:"totalAmount"`
	PaymentMethodCode  string    `json:"paymentMethodCode"`
	PaymentMethodName  string    `json:"paymentMethodName"`
	PaymentStatusCode  string    `json:"paymentStatusCode"`
	PaymentStatusName  string    `json:"paymentStatusName"`
	QRString           string    `json:"qrString"`
	ExpiredAt          time.Time `json:"expiredAt"`
	CreatedAt          time.Time `json:"createdAt"`
}

// ==============================================================================
// 3. GET /topup/{topupCode}/status
// ==============================================================================
// Request: Menggunakan path parameter (tidak butuh struct khusus) atau URI struct

type TopupStatusResponseDto struct {
	TopupTransactionID int64      `json:"topupTransactionId"`
	TopupCode          string     `json:"topupCode"`
	Amount             int64      `json:"amount"`
	AdminFee           int64      `json:"adminFee"`
	TotalAmount        int64      `json:"totalAmount"`
	PaymentMethodCode  string     `json:"paymentMethodCode"`
	PaymentStatusCode  string     `json:"paymentStatusCode"`
	PaymentMethodName  string     `json:"paymentMethodName"`
	QRString           string     `json:"qrString"`
	PaidAt             *time.Time `json:"paidAt,omitempty"` // Pointer karena bisa null jika belum dibayar
	ExpiredAt          time.Time  `json:"expiredAt"`
	FailedReason       string     `json:"failedReason,omitempty"`
	CompletedAt        *time.Time `json:"completedAt,omitempty"` // Pointer karena bisa null
	CurrentBalance     int64      `json:"currentBalance"`
}

// ==============================================================================
// 4. POST /payment/qris/callback
// ==============================================================================
// Midtrans transaction notification payload.

type QrisCallbackRequestDto struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status,omitempty"`
	Currency          string `json:"currency,omitempty"`
	Acquirer          string `json:"acquirer,omitempty"`
	SettlementTime    string `json:"settlement_time,omitempty"`
}

type QrisCallbackResponseDto struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status,omitempty"`
	Currency          string `json:"currency,omitempty"`
	Acquirer          string `json:"acquirer,omitempty"`
	SettlementTime    string `json:"settlement_time,omitempty"`
}
