package model

type LaporanFilterRequestModel struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	RoleID    int64  `json:"role_id"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Lokasi    string `json:"lokasi"`
}

type LaporanModel struct {
	TanggalTerpilih    *string                         `json:"tanggal_terpilih,omitempty"`
	Periode            *LaporanDateRangeModel          `json:"periode,omitempty"`
	Summary            *LaporanSummaryModel            `json:"summary,omitempty"`
	ChartBars          []LaporanChartBarModel          `json:"chart_bars,omitempty"`
	PaymentSummaries   []LaporanPaymentSummaryModel    `json:"payment_summaries,omitempty"`
	RecentTransactions []LaporanRecentTransactionModel `json:"recent_transactions,omitempty"`
}

type LaporanDateRangeModel struct {
	StartDate *string `json:"start_date,omitempty"`
	EndDate   *string `json:"end_date,omitempty"`
	Label     *string `json:"label,omitempty"`
}

type LaporanSummaryModel struct {
	TotalTransaksi    *int   `json:"total_transaksi,omitempty"`
	TotalPendapatan   *int64 `json:"total_pendapatan,omitempty"`
	RataRataTransaksi *int64 `json:"rata_rata_transaksi,omitempty"`
}

type LaporanChartBarModel struct {
	Tanggal         *string  `json:"tanggal,omitempty"`
	Amount          *int64   `json:"amount,omitempty"`
	Value           *float64 `json:"value,omitempty"`
	PeriodLabel     *string  `json:"period_label,omitempty"`
	PeriodStartDate *string  `json:"period_start_date,omitempty"`
	PeriodEndDate   *string  `json:"period_end_date,omitempty"`
}

type LaporanPaymentSummaryModel struct {
	Label      *string `json:"label,omitempty"`
	Amount     *int64  `json:"amount,omitempty"`
	Percentage *int    `json:"percentage,omitempty"`
}

type LaporanRecentTransactionModel struct {
	Code       *string `json:"code,omitempty"`
	Time       *string `json:"time,omitempty"`
	Total      *int64  `json:"total,omitempty"`
	PaymentTag *string `json:"payment_tag,omitempty"`
}
