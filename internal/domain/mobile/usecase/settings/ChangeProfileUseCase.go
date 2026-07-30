package settings

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"modulegue/config"
	model "modulegue/internal/domain/mobile/model/settings"
	repository "modulegue/internal/domain/mobile/repository"
)

type ChangeProfileUseCase struct {
	changeProfileRepository repository.SettingsRepository
}

func NewChangeProfileUseCase(
	changeProfileRepository repository.SettingsRepository,
) *ChangeProfileUseCase {
	return &ChangeProfileUseCase{
		changeProfileRepository: changeProfileRepository,
	}
}

func (uc *ChangeProfileUseCase) Execute(ctx context.Context, req model.SettingsModel) error {
	// If photo bytes are provided, save to D:/parking_data/images and set FotoProfil URL
	if len(req.FotoProfilBytes) > 0 {
		dir := "D:/parking_data/images"
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create directory: %w", err)
		}

		ext := filepath.Ext(req.FotoProfilName)
		if ext == "" {
			ext = ".jpg"
		}
		filename := fmt.Sprintf("photo_%d_%d%s", time.Now().UnixNano(), req.UserId, ext)
		filePath := filepath.Join(dir, filename)

		if err := os.WriteFile(filePath, req.FotoProfilBytes, 0644); err != nil {
			return fmt.Errorf("save profile photo: %w", err)
		}

		// Retrieve base URL from config and construct the public image URL
		cfg := config.Load()
		publicURL := fmt.Sprintf("%s/images/%s", cfg.ImageBaseURL, filename)
		req.FotoProfil = &publicURL
	}

	err := uc.changeProfileRepository.ChangeProfile(ctx, req)
	if err != nil {
		return fmt.Errorf("change profile: %w", err)
	}

	return nil
}
