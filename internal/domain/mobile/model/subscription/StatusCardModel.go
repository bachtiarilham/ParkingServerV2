package subscription

type StatusCardModel struct {
	PaketAktif *string `json:"paket_aktif,omitempty"`
	Kadaluarsa *string `json:"kadaluarsa,omitempty"`
	Benefit    *string `json:"benefit,omitempty"`
}
