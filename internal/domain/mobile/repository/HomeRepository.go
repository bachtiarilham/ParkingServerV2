package repository

import (
	"context"
	model "modulegue/internal/domain/mobile/model/home"
)

type HomeRepository interface {
	GetJukirHome(ctx context.Context, reqModel model.GetHomeReqModel) (*model.JukirHomeModel, error)
	GetCustomerHome(ctx context.Context, reqModel model.GetHomeReqModel) (*model.CustomerHomeModel, error)
}
