package dto

import (
	"modulegue/internal/domain/web/location"
	"modulegue/internal/domain/web/metrics"
	"modulegue/internal/domain/web/officer"
	"modulegue/internal/domain/web/settings"
)

type MonitoringOverview struct {
	TopFilters            settings.TopFilters            `json:"topFilters"`
	MonitoringZones       []metrics.RowItem              `json:"monitoringZones"`
	ParkingLocations      []location.ParkingLocation     `json:"parkingLocations"`
	ParkingOfficerOptions []officer.ParkingOfficerOption `json:"parkingOfficerOptions"`
}
