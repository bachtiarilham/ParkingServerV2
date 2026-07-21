package riwayat

import "time"

type RiwayatRequestModel struct {
	UserID    int64     `json:"user_id"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}
