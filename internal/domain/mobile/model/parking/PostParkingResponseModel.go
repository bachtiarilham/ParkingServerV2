package parking

import "time"

type PostParkingResponseModel struct {
	SessionId         int64
	SessionCode       string
	TransactionCode   string
	PlateNumber       string
	VehicleTypeCode   string
	VehicleTypeName   string
	ZoneId            int64
	ZoneName          string
	LocationId        int64
	LocationName      string
	Address           string
	AreaId            int64
	AreaName          string
	Amount            int64
	QrString          string
	QrExpiredAt       time.Time
	PaymentCode       string
	PaymentStatusCode string
	PaymentStatusName string
}
