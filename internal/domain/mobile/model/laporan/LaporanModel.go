package laporan

type LaporanModel struct {
	TanggalTerpilih    *string                        `json:"tanggal_terpilih,omitempty"`
	Periode            *LaporanDateRangeModel         `json:"periode,omitempty"`
	Summary            *LaporanSummaryModel           `json:"summary,omitempty"`
	ChartBars          *LaporanChartBarModel          `json:"chart_bars,omitempty"`
	PaymentSummaries   *LaporanPaymentSummaryModel    `json:"payment_summaries,omitempty"`
	RecentTransactions *LaporanRecentTransactionModel `json:"recent_transactions,omitempty"`
}
