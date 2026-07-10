package riwayat

import (
	"context"
	"fmt"

	model "modulegue/internal/domain/mobile/model/riwayat"
	"modulegue/internal/domain/mobile/repository"
)

type GetRiwayatUseCase struct {
	riwayatRepo repository.RiwayatRepository
}

func NewGetRiwayatUseCase(
	riwayatRepo repository.RiwayatRepository,
) *GetRiwayatUseCase {
	return &GetRiwayatUseCase{
		riwayatRepo: riwayatRepo,
	}
}

func (uc *GetRiwayatUseCase) Execute(ctx context.Context, reqModel model.RiwayatRequestModel) (*model.RiwayatModel, error) {
	result, err := uc.riwayatRepo.GetRiwayat(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("get riwayat: %w", err)
	}

	if result == nil {
		return &model.RiwayatModel{
			Sections: []model.RiwayatSectionModel{},
		}, nil
	}

	if result.Sections == nil {
		result.Sections = []model.RiwayatSectionModel{}
	}

	return result, nil
}
