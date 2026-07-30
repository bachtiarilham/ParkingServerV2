package helper

type NominalPaymentModel struct {
	Nominal *[]NominalItemModel `json:"nominal"`
}

type NominalItemModel struct {
	OptionID      int64  `json:"optionId"`
	NominalAmount int64  `json:"nominalAmout"`
	Label         string `json:"label"`
}
