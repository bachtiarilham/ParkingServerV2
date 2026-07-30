package helper

type PaymentMethodResponseDto struct {
	MetodeBayar *[]MetodeItemDto `json:"metode_bayar"`
}

type MetodeItemDto struct {
	PaymentMethodId int64  `json:"method_id"`
	NamaPayment     string `json:"nama_payment"`
	CodePayment     string `json:"code_payment"`
	LogoPayment     string `json:"logo_payment"`
}
