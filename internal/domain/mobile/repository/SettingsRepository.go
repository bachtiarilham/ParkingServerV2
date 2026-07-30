package repository

import (
	"context"
	model "modulegue/internal/domain/mobile/model/settings"
)

type SettingsRepository interface {
	ChangeProfile(ctx context.Context, req model.SettingsModel) error
}
