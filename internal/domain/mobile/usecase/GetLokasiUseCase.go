package usecase

import (
	"context"

	"modulegue/internal/domain/mobile/model"
	"modulegue/internal/domain/mobile/repository"
)

type GetLokasiUseCase struct {
	getLokasiRepo repository.GetLokasiRepository
}

func NewGetLokasiUseCase(getLokasiRepo repository.GetLokasiRepository) *GetLokasiUseCase {
	return &GetLokasiUseCase{
		getLokasiRepo: getLokasiRepo,
	}
}

func (uc *GetLokasiUseCase) Execute(ctx context.Context) (*model.LokasiModel, error) {
	return uc.getLokasiRepo.GetLokasi(ctx)
}
