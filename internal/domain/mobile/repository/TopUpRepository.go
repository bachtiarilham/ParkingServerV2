package repository

import (
	"context"
	model "modulegue/internal/domain/mobile/model/topup"
)

type TopUpRepository interface {
	TopUpCreate(ctx context.Context, reqModel model.TopupCreateRequestModel) (*model.TopupCreateResponseModel, error)
	TopUpStatus(ctx context.Context, req string) (*model.TopupStatusResponseModel, error)
	TopUpCallback(ctx context.Context, reqModel model.QrisCallbackRequestModel) (*model.QrisCallbackResponseModel, error)
}
