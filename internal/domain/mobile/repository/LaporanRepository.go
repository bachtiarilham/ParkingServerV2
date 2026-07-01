package repository

import (
	"context"

	"modulegue/internal/domain/mobile/model"
)

type LaporanRepository interface {
	GetLaporan(ctx context.Context, filter model.LaporanFilterRequestModel) (*model.LaporanModel, error)
}
