package setting

type SaveTarifRequestDto struct {
	Tarif           int    `json:"Tarif"`
	KeteranganTarif string `json:"KeteranganTarif"`
	IDLokasi        int    `json:"IDLokasi"`
}

// Gunakan AsyncTaskResponse yang sudah dideklarasikan di endpoint Register
