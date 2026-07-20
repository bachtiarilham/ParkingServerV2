package setting

import (
	"context"
	"fmt"
	model "modulegue/internal/domain/web/model/setting"
	"modulegue/internal/domain/web/repository"
)

type ShowSelectedJukirUseCase struct {
	settingRepo repository.SettingRepository
}

func NewShowSelectedJukirUseCase(
	settingRepo repository.SettingRepository,
) *ShowSelectedJukirUseCase {
	return &ShowSelectedJukirUseCase{
		settingRepo: settingRepo,
	}
}

func (uc *ShowSelectedJukirUseCase) Execute(ctx context.Context, reqModel model.ShowSelectedJukirRequestModel) (*model.ShowSelectedJukirResponseModel, error) {
	resp, err := uc.settingRepo.ShowSelectedJukir(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("gagal memunculkan jukir : %w", err)
	}
	return resp, nil
}
