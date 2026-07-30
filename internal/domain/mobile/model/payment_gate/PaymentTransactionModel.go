package paymentgate

type PaymentTransactionModel struct {
	PaymentType  string `json:"payment_type"`
	UserID       int64  `json:"user_id"`
	ReferenceID  int64  `json:"reference_id"`
	Amount       int64  `json:"amount"`
	PartnerShare int64
}
