package repository

import (
	"context"

	"modulegue/internal/domain/mobile/model"
)

type LokasiRepository interface {
	GetLokasi(ctx context.Context) (*model.LokasiModel, error)
}
