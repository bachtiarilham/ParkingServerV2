package helper

type TopupOptionsResponseDto struct {
	Nominal *[]TopupOptionItemDto `json:"nominal"`
}

type TopupOptionItemDto struct {
	OptionID      int64  `json:"optionId"`
	NominalAmount int64  `json:"nominalAmout"` // Typo dari request (nominalAmout) dipertahankan di json tag
	Label         string `json:"label"`
}
