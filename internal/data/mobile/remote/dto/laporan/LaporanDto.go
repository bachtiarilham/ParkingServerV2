package laporan

type LaporanFilterRequestDto struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	RoleID    int64  `json:"role_id"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Lokasi    string `json:"lokasi"`
}

type LaporanDto struct {
	TanggalTerpilih    *string                       `json:"tanggal_terpilih,omitempty"`
	Periode            *LaporanDateRangeDto          `json:"periode,omitempty"`
	Summary            *LaporanSummaryDto            `json:"summary,omitempty"`
	ChartBars          []LaporanChartBarDto          `json:"chart_bars,omitempty"`
	PaymentSummaries   []LaporanPaymentSummaryDto    `json:"payment_summaries,omitempty"`
	RecentTransactions []LaporanRecentTransactionDto `json:"recent_transactions,omitempty"`
}

type LaporanDateRangeDto struct {
	StartDate *string `json:"start_date,omitempty"`
	EndDate   *string `json:"end_date,omitempty"`
	Label     *string `json:"label,omitempty"`
}

type LaporanSummaryDto struct {
	TotalTransaksi    *int   `json:"total_transaksi,omitempty"`
	TotalPendapatan   *int64 `json:"total_pendapatan,omitempty"`
	RataRataTransaksi *int64 `json:"rata_rata_transaksi,omitempty"`
}

type LaporanChartBarDto struct {
	Tanggal         *string  `json:"tanggal,omitempty"`
	Amount          *int64   `json:"amount,omitempty"`
	Value           *float64 `json:"value,omitempty"`
	PeriodLabel     *string  `json:"period_label,omitempty"`
	PeriodStartDate *string  `json:"period_start_date,omitempty"`
	PeriodEndDate   *string  `json:"period_end_date,omitempty"`
}

type LaporanPaymentSummaryDto struct {
	Label      *string `json:"label,omitempty"`
	Amount     *int64  `json:"amount,omitempty"`
	Percentage *int    `json:"percentage,omitempty"`
}

type LaporanRecentTransactionDto struct {
	Code       *string `json:"code,omitempty"`
	Time       *string `json:"time,omitempty"`
	Total      *int64  `json:"total,omitempty"`
	PaymentTag *string `json:"payment_tag,omitempty"`
}
