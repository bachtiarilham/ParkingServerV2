package topup

type TopUpRequestDto struct {
	NominalTopUp int `json:"nominal_topup"`
	IDUser       int `json:"id_user"`
}
