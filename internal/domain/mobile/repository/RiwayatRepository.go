package repository

import (
	"context"

	model "modulegue/internal/domain/mobile/model/riwayat"
)

type RiwayatRepository interface {
	GetRiwayat(ctx context.Context, req model.RiwayatRequestModel) (*model.RiwayatModel, error)
}
