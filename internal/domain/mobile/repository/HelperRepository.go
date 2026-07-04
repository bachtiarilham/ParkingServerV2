package repository

import (
	"context"

	model "modulegue/internal/domain/mobile/model/helper"
)

type HelperRepository interface {
	GetLokasi(ctx context.Context) (*model.LokasiModel, error)
}
