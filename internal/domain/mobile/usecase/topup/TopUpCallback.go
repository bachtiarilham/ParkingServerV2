package topup

import (
	"context"
	"fmt"

	model "modulegue/internal/domain/mobile/model/topup"
	"modulegue/internal/domain/mobile/repository"
)

type TopUpCallbackUseCase struct {
	topupCallbackRepo repository.TopUpRepository
}

func NewTopUpCallbackUseCase(
	topupCallbackRepo repository.TopUpRepository,
) *TopUpCallbackUseCase {
	return &TopUpCallbackUseCase{
		topupCallbackRepo: topupCallbackRepo,
	}
}

func (uc *TopUpCallbackUseCase) Execute(ctx context.Context, reqModel model.QrisCallbackRequestModel) (*model.QrisCallbackResponseModel, error) {
	result, err := uc.topupCallbackRepo.TopUpCallback(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("topup callback: %w", err)
	}
	return result, nil
}
