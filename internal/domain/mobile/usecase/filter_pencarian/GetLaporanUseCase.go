package filterpencarian

import (
	"context"
	"fmt"

	req "modulegue/internal/domain/mobile/model/filter_pencarian"
	resp "modulegue/internal/domain/mobile/model/laporan"
	"modulegue/internal/domain/mobile/repository"
)

type GetLaporanUseCase struct {
	repo repository.FilterPencarianRepository
}

func NewGetLaporanUseCase(
	repo repository.FilterPencarianRepository,
) *GetLaporanUseCase {
	return &GetLaporanUseCase{
		repo: repo,
	}
}

func (uc *GetLaporanUseCase) Execute(ctx context.Context, reqModel req.FilterPencarianModel) (*resp.LaporanModel, error) {

	result, err := uc.repo.GetLaporan(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("get laporan: %w", err)
	}

	return result, nil
}
