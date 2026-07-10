package riwayat

type RiwayatRequestModel struct {
	UserID      int64  `json:"user_id"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	PaymentCode string `json:"payment"`
	VehicleCode string `json:"vehicle"`
	LokasiCode  string `json:"lokasi"`
}
