package repository

import (
	"context"
	"modulegue/internal/delivery/web/dto"
)

type dashboardRepo interface {
	GetDashboardOverview(ctx context.Context) (dto.DashboardOverview, error)
	GetMonitoringOverview(ctx context.Context) (dto.MonitoringOverview, error)
	GetOfficerOverview(ctx context.Context) (dto.OfficerOverview, error)
	GetTransactionsOverview(ctx context.Context) (dto.TransactionsOverview, error)
	GetSettingsOverview(ctx context.Context) (dto.SettingsOverview, error)
}
