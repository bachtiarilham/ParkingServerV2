package repository

import (
	"context"

	"modulegue/internal/domain/mobile/model"
)

type RegisterRepository interface {
	Register(ctx context.Context, req model.RegisterRequestModel) (*model.RegisterResponseModel, error)
}
