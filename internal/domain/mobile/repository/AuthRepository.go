package repository

import (
	"context"

	model "modulegue/internal/domain/mobile/model/auth"
)

type AuthRepository interface {
	LoginUser(ctx context.Context, reqModel model.LoginRequestModel) (*model.TokenSetModel, *model.LoginRespModel, error)
	ExistsByIdentity(ctx context.Context, nik, email, username, phone string) (*model.UserExistRespModel, error)

	LogoutUser(ctx context.Context, reqModel model.LogoutReqModel) error
	FindByID(ctx context.Context, id int64) (*model.LoginRespModel, error)
	RegisterUser(ctx context.Context, req model.RegisterRequestModel) error
	ChangePasswordUser(ctx context.Context, userID int64, newPasswordHash string) error
}
