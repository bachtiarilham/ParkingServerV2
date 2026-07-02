package usecase

import (
	"context"
	"fmt"

	model "modulegue/internal/domain/mobile/model/riwayat"
	"modulegue/internal/domain/mobile/repository"
	"modulegue/internal/middleware"
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
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok && reqModel.UserID == 0 {
		return nil, fmt.Errorf("user not authenticated")
	}

	roleID, ok := middleware.RoleIDFromContext(ctx)
	if !ok && reqModel.RoleID == 0 {
		return nil, fmt.Errorf("role not found in context")
	}

	if reqModel.UserID == 0 {
		reqModel.UserID = userID
	}
	if reqModel.RoleID == 0 {
		reqModel.RoleID = roleID
	}

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
