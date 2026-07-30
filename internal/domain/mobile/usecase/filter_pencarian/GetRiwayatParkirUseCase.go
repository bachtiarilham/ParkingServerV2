package filterpencarian

import (
	"context"
	"fmt"

	req "modulegue/internal/domain/mobile/model/filter_pencarian"
	resp "modulegue/internal/domain/mobile/model/riwayat"
	"modulegue/internal/domain/mobile/repository"
)

type GetRiwayatParkirUseCase struct {
	repo repository.FilterPencarianRepository
}

func NewGetRiwayatParkirUseCase(
	repo repository.FilterPencarianRepository,
) *GetRiwayatParkirUseCase {
	return &GetRiwayatParkirUseCase{
		repo: repo,
	}
}

func (uc *GetRiwayatParkirUseCase) Execute(ctx context.Context, reqModel req.FilterPencarianModel) (*resp.RiwayatParkirModel, error) {

	result, err := uc.repo.GetRiwayatParkir(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("get riwayat: %w", err)
	}

	return result, nil
}
