package riwayat

type RiwayatItemModel struct {
	Code        *string `json:"code,omitempty"`
	PlateNumber *string `json:"plate_number,omitempty"`
	VehicleType *string `json:"vehicle_type,omitempty"`
	Time        *string `json:"time,omitempty"`
	Amount      *int64  `json:"amount,omitempty"`
	IsEntry     *bool   `json:"is_entry,omitempty"`
}
