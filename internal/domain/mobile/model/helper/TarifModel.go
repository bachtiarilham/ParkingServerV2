package helper

type TarifModel struct {
	ItemTarif []TarifModelItem `json:"itemTarif,omitempty"`
}

type TarifModelItem struct {
	Kendaraan int64  `json:"kendaraan"`
	Nominal   string `json:"nominal"`
}
