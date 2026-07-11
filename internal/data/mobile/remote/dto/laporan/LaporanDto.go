package laporan

import "time"

type LaporanFilterRequestDto struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type LaporanResponseDto struct {
	TanggalAwal          time.Time      `json:"tanggal_awal"`
	TanggalAkhir         time.Time      `json:"tanggal_akhir"`
	TotalTransaksi       int64          `json:"total_transaksi"`
	TotalPendapatanJukir int64          `json:"total_pendapatan_jukir"`
	PendapatanPerTanggal *[]LaporanItem `json:"pendapatan_per_tanggal"`
}
type LaporanItem struct {
	Tanggal              time.Time `json:"tanggal"`
	TotalTransaksi       int64     `json:"total_transaksi"`
	TotalPendapatanJukir int64     `json:"total_pendapatan_jukir"`
	MotorCount           int64     `json:"motor_count"`
	CarCount             int64     `json:"car_count"`
	QrisCount            int64     `json:"qris_count"`
	CashCount            int64     `json:"cash_count"`
}
