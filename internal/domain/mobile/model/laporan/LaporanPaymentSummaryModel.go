package laporan

type LaporanPaymentSummaryModel struct {
	PaymentSummaries []LaporanPaymentSummaryItemModel `json:"payment_summaries,omitempty"`
}

type LaporanPaymentSummaryItemModel struct {
	Label      *string `json:"label,omitempty"`
	Amount     *int64  `json:"amount,omitempty"`
	Percentage *int    `json:"percentage,omitempty"`
}
