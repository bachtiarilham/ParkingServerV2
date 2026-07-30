package repository

import (
	"context"

	req "modulegue/internal/domain/mobile/model/filter_pencarian"
	resplap "modulegue/internal/domain/mobile/model/laporan"
	resp "modulegue/internal/domain/mobile/model/riwayat"
)

type FilterPencarianRepository interface {
	GetRiwayatTransaksi(ctx context.Context, req req.FilterPencarianModel) (*resp.RiwayatTransaksiModel, error)
	GetRiwayatMembership(ctx context.Context, req req.FilterPencarianModel) (*resp.RiwayatMembershipModel, error)
	GetRiwayatParkir(ctx context.Context, req req.FilterPencarianModel) (*resp.RiwayatParkirModel, error)
	GetLaporan(ctx context.Context, req req.FilterPencarianModel) (*resplap.LaporanModel, error)
}
