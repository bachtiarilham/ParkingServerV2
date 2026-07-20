package setting

import (
	"context"
	"fmt"
	model "modulegue/internal/domain/web/model/setting"
	"modulegue/internal/domain/web/repository"
)

type SaveScheduleUseCase struct {
	settingRepo repository.SettingRepository
}

func NewSaveScheduleUseCase(
	settingRepo repository.SettingRepository,
) *SaveScheduleUseCase {
	return &SaveScheduleUseCase{
		settingRepo: settingRepo,
	}
}

func (uc *SaveScheduleUseCase) Execute(ctx context.Context, reqModel model.SaveScheduleRequestModel) error {
	err := uc.settingRepo.SaveSchedule(ctx, reqModel)
	if err != nil {
		return fmt.Errorf("gagal menyimpan schedule jukir : %w", err)
	}
	return nil
}
