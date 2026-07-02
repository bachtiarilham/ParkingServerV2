package laporan

type LaporanChartBarModel struct {
	ChartBars []LaporanChartBarItemModel `json:"chart_bars,omitempty"`
}

type LaporanChartBarItemModel struct {
	Tanggal         *string  `json:"tanggal,omitempty"`
	Amount          *int64   `json:"amount,omitempty"`
	Value           *float64 `json:"value,omitempty"`
	PeriodLabel     *string  `json:"period_label,omitempty"`
	PeriodStartDate *string  `json:"period_start_date,omitempty"`
	PeriodEndDate   *string  `json:"period_end_date,omitempty"`
}
