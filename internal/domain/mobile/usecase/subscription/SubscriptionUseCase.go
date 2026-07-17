package subscription

import (
	"context"
	"fmt"

	model "modulegue/internal/domain/mobile/model/subscription"
	"modulegue/internal/domain/mobile/repository"
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

func (uc *SubscriptionUseCase) Execute(ctx context.Context, userId int64) (*model.SubscriptionResponseModel, error) {
	result, err := uc.subsRepo.GetSubscribe(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("get subscription: %w", err)
	}

	if result == nil {
		return emptySubscriptionResponse(), nil
	}

	if result.ActivePackageBenefit == nil {
		result.ActivePackageBenefit = []string{}
	}
	if result.ListPaket.Bulanan == nil {
		result.ListPaket.Bulanan = []model.DetailPaket{}
	}
	if result.ListPaket.EnamBulan == nil {
		result.ListPaket.EnamBulan = []model.DetailPaket{}
	}
	if result.ListPaket.Tahunan == nil {
		result.ListPaket.Tahunan = []model.DetailPaket{}
	}
	if result.PromoTersedia.SyaratDanKetentuan == nil {
		result.PromoTersedia.SyaratDanKetentuan = []string{}
	}
	if result.PromoTersedia.EachPromo == nil {
		result.PromoTersedia.EachPromo = []model.DetailPromo{}
	}

	return result, nil
}

func emptySubscriptionResponse() *model.SubscriptionResponseModel {
	return &model.SubscriptionResponseModel{
		ActivePackageBenefit: []string{},
		ListPaket: model.ListPaket{
			Bulanan:   []model.DetailPaket{},
			EnamBulan: []model.DetailPaket{},
			Tahunan:   []model.DetailPaket{},
		},
		PromoTersedia: model.PromoTersedia{
			SyaratDanKetentuan: []string{},
			EachPromo:          []model.DetailPromo{},
		},
	}
}
