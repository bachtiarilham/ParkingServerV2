package helper

type GetTarifResponseModel struct {
	Tarif *[]TarifItemModel `json:"tarif"`
}

type TarifItemModel struct {
	KetTarif string `json:"ket_tarif"`
	ID       int    `json:"id"`
	Tarif    int    `json:"tarif"`
}
