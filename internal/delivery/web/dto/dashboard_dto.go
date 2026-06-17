package dto

import (
	"modulegue/internal/domain/web/location"
	"modulegue/internal/domain/web/metrics"
)

type DashboardOverview struct {
	DashboardStats        []metrics.StatCard           `json:"dashboardStats"`
	DashboardTransactions []metrics.RowItem            `json:"dashboardTransactions"`
	FieldOfficers         []metrics.RowItem            `json:"fieldOfficers"`
	DashboardAlerts       []metrics.AlertItem          `json:"dashboardAlerts"`
	HourlyTraffic         []metrics.HourlyTrafficPoint `json:"hourlyTraffic"`
	RevenueByLocation     []metrics.LocationMetric     `json:"revenueByLocation"`
	OccupancyByLocation   []metrics.LocationMetric     `json:"occupancyByLocation"`
	ComparisonMetrics     []metrics.ComparisonMetric   `json:"comparisonMetrics"`
	ParkingHeatmap        []metrics.HeatmapPoint       `json:"parkingHeatmap"`
	PriorityActions       []metrics.ActionItem         `json:"priorityActions"`
	ParkingLocations      []location.ParkingLocation   `json:"parkingLocations"`
}
