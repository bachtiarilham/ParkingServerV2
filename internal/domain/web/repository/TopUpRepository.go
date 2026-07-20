package repository

import (
	"context"
	model "modulegue/internal/domain/web/model/topup"
)

type TopUpRepository interface {
	TopUp(context.Context, model.TopUpRequestModel) (*model.TopUpResponseModel, error)
}
