package parking

import "time"

type PostParkingResponseDto struct {
	SessionCode string    `json:"session_code"`
	PlateNumber string    `json:"plate_nums"`
	Waktu       time.Time `json:"waktu"`
	QrExpired   time.Time `json:"qr_expiredat"`
	BiayaParkir int64     `json:"biaya_parkir"`
}
