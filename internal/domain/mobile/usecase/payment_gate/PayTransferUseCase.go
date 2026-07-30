package paymentgate

import (
	"context"
	"fmt"

	model "modulegue/internal/domain/mobile/model/payment_gate"
	"modulegue/internal/domain/mobile/repository"
)

type PayTransferUseCase struct {
	repo repository.PaymentGateRepository
}

func NewPayTransferUseCase(
	repo repository.PaymentGateRepository,
) *PayTransferUseCase {
	return &PayTransferUseCase{
		repo: repo,
	}
}

func (uc *PayTransferUseCase) Execute(ctx context.Context, reqmodel model.PayRequestModel) (*model.PayResponseModel, error) {
	result, err := uc.repo.PayTransfer(ctx, reqmodel)
	if err != nil {
		return nil, fmt.Errorf("pay transfer repo: %w", err)
	}
	if result == nil {
		return nil, fmt.Errorf("pay transfer returned empty result")
	}

	return result, nil
}
