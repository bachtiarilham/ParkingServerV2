package repository

import (
	"context"

	"modulegue/internal/domain/mobile/model"
)

type SubscribeRepository interface {
	GetSubscribe(ctx context.Context) (*model.SubscribeModel, error)
}
