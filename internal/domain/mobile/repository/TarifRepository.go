package repository

import (
	"context"

	"modulegue/internal/domain/mobile/model"
)

type TarifRepository interface {
	GetTarif(ctx context.Context) ([]model.TarifModel, error)
}
