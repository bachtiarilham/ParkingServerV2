package home

import (
	"context"

	model "modulegue/internal/domain/mobile/model/home"
	profileModel "modulegue/internal/domain/mobile/model/profile"
	"modulegue/internal/domain/mobile/repository"
)

type GetHomeUseCase struct {
	homeRepo repository.HomeRepository
}

func NewGetHomeUseCase(
	homeRepo repository.HomeRepository,
) *GetHomeUseCase {
	return &GetHomeUseCase{
		homeRepo: homeRepo,
	}
}

func (uc *GetHomeUseCase) ExecuteCustomerHome(ctx context.Context, reqModel model.GetHomeReqModel) (*model.CustomerHomeModel, error) {
	result, err := uc.homeRepo.GetCustomerHome(ctx, reqModel)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return &model.CustomerHomeModel{
			Profile:          &profileModel.CustomerModel{},
			Contents:         &[]model.ContentsModel{},
			UnreadNotifCount: 0,
		}, nil
	}
	if result.Profile == nil {
		result.Profile = &profileModel.CustomerModel{}
	}
	if result.Contents == nil {
		result.Contents = &[]model.ContentsModel{}
	}
	return result, nil
}

func (uc *GetHomeUseCase) ExecuteJukirHome(ctx context.Context, reqModel model.GetHomeReqModel) (*model.JukirHomeModel, error) {
	result, err := uc.homeRepo.GetJukirHome(ctx, reqModel)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return &model.JukirHomeModel{
			Profile:          &profileModel.JukirModel{},
			Contents:         &[]model.ContentsModel{},
			UnreadNotifCount: 0,
		}, nil
	}
	if result.Profile == nil {
		result.Profile = &profileModel.JukirModel{}
	}
	if result.Contents == nil {
		result.Contents = &[]model.ContentsModel{}
	}
	return result, nil
}
