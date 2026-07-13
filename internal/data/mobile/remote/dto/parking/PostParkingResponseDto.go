package parking

import "time"

type PostParkingResponseDto struct {
	SessionId         int64     `json:"session_id"`
	SessionCode       string    `json:"session_code"`
	TransactionCode   string    `json:"trx_code"`
	PlateNumber       string    `json:"plate_nums"`
	VehicleTypeCode   string    `json:"vhc_type_code"`
	VehicleTypeName   string    `json:"vhc_type_name"`
	ZoneId            int64     `json:"zone_id"`
	ZoneName          string    `json:"zone_name"`
	LocationId        int64     `json:"loc_id"`
	LocationName      string    `json:"loc_name"`
	Address           string    `json:"address"`
	AreaId            int64     `json:"adea_id"`
	AreaName          string    `json:"area_name"`
	Amount            int64     `json:"amount"`
	QrString          string    `json:"qr_str"`
	QrExpiredAt       time.Time `json:"qr_expiredat"`
	PaymentCode       string    `json:"payment_code"`
	PaymentStatusCode string    `json:"payment_stat_code"`
	PaymentStatusName string    `json:"payment_stat_name"`
}
