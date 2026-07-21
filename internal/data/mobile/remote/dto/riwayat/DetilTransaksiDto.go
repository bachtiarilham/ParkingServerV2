package riwayat

type DetilTransaksiDto struct {
	Tanggal            string `json:"tanggal"`
	TopUpTransactionID int64  `json:"topup_transaction_id"`
	TopUpCode          string `json:"topup_code"`
	UserID             int64  `json:"user_id"`
	WalletID           int64  `json:"wallet_id"`
	PaymentMethodID    int64  `json:"payment_method_id"`
	PaymentMethodCode  string `json:"payment_method_code"`
	PaymentMethodName  string `json:"payment_method_name"`
	Amount             int64  `json:"amount"`
	AdminFee           int64  `json:"admin_fee"`
	TotalAmount        int64  `json:"total_amount"`
	TransactionStatus  string `json:"transaction_status"`
	ExternalReference  string `json:"external_reference,omitempty"`
	ProviderName       string `json:"provider_name,omitempty"`
	CreatedAt          string `json:"created_at"`
	ExpiredAt          string `json:"expired_at,omitempty"`
	PaidAt             string `json:"paid_at,omitempty"`
	CompletedAt        string `json:"completed_at,omitempty"`
	FailedReason       string `json:"failed_reason,omitempty"`
}
