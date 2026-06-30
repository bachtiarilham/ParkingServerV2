package web

import (
	"context"
	"fmt"
	"strings"

	"modulegue/internal/domain/dashboard"
)

type DashboardUseCase struct {
	repo dashboard.Repository
}

func NewDashboardUseCase(repo dashboard.Repository) *DashboardUseCase {
	return &DashboardUseCase{repo: repo}
}

func (uc *DashboardUseCase) GetDashboardOverview(ctx context.Context, actorID int64) (dashboard.DashboardOverview, error) {
	if _, err := uc.authorizeAdmin(ctx, actorID); err != nil {
		return dashboard.DashboardOverview{}, err
	}
	return uc.repo.GetDashboardOverview(ctx)
}

func (uc *DashboardUseCase) GetMonitoringOverview(ctx context.Context, actorID int64) (dashboard.MonitoringOverview, error) {
	if _, err := uc.authorizeAdmin(ctx, actorID); err != nil {
		return dashboard.MonitoringOverview{}, err
	}
	return uc.repo.GetMonitoringOverview(ctx)
}

func (uc *DashboardUseCase) GetOfficerOverview(ctx context.Context, actorID int64) (dashboard.OfficerOverview, error) {
	if _, err := uc.authorizeAdmin(ctx, actorID); err != nil {
		return dashboard.OfficerOverview{}, err
	}
	return uc.repo.GetOfficerOverview(ctx)
}

func (uc *DashboardUseCase) GetTransactionsOverview(ctx context.Context, actorID int64) (dashboard.TransactionsOverview, error) {
	if _, err := uc.authorizeAdmin(ctx, actorID); err != nil {
		return dashboard.TransactionsOverview{}, err
	}
	return uc.repo.GetTransactionsOverview(ctx)
}

func (uc *DashboardUseCase) GetSettingsOverview(ctx context.Context, actorID int64) (dashboard.SettingsOverview, error) {
	if _, err := uc.authorizeAdmin(ctx, actorID); err != nil {
		return dashboard.SettingsOverview{}, err
	}
	return uc.repo.GetSettingsOverview(ctx)
}

func (uc *DashboardUseCase) authorizeAdmin(ctx context.Context, actorID int64) (*dashboard.AuthRecord, error) {
	admin, err := uc.repo.FindAdminByID(ctx, actorID)
	if err != nil {
		return nil, fmt.Errorf("admin tidak ditemukan: %w", err)
	}
	if !isAdminRole(admin.RoleCode) {
		return nil, ErrForbiddenRole
	}
	if strings.TrimSpace(admin.Username) == "" {
		return nil, ErrForbiddenRole
	}
	return admin, nil
}
