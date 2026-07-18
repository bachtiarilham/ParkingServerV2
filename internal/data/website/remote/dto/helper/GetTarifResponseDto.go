package helper

type GetTarifResponseDto struct {
	Tarif []TarifItemDto `json:"tarif"`
}

type TarifItemDto struct {
	KetTarif string `json:"ket_tarif"`
	ID       int    `json:"id"`
	Tarif    int    `json:"tarif"`
}
