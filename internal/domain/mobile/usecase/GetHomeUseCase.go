package usecase

import (
	"context"

	"modulegue/internal/domain/mobile/model"
	"modulegue/internal/domain/mobile/repository"
)

type GetHomeInput struct {
	UserID int64
	RoleID int64
}

type GetHomeUseCase struct {
	homeRepo repository.HomeRepository
}

func NewGetHomeUseCase(homeRepo repository.HomeRepository) *GetHomeUseCase {
	return &GetHomeUseCase{
		homeRepo: homeRepo,
	}
}

func (uc *GetHomeUseCase) Execute(ctx context.Context, input GetHomeInput) (*model.HomeModel, error) {
	_ = input
	return uc.homeRepo.GetHome(ctx)
}
