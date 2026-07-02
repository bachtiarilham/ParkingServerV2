package laporan

type LaporanRequestModel struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	RoleID    int64  `json:"role_id"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Lokasi    string `json:"lokasi"`
}
