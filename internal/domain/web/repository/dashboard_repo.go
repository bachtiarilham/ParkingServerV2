package repository

import (
	"context"
	"modulegue/internal/delivery/web/dto"
	"modulegue/internal/delivery/web/dto/request"
	"modulegue/internal/domain/web/finance"
	"modulegue/internal/domain/web/location"
	"modulegue/internal/domain/web/metrics"
	"modulegue/internal/domain/web/officer"
	"modulegue/internal/domain/web/settings"
	"modulegue/internal/domain/web/transaction"
	"time"
)

type DashboardRepository interface {
	GetDashboardOverview(ctx context.Context) (dto.DashboardOverview, error)
	GetMonitoringOverview(ctx context.Context) (dto.MonitoringOverview, error)
	GetOfficerOverview(ctx context.Context) (dto.OfficerOverview, error)
	GetTransactionsOverview(ctx context.Context) (dto.TransactionsOverview, error)
	GetSettingsOverview(ctx context.Context) (dto.SettingsOverview, error)
}

type LocationRepository interface {
	GetLocations(ctx context.Context) ([]location.ParkingLocation, error)
	GetLocationByID(ctx context.Context, locationID string) (location.ParkingLocation, error)
	UpdateLocationSettings(ctx context.Context, adminID int64, locationID string, tariffMotor, tariffMobil int64, operationalNote string) (location.ParkingLocation, error)
	SaveShiftTemplates(ctx context.Context, adminID int64, locationID string, templates []location.ParkingShiftTemplate) error
	GetHourlyTraffic(ctx context.Context, start, end time.Time) ([]metrics.HourlyTrafficPoint, error)
	GetHeatmapData(ctx context.Context, now time.Time) ([]metrics.HeatmapPoint, error)
}

type OfficerRepository interface {
	GetOfficerOptions(ctx context.Context) ([]officer.ParkingOfficerOption, error)
	UpdateOfficerStatus(ctx context.Context, adminID int64, officerID string, status string) (officer.ParkingOfficerOption, error)
	ApplyOfficerMutation(ctx context.Context, adminID int64, req request.ApplyOfficerMutationRequest) (officer.ParkingOfficerOption, error)
}

type SettingsRepository interface {
	UpdateSettingsOverview(ctx context.Context, adminID int64, req request.UpdateSettingsOverviewRequest) (dto.SettingsOverview, error)
	ListAlertRules(ctx context.Context) ([]settings.AlertRuleItem, error)
	ListShiftTemplates(ctx context.Context) ([]settings.ShiftTemplateItem, error)
	ListNotifications(ctx context.Context) ([]settings.NotificationItem, error)
	ListPaymentMethods(ctx context.Context) ([]settings.PaymentMethodItem, error)
	ListDefaultTariffs(ctx context.Context) ([]settings.DefaultTariffItem, error)
}

type TransactionRepository interface {
	CreateDisputeCase(ctx context.Context, adminID int64, req request.CreateDisputeCaseRequest) (transaction.DisputeCaseSummary, error)
	UpdateDisputeCaseStatus(ctx context.Context, adminID int64, disputeID string, req request.UpdateDisputeCaseStatusRequest) (transaction.DisputeCaseSummary, error)
	CreateRefundTransaction(ctx context.Context, adminID int64, req request.CreateRefundTransactionRequest) (finance.RefundTransactionSummary, error)
	UpdateRefundTransactionStatus(ctx context.Context, adminID int64, refundID string, req request.UpdateRefundStatusRequest) (finance.RefundTransactionSummary, error)
	CreateClosingBatch(ctx context.Context, adminID int64, req request.CreateClosingBatchRequest) (finance.ClosingBatchSummary, error)
	UpdateClosingBatchStatus(ctx context.Context, adminID int64, closingID string, req request.UpdateClosingStatusRequest) (finance.ClosingBatchSummary, error)
}
