package riwayat

// RiwayatModel adalah wrapper data riwayat transaksi.
// Detail item dan section disimpan di file terpisah, tapi tetap dalam package yang sama.
type RiwayatModel struct {
	ParkirSections []RiwayatSectionModel `json:"sections,omitempty"`
	TopUpSections  []TopUpSectionModel   `json:"topup_sections,omitempty"`
}

type RiwayatSectionModel struct {
	Date  string             `json:"date,omitempty"`
	Items []RiwayatItemModel `json:"items,omitempty"`
}

type RiwayatItemModel struct {
	Code        string `json:"code,omitempty"`
	PlateNumber string `json:"plate_number,omitempty"`
	VehicleType string `json:"vehicle_type,omitempty"`
	Time        string `json:"time,omitempty"`
	Amount      int64  `json:"amount,omitempty"`
	IsEntry     bool   `json:"is_entry,omitempty"`
}

type TopUpSectionModel struct {
	Date  string           `json:"date,omitempty"`
	Items []TopUpItemModel `json:"items,omitempty"`
}

type TopUpItemModel struct {
	Code              string `json:"code,omitempty"`
	PaymentMethodName string `json:"payment_name,omitempty"`
	TransactionStatus string `json:"transaction_status,omitempty"`
	ProviderName      string `json:"provider_name,omitempty"`
	Time              string `json:"time,omitempty"`
	Amount            int64  `json:"amount,omitempty"`
}
