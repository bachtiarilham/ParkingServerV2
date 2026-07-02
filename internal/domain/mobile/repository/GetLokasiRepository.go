package repository

import (
	"context"

	"modulegue/internal/domain/mobile/model"
)

type GetLokasiRepository interface {
	GetLokasi(ctx context.Context) (*model.LokasiModel, error)
}
