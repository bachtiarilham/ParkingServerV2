package riwayat

type RiwayatRequestDto struct {
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	PaymentCode string `json:"paymentCode"`
	VehicleCode string `json:"vehicleCode"`
	LokasiCode  string `json:"lokasiCode"`
}

type RiwayatResponseDto struct {
	Sections []RiwayatSectionDto `json:"sections"`
}

type RiwayatSectionDto struct {
	Date  string           `json:"date,omitempty"`
	Items []RiwayatItemDto `json:"items,omitempty"`
}

type RiwayatItemDto struct {
	Code        string `json:"code,omitempty"`
	PlateNumber string `json:"plate_number,omitempty"`
	VehicleType string `json:"vehicle_type,omitempty"`
	Time        string `json:"time,omitempty"`
	Amount      int64  `json:"amount,omitempty"`
	IsEntry     bool   `json:"is_entry,omitempty"`
}
