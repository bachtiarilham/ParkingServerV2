package topup

type TopUpRequestModel struct {
	IDUser       int `json:"id_user"`
	NominalTopUp int `json:"nominal_topup"`
}
