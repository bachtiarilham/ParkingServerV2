package laporan

type LaporanSummaryModel struct {
	TotalTransaksi    *int   `json:"total_transaksi,omitempty"`
	TotalPendapatan   *int64 `json:"total_pendapatan,omitempty"`
	RataRataTransaksi *int64 `json:"rata_rata_transaksi,omitempty"`
}
