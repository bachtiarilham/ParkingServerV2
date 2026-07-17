package topup

import (
	"context"
	"fmt"

	model "modulegue/internal/domain/mobile/model/topup"
	"modulegue/internal/domain/mobile/repository"
)

type GetTopUpStatusUseCase struct {
	topupGetTopUpStatusRepo repository.TopUpRepository
}

func NewGetTopUpStatusUseCase(
	topupGetTopUpStatusRepo repository.TopUpRepository,
) *GetTopUpStatusUseCase {
	return &GetTopUpStatusUseCase{
		topupGetTopUpStatusRepo: topupGetTopUpStatusRepo,
	}
}

func (uc *GetTopUpStatusUseCase) Execute(ctx context.Context, req string) (*model.TopupStatusResponseModel, error) {
	result, err := uc.topupGetTopUpStatusRepo.TopUpStatus(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get subscription: %w", err)
	}
	return result, nil
}
