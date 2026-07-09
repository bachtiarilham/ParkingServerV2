package helper

type TarifModel struct {
	TarifItem *[]TarifItem
}

type TarifItem struct {
	KendaraanId   int64
	KendaraanKode string
	KendaraanNama string `json:"kendaraan"`
	Nominal       int64  `json:"nominal"`
}
