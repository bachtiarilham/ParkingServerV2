package setting

type SaveTarifRequestModel struct {
	Tarif           int    `json:"Tarif"`
	KeteranganTarif string `json:"KeteranganTarif"`
	IDLokasi        int    `json:"IDLokasi"`
}
