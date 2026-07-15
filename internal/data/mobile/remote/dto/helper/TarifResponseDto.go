package helper

type TarifResponseDto struct {
	TarifResponseItemDto []TarifResponseItemDto `json:"tarif_item"`
}

type TarifResponseItemDto struct {
	KendaraanId   int64  `json:"kendaraan_id"`
	KendaraanKode string `json:"kendaraan_kode"`
	KendaraanNama string `json:"kendaraan_nama"`
	Nominal       int64  `json:"nominal"`
}
