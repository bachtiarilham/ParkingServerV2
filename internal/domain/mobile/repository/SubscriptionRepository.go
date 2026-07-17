package repository

import (
	"context"

	model "modulegue/internal/domain/mobile/model/subscription"
)

type SubscriptionRepository interface {
	GetSubscribe(ctx context.Context, userId int64) (*model.SubscriptionResponseModel, error)
}
