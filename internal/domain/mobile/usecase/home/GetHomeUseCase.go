package home

import (
	"context"
	"fmt"

	model "modulegue/internal/domain/mobile/model/home"
	"modulegue/internal/domain/mobile/repository"
	"modulegue/internal/middleware"
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
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("user not authenticated")
	}

	roleID, ok := middleware.RoleIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("role not found in context")
	}

	if reqModel.UserID == 0 {
		reqModel.UserID = userID
	}
	if reqModel.RoleID == 0 {
		reqModel.RoleID = roleID
	}

	result, err := uc.homeRepo.GetHome(ctx, userID, roleID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return &model.HomeModel{
			Events:   &model.EventsModel{Events: []model.EventsItemModel{}},
			News:     &model.NewsModel{News: []model.NewsModel{}},
			Warnings: &model.WarningsModel{},
		}, nil
	}
	if result.Events == nil {
		result.Events = &model.EventsModel{Events: []model.EventsItemModel{}}
	}
	if result.News == nil {
		result.News = &model.NewsModel{News: []model.NewsModel{}}
	}
	if result.Warnings == nil {
		result.Warnings = &model.WarningsModel{}
	}
	return result, nil
}
