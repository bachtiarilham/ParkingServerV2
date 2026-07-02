package repository

import (
	"context"

	model "modulegue/internal/domain/mobile/model/laporan"
)

type LaporanRepository interface {
	GetLaporan(ctx context.Context, filter model.LaporanRequestModel) (*model.LaporanModel, error)
}
