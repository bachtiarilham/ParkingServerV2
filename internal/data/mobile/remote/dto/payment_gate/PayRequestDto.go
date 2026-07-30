package paymentgate

type PayRequestDto struct {
	TargetID          *string `json:"target_id"`
	PaymentType       string  `json:"payment_type"`        // PARKIR, TOP_UP, MEMBERSHIP
	PaymentMethodCode string  `json:"payment_method_code"` // WALLET, BCA_VA, GOPAY, dll
	Amount            int64   `json:"amount"`              // Jumlah Nominal
	PromoCode         string  `json:"promo_code"`
}
