package parking

import "time"

type PostParkingResponseModel struct {
	SessionCode string
	PlateNumber string
	Waktu       time.Time
	QrExpired   time.Time
	BiayaParkir int64
}
