package paymentgate

type PayRequestModel struct {
	PaymentType       string  `json:"payment_type"`
	TargetID          *string `json:"target_id"`
	Amount            int64   `json:"amount"`
	PaymentMethodCode string  `json:"payment_method_code"`
	PromoCode         string  `json:"promo_code"`
	UserID            int64   `json:"user_id"`

	JukirShare    int64
	CompanyShare  int64
	GovShare      int64
	MidtransShare int64
}
