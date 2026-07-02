package riwayat

// RiwayatModel adalah wrapper data riwayat transaksi.
// Detail item dan section disimpan di file terpisah, tapi tetap dalam package yang sama.
type RiwayatModel struct {
	Sections []RiwayatSectionModel `json:"sections,omitempty"`
}
