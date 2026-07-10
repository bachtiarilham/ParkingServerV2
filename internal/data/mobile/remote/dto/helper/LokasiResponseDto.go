package helper

type LokasiResponseDto struct {
	LokasiItem *[]LokasiItemResponseDto `json:"lokasi_item,omitempty"`
}

type LokasiItemResponseDto struct {
	NamaLokasi string `json:"nama_lokasi"`
	LokasiId   int64  `json:"lokasi_id"`
	AreaId     int64  `json:"area_id"`
	NamaArea   string `json:"nama_area"`
	ZonaId     int64  `json:"zona_id"`
	NamaZona   string `json:"nama_zona"`
	Address    string `json:"address"`
}
