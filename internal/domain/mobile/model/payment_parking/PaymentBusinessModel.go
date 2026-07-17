package payment

import "time"

type PaymentBusinessModel struct {
	SessionId       int64
	SessionCode     string
	TransactionCode string
	CustomerUserId  int64
	OfficerUserId   int64
	VehicleTypeId   int64
	VehicleTypeCode string
	VehicleTypeName string
	PlateNumber     string
	LocationId      int64
	LocationName    string
	Address         string
	AreaId          int64
	AreaName        string
	ZoneId          int64
	ZoneName        string
	Amount          int64
	QrExpiredAt     time.Time
	PaymentCode     string
	ParkingStatusId string
	PaymentStatusId string
	FailedReason    string
}
