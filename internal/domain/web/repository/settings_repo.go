package repository

import (
	"context"
	"modulegue/internal/delivery/web/dto"
	"modulegue/internal/delivery/web/dto/request"
	"modulegue/internal/domain/web/location"
)

type settingsRepo interface {
	UpdateSettingsOverview(ctx context.Context, adminID int64, req request.UpdateSettingsOverviewRequest) (dto.SettingsOverview, error)
	UpdateLocationSettings(ctx context.Context, adminID int64, locationID string, req request.UpdateLocationSettingsRequest) (location.ParkingLocation, error)
	SaveLocationShiftTemplates(ctx context.Context, adminID int64, locationID string, templates []location.ParkingShiftTemplate) error
}
