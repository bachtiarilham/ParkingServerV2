package helper

type TopupResponseModel struct {
	Nominal     *[]NominalItemModel `json:"nominal"`
	MetodeBayar *[]MetodeItemModel  `json:"metode_bayar"`
}

type NominalItemModel struct {
	OptionID      int64  `json:"optionId"`
	NominalAmount int64  `json:"nominalAmout"`
	Label         string `json:"label"`
}

type MetodeItemModel struct {
	PaymentMethodId int64  `json:"method_id"`
	NamaPayment     string `json:"nama_payment"`
	CodePayment     string `json:"code_payment"`
	LogoPayment     string `json:"logo_payment"`
}
