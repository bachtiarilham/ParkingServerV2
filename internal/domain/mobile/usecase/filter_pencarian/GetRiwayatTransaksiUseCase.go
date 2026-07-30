package filterpencarian

import (
	"context"
	"fmt"

	req "modulegue/internal/domain/mobile/model/filter_pencarian"
	resp "modulegue/internal/domain/mobile/model/riwayat"
	"modulegue/internal/domain/mobile/repository"
)

type GetRiwayatTransaksiUseCase struct {
	repo repository.FilterPencarianRepository
}

func NewGetRiwayatTransaksiUseCase(
	repo repository.FilterPencarianRepository,
) *GetRiwayatTransaksiUseCase {
	return &GetRiwayatTransaksiUseCase{
		repo: repo,
	}
}

func (uc *GetRiwayatTransaksiUseCase) Execute(ctx context.Context, reqModel req.FilterPencarianModel) (*resp.RiwayatTransaksiModel, error) {

	result, err := uc.repo.GetRiwayatTransaksi(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("get riwayat: %w", err)
	}

	return result, nil
}
