package riwayat

import (
	"context"
	"fmt"

	model "modulegue/internal/domain/mobile/model/riwayat"
	"modulegue/internal/domain/mobile/repository"
)

type GetParkirDetilUseCase struct {
	riwayatRepo repository.RiwayatRepository
}

func NewGetParkirDetilUseCase(
	riwayatRepo repository.RiwayatRepository,
) *GetParkirDetilUseCase {
	return &GetParkirDetilUseCase{
		riwayatRepo: riwayatRepo,
	}
}

func (uc *GetParkirDetilUseCase) Execute(ctx context.Context, reqModel model.DetilParkirRequestModel) (*model.DetilParkirModel, error) {
	result, err := uc.riwayatRepo.GetParkirDetil(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("get riwayat: %w", err)
	}

	return result, nil
}
