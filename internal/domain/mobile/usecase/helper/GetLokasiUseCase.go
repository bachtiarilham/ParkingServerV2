package helper

import (
	"context"

	model "modulegue/internal/domain/mobile/model/helper"
	"modulegue/internal/domain/mobile/repository"
)

type GetLokasiUseCase struct {
	getLokasiRepo repository.HelperRepository
}

func NewGetLokasiUseCase(getLokasiRepo repository.HelperRepository) *GetLokasiUseCase {
	return &GetLokasiUseCase{
		getLokasiRepo: getLokasiRepo,
	}
}

func (uc *GetLokasiUseCase) Execute(ctx context.Context, userId int64) (*model.LokasiModel, error) {
	return uc.getLokasiRepo.GetLokasi(ctx, userId)
}
