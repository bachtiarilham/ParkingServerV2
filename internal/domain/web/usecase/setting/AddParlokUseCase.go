package setting

import (
	"context"
	"fmt"
	model "modulegue/internal/domain/web/model/setting"
	"modulegue/internal/domain/web/repository"
)

type AddParlokUseCase struct {
	settingRepo repository.SettingRepository
}

func NewAddParlokUseCase(
	settingRepo repository.SettingRepository,
) *AddParlokUseCase {
	return &AddParlokUseCase{
		settingRepo: settingRepo,
	}
}

func (uc *AddParlokUseCase) Execute(ctx context.Context, reqModel model.AddParlokRequestModel) error {
	err := uc.settingRepo.AddParlok(ctx, reqModel)
	if err != nil {
		return fmt.Errorf("gagal menambahkan parkir lokasi : %w", err)
	}
	return nil
}
