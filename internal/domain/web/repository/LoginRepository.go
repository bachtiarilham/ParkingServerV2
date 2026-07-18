package repository

import (
	"context"
	model "modulegue/internal/domain/shared/model/auth"
	modelauth "modulegue/internal/domain/shared/model/auth"
	respmodel "modulegue/internal/domain/web/model/home"
)

type LoginRepository interface {
	Login(ctx context.Context, reqModel model.LoginRequestModel) (*modelauth.TokenSetModel, *respmodel.HomeResponseModel, error)
}
