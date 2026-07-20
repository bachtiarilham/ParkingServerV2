package helper

import (
	"context"
	"fmt"
	model "modulegue/internal/domain/web/model/helper"
	"modulegue/internal/domain/web/repository"
)

type GetTarifUseCase struct {
	helperRepo repository.HelperRepository
}

func NewGetTarifUseCase(
	helperRepo repository.HelperRepository,
) *GetTarifUseCase {
	return &GetTarifUseCase{
		helperRepo: helperRepo,
	}
}

func (uc *GetTarifUseCase) Execute(ctx context.Context, reqModel model.GetTarifRequestModel) (*model.GetTarifResponseModel, error) {
	resp, err := uc.helperRepo.GetTarif(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data monitoring : %w", err)
	}
	return resp, nil
}
