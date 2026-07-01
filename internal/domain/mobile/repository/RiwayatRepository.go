package repository

import (
	"context"

	"modulegue/internal/domain/mobile/model"
)

type RiwayatRepository interface {
	GetRiwayat(ctx context.Context, req model.RiwayatRequestModel) (*model.RiwayatModel, error)
}
