package helper

type PaymentMethodModel struct {
	MetodeBayar *[]MetodeItemModel `json:"metode_bayar"`
}

type MetodeItemModel struct {
	PaymentMethodId int64  `json:"method_id"`
	NamaPayment     string `json:"nama_payment"`
	CodePayment     string `json:"code_payment"`
	LogoPayment     string `json:"logo_payment"`
}
