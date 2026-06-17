package dto

import (
	"modulegue/internal/domain/web/location"
	"modulegue/internal/domain/web/metrics"
	"modulegue/internal/domain/web/officer"
)

type OfficerOverview struct {
	OfficerStats          []metrics.StatCard             `json:"officerStats"`
	ParkingOfficerOptions []officer.ParkingOfficerOption `json:"parkingOfficerOptions"`
	ParkingLocations      []location.ParkingLocation     `json:"parkingLocations"`
}
