package usecase

import (
	"context"
	"fmt"

	model "modulegue/internal/domain/mobile/model/subscription"
	"modulegue/internal/domain/mobile/repository"
	"modulegue/internal/middleware"
)

type SubscriptionUseCase struct {
	subsRepo repository.SubscriptionRepository
}

func NewSubscriptionUseCase(
	subsRepo repository.SubscriptionRepository,
) *SubscriptionUseCase {
	return &SubscriptionUseCase{
		subsRepo: subsRepo,
	}
}

func (uc *SubscriptionUseCase) Execute(ctx context.Context) (*model.SubscribeModel, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("user not authenticated")
	}

	result, err := uc.subsRepo.GetSubscribe(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get subscription: %w", err)
	}

	if result == nil {
		return &model.SubscribeModel{
			PackageCard: []model.PackageCardModel{},
			Promo:       []model.PromoModel{},
		}, nil
	}

	if result.PackageCard == nil {
		result.PackageCard = []model.PackageCardModel{}
	}
	if result.Promo == nil {
		result.Promo = []model.PromoModel{}
	}

	return result, nil
}
