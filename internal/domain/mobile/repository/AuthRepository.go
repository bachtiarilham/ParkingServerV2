package repository

import (
	"context"

	model "modulegue/internal/domain/mobile/model/auth"
)

type AuthRepository interface {
	LoginUser(ctx context.Context, identifier, password string) (*model.UserModel, error)
	ExistsByEmailOrUsernameOrPhone(ctx context.Context, email, username, phone string) (*model.UserExistResult, error)

	// LoginUser(ctx context.Context, email, username, phone, password string) (*model.UserModel, error)
	LogoutUser(ctx context.Context, reqModel model.LogoutReqModel) error
	FindByID(ctx context.Context, id int64) (*model.UserModel, error)
	RegisterUser(ctx context.Context, req model.RegisterRequestModel) (*model.RegisterResponseModel, error)
	ChangePasswordUser(ctx context.Context, userID int64, newPasswordHash string) error
	CreateUser(ctx context.Context, user *model.UserModel) (*model.UserModel, error)
	UpdateUser(ctx context.Context, user *model.UserModel) error

	SaveSession(ctx context.Context, s model.SessionModel) error
	FindSessionByRefreshToken(ctx context.Context, token string) (model.SessionModel, error)
	RotateSession(ctx context.Context, oldRefreshToken string, newSession model.SessionModel) error
	DeleteSession(ctx context.Context, refreshToken string) error
	DeleteAllSessions(ctx context.Context, userID int64) error
}
