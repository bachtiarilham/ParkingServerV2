package monitoring

type MonitoringRequestModel struct {
	TglAwal     string `json:"tgl_awal"`
	TglAkhir    string `json:"tgl_akhir"`
	IDLokasi    int    `json:"idlokasi"`
	NamaPetugas string `json:"nama_petugas"`
}
