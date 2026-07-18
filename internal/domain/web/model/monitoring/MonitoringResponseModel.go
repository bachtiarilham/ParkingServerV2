package monitoring

type MonitoringResponseModel struct {
	Parlok    *[]ParlokItemModel    `json:"parlok"`
	Transaksi *[]TransaksiItemModel `json:"transaksi"`
}

type ParlokItemModel struct {
	NamaParlok      string `json:"nama_parlok"`
	IDZona          int    `json:"id_zona"`
	NamaZona        string `json:"nama_zona"`
	LatMin          string `json:"lat_min"`
	LatMax          string `json:"lat_max"`
	LngMin          string `json:"lng_min"`
	LngMax          string `json:"lng_max"`
	CenterX         string `json:"center_x"`
	CenterY         string `json:"center_y"`
	Altitude        string `json:"altitude"`
	PendapatanMotor int    `json:"pendapatan_motor"`
	PendapatanMobil int    `json:"pendapatan_mobil"`
	TotalPendapatan int    `json:"total_pendapatan"`
}

type TransaksiItemModel struct {
	NamaJukir  string `json:"nama_jukir"`
	Parlok     string `json:"parlok"`
	Zona       string `json:"zona"`
	Plat       string `json:"plat"`
	Waktu      string `json:"waktu"`
	Kendaraan  string `json:"kendaraan"`
	Pembayaran string `json:"pembayaran"`
	Tarif      int    `json:"tarif"`
}
