package setting

import (
	"context"
	"fmt"
	model "modulegue/internal/domain/web/model/setting"
	"modulegue/internal/domain/web/repository"
)

type UpdateProfileUseCase struct {
	settingRepo repository.SettingRepository
}

func NewUpdateProfileUseCase(
	settingRepo repository.SettingRepository,
) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{
		settingRepo: settingRepo,
	}
}

func (uc *UpdateProfileUseCase) Execute(ctx context.Context, reqModel model.UpdateProfilRequestModel) error {
	err := uc.settingRepo.UpdateProfil(ctx, reqModel)
	if err != nil {
		return fmt.Errorf("gagal update profil : %w", err)
	}
	return nil
}
