package repository

import (
	"context"

	model "modulegue/internal/domain/mobile/model/subscription"
)

type SubscriptionRepository interface {
	GetSubscribe(ctx context.Context) (*model.SubscribeModel, error)
}
