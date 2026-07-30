package filterpencarian

import (
	"context"
	"fmt"

	req "modulegue/internal/domain/mobile/model/filter_pencarian"
	resp "modulegue/internal/domain/mobile/model/riwayat"

	"modulegue/internal/domain/mobile/repository"
)

type GetRiwayatMembershipUseCase struct {
	repo repository.FilterPencarianRepository
}

func NewGetRiwayatMembershipUseCase(
	repo repository.FilterPencarianRepository,
) *GetRiwayatMembershipUseCase {
	return &GetRiwayatMembershipUseCase{
		repo: repo,
	}
}

func (uc *GetRiwayatMembershipUseCase) Execute(ctx context.Context, reqModel req.FilterPencarianModel) (*resp.RiwayatMembershipModel, error) {

	result, err := uc.repo.GetRiwayatMembership(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("get riwayat: %w", err)
	}

	return result, nil
}
