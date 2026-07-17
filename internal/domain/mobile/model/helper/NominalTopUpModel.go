package helper

type TopupOptionsResponseModel struct {
	Nominal *[]TopupOptionItemModel `json:"nominal"`
}

type TopupOptionItemModel struct {
	OptionID      int64  `json:"optionId"`
	NominalAmount int64  `json:"nominalAmout"`
	Label         string `json:"label"`
}
