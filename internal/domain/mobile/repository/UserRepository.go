package repository

import (
	"context"

	"modulegue/internal/domain/mobile/model"
)

type UserRepository interface {
	GetUser(ctx context.Context, userID int64) (*model.LoginModel, error)
	UpdatePassword(ctx context.Context, userID int64, newPasswordHash string) error
	Create(ctx context.Context, user *model.UserModel) error
	Update(ctx context.Context, user *model.UserModel) error
}
