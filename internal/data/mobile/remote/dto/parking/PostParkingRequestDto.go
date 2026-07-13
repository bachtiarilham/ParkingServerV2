package parking

type PostParkingRequestDto struct {
	PlateNumber     string `json:"nomor_polisi"`
	VehicleTypeCode string `json:"jenis_kendaraan"`
	SelectedAreaId  int64  `json:"area_parkir"`
}
