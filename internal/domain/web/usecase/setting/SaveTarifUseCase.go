package setting

import (
	"context"
	"fmt"
	model "modulegue/internal/domain/web/model/setting"
	"modulegue/internal/domain/web/repository"
)

type SaveTarifUseCase struct {
	settingRepo repository.SettingRepository
}

func NewSaveTarifUseCase(
	settingRepo repository.SettingRepository,
) *SaveTarifUseCase {
	return &SaveTarifUseCase{
		settingRepo: settingRepo,
	}
}

func (uc *SaveTarifUseCase) Execute(ctx context.Context, reqModel model.SaveTarifRequestModel) error {
	err := uc.settingRepo.SaveTarif(ctx, reqModel)
	if err != nil {
		return fmt.Errorf("gagal menyimpan tarif : %w", err)
	}
	return nil
}
