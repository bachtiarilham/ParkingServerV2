package helper

import (
	"context"

	model "modulegue/internal/domain/mobile/model/helper"
	"modulegue/internal/domain/mobile/repository"
)

type GetNominalTopUpUseCase struct {
	getNominalTopUpRepo repository.HelperRepository
}

func NewNominalTopUpUseCase(getNominalTopUpRepo repository.HelperRepository) *GetNominalTopUpUseCase {
	return &GetNominalTopUpUseCase{
		getNominalTopUpRepo: getNominalTopUpRepo,
	}
}

func (uc *GetNominalTopUpUseCase) Execute(ctx context.Context) (*model.NominalPaymentModel, error) {
	return uc.getNominalTopUpRepo.GetNominalPayment(ctx)
}
