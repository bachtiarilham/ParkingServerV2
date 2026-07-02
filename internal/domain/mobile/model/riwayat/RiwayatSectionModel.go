package riwayat

type RiwayatSectionModel struct {
	Date  *string            `json:"date,omitempty"`
	Items []RiwayatItemModel `json:"items,omitempty"`
}
