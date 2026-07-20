package helper

import (
	"context"
	"fmt"
	model "modulegue/internal/domain/web/model/helper"
	"modulegue/internal/domain/web/repository"
)

type GetLokasiUseCase struct {
	helperRepo repository.HelperRepository
}

func NewGetLokasiUseCase(
	helperRepo repository.HelperRepository,
) *GetLokasiUseCase {
	return &GetLokasiUseCase{
		helperRepo: helperRepo,
	}
}

func (uc *GetLokasiUseCase) Execute(ctx context.Context, reqModel model.GetLokasiRequestModel) (*model.GetLokasiResponseModel, error) {
	resp, err := uc.helperRepo.GetLokasi(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data monitoring : %w", err)
	}
	return resp, nil
}
