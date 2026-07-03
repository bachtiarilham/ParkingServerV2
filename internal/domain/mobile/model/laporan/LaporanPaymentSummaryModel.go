package laporan

type LaporanPaymentSummaryModel struct {
	Label      *string `json:"label,omitempty"`
	Amount     *int64  `json:"amount,omitempty"`
	Percentage *int    `json:"percentage,omitempty"`
}
