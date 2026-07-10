package laporan

import "time"

type LaporanModel struct {
	TanggalAwal          time.Time
	TanggalAkhir         time.Time
	TotalTransaksi       int64
	TotalPendapatanJukir int64
	PendapatanPerTanggal *[]LaporanItem
}

type LaporanItem struct {
	Tanggal              time.Time
	TotalTransaksi       int64
	TotalPendapatanJukir int64
	MotorCount           int64
	CarCount             int64
	QrisCount            int64
	CashCount            int64
}
