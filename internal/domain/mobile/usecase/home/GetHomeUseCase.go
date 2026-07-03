package home

import (
	"context"

	model "modulegue/internal/domain/mobile/model/home"
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

func (uc *GetHomeUseCase) Execute(ctx context.Context, reqModel model.GetHomeReqModel) (*model.HomeModel, error) {
	result, err := uc.homeRepo.GetHome(ctx, reqModel)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return &model.HomeModel{
			Events:   []model.EventsModel{},
			News:     []model.NewsModel{},
			Warnings: &model.WarningsModel{},
		}, nil
	}
	if result.Events == nil {
		result.Events = []model.EventsModel{}
	}
	if result.News == nil {
		result.News = []model.NewsModel{}
	}
	if result.Warnings == nil {
		result.Warnings = &model.WarningsModel{}
	}
	return result, nil
}
