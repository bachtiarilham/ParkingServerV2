package dto

type RiwayatRequestDto struct {
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
	RoleID      int64  `json:"role_id"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Transaction string `json:"transaction"`
	Payment     string `json:"payment"`
	Vehicle     string `json:"vehicle"`
	Lokasi      string `json:"lokasi"`
}

type RiwayatDto struct {
	Sections []RiwayatSectionDto `json:"sections,omitempty"`
}

type RiwayatSectionDto struct {
	Date  *string          `json:"date,omitempty"`
	Items []RiwayatItemDto `json:"items,omitempty"`
}

type RiwayatItemDto struct {
	Code        *string `json:"code,omitempty"`
	PlateNumber *string `json:"plate_number,omitempty"`
	VehicleType *string `json:"vehicle_type,omitempty"`
	Time        *string `json:"time,omitempty"`
	Amount      *int64  `json:"amount,omitempty"`
	IsEntry     *bool   `json:"is_entry,omitempty"`
}
