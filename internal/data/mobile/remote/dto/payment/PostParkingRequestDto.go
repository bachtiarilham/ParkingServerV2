package payment

type PostParkingRequestDto struct {
	NomorPolisi    int64  `json:"nomor_polisi"`
	JenisKendaraan string `json:"jenis_kendaraan"`
	WaktuMasuk     string `json:"waktu_masuk"`
	ZonaParkir     string `json:"zona_parkir"`
	LokasiParkir   string `json:"lokasi_parkir"`
}
