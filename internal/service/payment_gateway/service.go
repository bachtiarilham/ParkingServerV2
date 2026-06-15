package payment_gateway

import "context"

type PaymentDetails struct {
	OrderID         string
	Amount          int64
	ItemDetails     []ItemDetail
	CustomerDetails *CustomerDetails
	ExpiryDuration  int64 // Durasi dalam detik
}

type ItemDetail struct {
	ID       string
	Name     string
	Price    int64
	Quantity int32
}

type CustomerDetails struct {
	FName string
	LName string
	Email string
	Phone string
}

type Service interface {
	RequestPayment(ctx context.Context, details PaymentDetails) (QRString string, VAString string, BankName string, ExpiredAt string, err error)
	VerifyPayment(ctx context.Context, orderID string) (Status string, PaidAt *string, err error) // PaidAt dalam format ISO8601 string
}
