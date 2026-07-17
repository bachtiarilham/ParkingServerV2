package topup

import "time"

// ==============================================================================
// 1. GET /topup/options
// ==============================================================================

// ==============================================================================
// 2. POST /topup/create
// ==============================================================================

type TopupCreateRequestModel struct {
	UserID            int64     `json:"userId"`
	Amount            int64     `json:"amount"`
	PaymentMethodCode string    `json:"paymentMethodCode"`
	AdminFee          int64     `json:"adminFee"`
	TopupCode         string    `json:"topupCode"`
	ExternalReference string    `json:"externalReference"`
	ProviderName      string    `json:"providerName"`
	QRString          string    `json:"qrString"`
	ExpiredAt         time.Time `json:"expiredAt"`
}

type TopupCreateResponseModel struct {
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

type TopupStatusRequestModel struct {
	TopupCode string `uri:"topupCode" binding:"required"`
}

type TopupStatusResponseModel struct {
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
// Digabung menjadi satu request struct karena payload dari provider
// biasanya seragam dan hanya berbeda di field yang terisi.

type QrisCallbackRequestModel struct {
	TransactionTime   string `json:"transactionTime"`
	TransactionStatus string `json:"transactionStatus"`
	TransactionID     string `json:"transactionID"`
	StatusMessage     string `json:"statusMessage"`
	StatusCode        string `json:"statusCode"`
	SignatureKey      string `json:"signatureKey"`
	PaymentType       string `json:"paymentType"`
	TopupCode         string `json:"topupCode"`
	MerchantID        string `json:"merchantID"`
	GrossAmount       string `json:"grossAmount"`
	FraudStatus       string `json:"fraudStatus"`
	Currency          string `json:"currency"`
	Acquirer          string `json:"acquirer"`
	SettlementTime    string `json:"settlementTime"`
	FailedReason      string `json:"failedReason,omitempty"`
}

type QrisCallbackResponseModel struct {
	TopupTransactionID int64  `json:"topupTransactionId"`
	UserID             int64  `json:"userId"`
	WalletID           int64  `json:"walletId"`
	Amount             int64  `json:"amount"`
	PaymentStatusID    int64  `json:"payment_status_id"` // Format snake_case sesuai prompt
	PaymentStatusCode  string `json:"paymentStatusCode"`
}
