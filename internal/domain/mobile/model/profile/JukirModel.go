package profile

import "time"

type JukirModel struct {
	UserId                  int64     `json:"user_id"`
	Nik                     string    `json:"nik"`
	FullName                string    `json:"full_name"`
	Username                string    `json:"username"`
	Email                   string    `json:"email"`
	Phone                   string    `json:"phone"`
	PhotoUrl                string    `json:"photo_url"`
	IsVerified              bool      `json:"is_verified"`
	RoleId                  int64     `json:"role_id"`
	RoleCode                string    `json:"role_code"`
	RoleName                string    `json:"role_name"`
	Saldo                   int64     `json:"saldo"`
	LocationId              int64     `json:"location_id"`
	LocationCode            string    `json:"location_code"`
	LocationName            string    `json:"location_name"`
	Address                 string    `json:"address"`
	MinLatitude             float64   `json:"minLatitude"`
	MaxLatitude             float64   `json:"maxLatitude"`
	MinLongitude            float64   `json:"minLongitude"`
	MaxLongitude            float64   `json:"maxLongitude"`
	CenterLatitude          float64   `json:"centerLatitude"`
	CenterLongitude         float64   `json:"centerLongitude"`
	RadiusMeter             int64     `json:"radius_meter"`
	AreaId                  int64     `json:"area_id"`
	AreaName                string    `json:"area_name"`
	ZoneId                  int64     `json:"zone_id"`
	ZoneName                string    `json:"zone_name"`
	AssignmentEffectiveFrom time.Time `json:"assignment_effective_from"`
	AssignmentEffectiveTo   time.Time `json:"assignment_effective_to"`
	TodayIncome             int64     `json:"today_income"`
	TotalIncome             int64     `json:"total_income"`
	TodayTransactionCount   int64     `json:"today_transaction_count"`
	UnreadNotificationCount int64     `json:"unread_notification_count"`
}
