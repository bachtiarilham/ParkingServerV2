package home

type CustomerSummaryModel struct {
	Saldo       int64  `json:"saldo"` // Asumsi dalam satuan terkecil (misalnya Rupiah)
	ExpiredDate string `json:"expiredDate"`
}
