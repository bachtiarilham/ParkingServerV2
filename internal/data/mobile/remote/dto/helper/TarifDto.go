package helper

type TarifDto struct {
	ItemTarif []TarifItemDto `json:"tarif"`
}
type TarifItemDto struct {
	Kendaraan int64  `json:"kendaraan"`
	Nominal   string `json:"nominal"`
}
