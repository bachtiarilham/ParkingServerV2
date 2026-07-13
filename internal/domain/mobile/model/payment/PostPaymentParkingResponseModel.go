package payment

import "time"

type PostPaymentParkingResponseModel struct {
	SessionId         int64
	SessionCode       string
	TransactionCode   string
	PlateNumber       string
	VehicleTypeCode   string
	VehicleTypeName   string
	LocationId        int64
	LocationName      string
	AreaId            int64
	AreaName          string
	Amount            int64
	ParkingStatusCode string
	ParkingStatusName string
	PaymentStatusCode string
	PaymentStatusName string
	PaymentCode       string
	FailedReason      string
	ReceiptNumber     int64
	StartedAt         time.Time
	PaidAt            time.Time
	QrExpiredAt       time.Time
}
