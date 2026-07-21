package riwayat

type RiwayatRequestDto struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type RiwayatResponseDto struct {
	ParkingSections []ParkingSectionDto `json:"parking_sections"`
	TopUpSections   []TopUpSectionDto   `json:"topup_sections"`
}

type ParkingSectionDto struct {
	Date  string           `json:"date,omitempty"`
	Items []ParkingItemDto `json:"items,omitempty"`
}

type ParkingItemDto struct {
	Code        string `json:"code,omitempty"`
	PlateNumber string `json:"plate_number,omitempty"`
	VehicleType string `json:"vehicle_type,omitempty"`
	Time        string `json:"time,omitempty"`
	Amount      int64  `json:"amount,omitempty"`
	IsEntry     bool   `json:"is_entry,omitempty"`
}

type TopUpSectionDto struct {
	Date  string         `json:"date,omitempty"`
	Items []TopUpItemDto `json:"items,omitempty"`
}

type TopUpItemDto struct {
	Code              string `json:"code,omitempty"`
	PaymentMethodName string `json:"payment_name,omitempty"`
	TransactionStatus string `json:"transaction_status,omitempty"`
	ProviderName      string `json:"provider_name,omitempty"`
	Time              string `json:"time,omitempty"`
	Amount            int64  `json:"amount,omitempty"`
}
