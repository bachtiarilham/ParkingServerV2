package helper

type NominalPaymentResponseDto struct {
	Nominal *[]NominalItemDto `json:"nominal"`
}

type NominalItemDto struct {
	OptionID      int64  `json:"optionId"`
	NominalAmount int64  `json:"nominalAmout"` // Typo dari request (nominalAmout) dipertahankan di json tag
	Label         string `json:"label"`
}
