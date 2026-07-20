package topup

type TopUpRequestModel struct {
	IDUser            int     `json:"id_user"`
	NominalTopUp      int     `json:"nominal_topup"`
	TopUpCode         string  `json:"-"`
	ExternalReference string  `json:"-"`
	BalanceBefore     float64 `json:"-"`
	BalanceAfter      float64 `json:"-"`
}
