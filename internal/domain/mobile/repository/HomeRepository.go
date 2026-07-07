package repository

import (
	"context"
	model "modulegue/internal/domain/mobile/model/home"
)

type HomeRepository interface {
	GetHome(ctx context.Context, reqModel model.GetHomeReqModel) (*model.HomeModel, error)
}
