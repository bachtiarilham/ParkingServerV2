package riwayat

import (
	"context"
	"fmt"

	model "modulegue/internal/domain/mobile/model/riwayat"
	"modulegue/internal/domain/mobile/repository"
)

type GetTransaksiDetilUseCase struct {
	riwayatRepo repository.RiwayatRepository
}

func NewGetTransaksiDetilUseCase(
	riwayatRepo repository.RiwayatRepository,
) *GetTransaksiDetilUseCase {
	return &GetTransaksiDetilUseCase{
		riwayatRepo: riwayatRepo,
	}
}

func (uc *GetTransaksiDetilUseCase) Execute(ctx context.Context, reqModel model.DetilTransaksiRequestModel) (*model.DetilTransaksiModel, error) {
	result, err := uc.riwayatRepo.GetTransaksiDetil(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("get riwayat: %w", err)
	}

	return result, nil
}
