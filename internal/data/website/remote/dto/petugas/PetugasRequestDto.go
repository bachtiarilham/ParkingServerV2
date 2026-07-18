package petugas

type PetugasResponseDto struct {
	Petugas *[]PetugasItemDto `json:"petugas"`
}

type PetugasItemDto struct {
	ID              int    `json:"id"`
	Nama            string `json:"nama"`
	JmlTransaksi    int    `json:"jml_transaksi"`
	TotalPendapatan int    `json:"total_pendapatan"`
	IsAktif         bool   `json:"is_aktif"`
	Parlok          string `json:"parlok"`
}
