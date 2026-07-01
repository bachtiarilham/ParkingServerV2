package repository

import (
	"context"

	"modulegue/internal/domain/mobile/model"
)

type LoginRepository interface {
	Authenticate(ctx context.Context, username string, password string) (*model.LoginModel, error)
}
