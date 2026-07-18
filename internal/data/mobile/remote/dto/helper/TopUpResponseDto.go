package helper

type TopupResponseDto struct {
	Nominal     *[]NominalItemDto `json:"nominal"`
	MetodeBayar *[]MetodeItemDto  `json:"metode_bayar"`
}

type NominalItemDto struct {
	OptionID      int64  `json:"optionId"`
	NominalAmount int64  `json:"nominalAmout"` // Typo dari request (nominalAmout) dipertahankan di json tag
	Label         string `json:"label"`
}

type MetodeItemDto struct {
	PaymentMethodId int64  `json:"method_id"`
	NamaPayment     string `json:"nama_payment"`
	CodePayment     string `json:"code_payment"`
	LogoPayment     string `json:"logo_payment"`
}
