package laporan

type LaporanRecentTransactionModel struct {
	Code       *string `json:"code,omitempty"`
	Time       *string `json:"time,omitempty"`
	Total      *int64  `json:"total,omitempty"`
	PaymentTag *string `json:"payment_tag,omitempty"`
}
