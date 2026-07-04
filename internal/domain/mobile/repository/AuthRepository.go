package repository

import (
	"context"

	model "modulegue/internal/domain/mobile/model/auth"
)

type AuthRepository interface {
	LoginUser(ctx context.Context, identifier, password string) (*model.UserModel, error)
	ExistsByEmailOrUsernameOrPhone(ctx context.Context, email, username, phone string) (*model.UserExistResult, error)

	LogoutUser(ctx context.Context, reqModel model.LogoutReqModel) error
	FindByID(ctx context.Context, id int64) (*model.UserModel, error)
	RegisterUser(ctx context.Context, req model.RegisterRequestModel) (*model.RegisterResponseModel, error)
	ChangePasswordUser(ctx context.Context, userID int64, newPasswordHash string) error
}
