package setting

import (
	"context"
	"fmt"
	model "modulegue/internal/domain/web/model/setting"
	"modulegue/internal/domain/web/repository"
)

type RegisterJukirUseCase struct {
	settingRepo repository.SettingRepository
}

func NewRegisterJukirUseCase(
	settingRepo repository.SettingRepository,
) *RegisterJukirUseCase {
	return &RegisterJukirUseCase{
		settingRepo: settingRepo,
	}
}

func (uc *RegisterJukirUseCase) Execute(ctx context.Context, reqModel model.RegisterRequestModel) error {
	err := uc.settingRepo.RegisterJukir(ctx, reqModel)
	if err != nil {
		return fmt.Errorf("gagal register jukir : %w", err)
	}
	return nil
}
