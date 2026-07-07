package riwayat

type RiwayatRequestModel struct {
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
	RoleID      int64  `json:"role_id"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Payment     string `json:"payment"`
	Vehicle     string `json:"vehicle"`
	Lokasi      string `json:"lokasi"`
}
