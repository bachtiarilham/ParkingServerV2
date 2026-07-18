package repository

import (
	"context"
	model "modulegue/internal/domain/web/model/topup"
)

type TopUpRepository interface {
	AddParlok(context.Context, model.TopUpRequestModel) error
}
