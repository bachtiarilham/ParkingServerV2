package paymentgate

type PaymentStatusModel struct {
	OrderID    string `json:"order_id"`
	Status     string `json:"status"` // "PAID", "PENDING", "FAILED"
	StatusText string `json:"status_text"`
}
