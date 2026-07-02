package laporan

type LaporanDateRangeModel struct {
	StartDate *string `json:"start_date,omitempty"`
	EndDate   *string `json:"end_date,omitempty"`
	Label     *string `json:"label,omitempty"`
}
