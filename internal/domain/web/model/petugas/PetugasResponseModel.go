package petugas

type PetugasResponseModel struct {
	Petugas *[]PetugasItemModel `json:"petugas"`
}

type PetugasItemModel struct {
	ID              int    `json:"id"`
	Nama            string `json:"nama"`
	JmlTransaksi    int    `json:"jml_transaksi"`
	TotalPendapatan int    `json:"total_pendapatan"`
	IsAktif         bool   `json:"is_aktif"`
	Parlok          string `json:"parlok"`
}
