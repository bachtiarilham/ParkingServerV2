package laporan

import "time"

type LaporanRequestModel struct {
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	RoleID    int64     `json:"role_id"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Lokasi    string    `json:"lokasi"`
}
