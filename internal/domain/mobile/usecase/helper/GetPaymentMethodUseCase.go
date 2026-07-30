package helper

import (
	"context"

	model "modulegue/internal/domain/mobile/model/helper"
	"modulegue/internal/domain/mobile/repository"
)

type GetPaymentMethodUseCase struct {
	getPaymentMethodRepo repository.HelperRepository
}

func NewGetPaymentMethodUseCase(getPaymentMethodRepo repository.HelperRepository) *GetPaymentMethodUseCase {
	return &GetPaymentMethodUseCase{
		getPaymentMethodRepo: getPaymentMethodRepo,
	}
}

func (uc *GetPaymentMethodUseCase) Execute(ctx context.Context) (*model.PaymentMethodModel, error) {
	return uc.getPaymentMethodRepo.GetPaymentMethod(ctx)
}
