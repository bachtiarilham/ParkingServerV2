package usecase

import (
	"context"

	model "modulegue/internal/domain/mobile/model/helper"
	"modulegue/internal/domain/mobile/repository"
)

type GetTarifUseCase struct {
	getTarifRepo repository.HelperRepository
}

func NewGetTarifUseCase(getTarifRepo repository.HelperRepository) *GetTarifUseCase {
	return &GetTarifUseCase{
		getTarifRepo: getTarifRepo,
	}
}

func (uc *GetTarifUseCase) Execute(ctx context.Context) (*model.TarifModel, error) {
	return uc.getTarifRepo.GetTarif(ctx)
}
