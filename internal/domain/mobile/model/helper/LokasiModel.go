package helper

type LokasiModel struct {
	LokasiItem *[]LokasiItem
}

type LokasiItem struct {
	NamaLokasi    string `json:"lokasi,omitempty"`
	LokasiId      int64
	AreaId        int64
	NamaArea      string
	ZonaId        int64
	NamaZona      string
	Address       string
	IsCurrentArea bool
}
