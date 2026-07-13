package payment

import (
	"time"
)

type PostPaymentParkingResponseDto struct {
	SessionId         int64     `json:"session_id"`
	SessionCode       string    `json:"session_code"`
	TransactionCode   string    `json:"trx_code"`
	PlateNumber       string    `json:"plate_number"`
	VehicleTypeCode   string    `json:"vhc_type_code"`
	VehicleTypeName   string    `json:"vhc_type_name"`
	LocationId        int64     `json:"loc_id"`
	LocationName      string    `json:"loc_name"`
	AreaId            int64     `json:"area_id"`
	AreaName          string    `json:"area_name"`
	Amount            int64     `json:"amount"`
	ParkingStatusCode string    `json:"parking_stat_code"`
	ParkingStatusName string    `json:"parking_stat_name"`
	PaymentStatusCode string    `json:"payment_stat_code"`
	PaymentStatusName string    `json:"payment_stat_name"`
	PaymentCode       string    `json:"payment_code"`
	FailedReason      string    `json:"failed_reason"`
	ReceiptNumber     int64     `json:"receipt_number"`
	StartedAt         time.Time `json:"startedat"`
	PaidAt            time.Time `json:"paidat"`
	QrExpiredAt       time.Time `json:"qr_expiredat"`
}
