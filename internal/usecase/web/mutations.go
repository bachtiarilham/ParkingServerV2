package web

import (
	"context"

	"modulegue/internal/domain/dashboard"
)

func (uc *DashboardUseCase) UpdateSettingsOverview(ctx context.Context, actorID int64, req dashboard.UpdateSettingsOverviewRequest) (dashboard.SettingsOverview, error) {
	if _, err := uc.authorizeAdmin(ctx, actorID); err != nil {
		return dashboard.SettingsOverview{}, err
	}
	return uc.repo.UpdateSettingsOverview(ctx, actorID, req)
}

func (uc *DashboardUseCase) UpdateLocationSettings(ctx context.Context, actorID int64, locationID string, req dashboard.UpdateLocationSettingsRequest) (dashboard.ParkingLocation, error) {
	if _, err := uc.authorizeAdmin(ctx, actorID); err != nil {
		return dashboard.ParkingLocation{}, err
	}
	return uc.repo.UpdateLocationSettings(ctx, actorID, locationID, req)
}

func (uc *DashboardUseCase) SaveLocationShiftTemplates(ctx context.Context, actorID int64, locationID string, templates []dashboard.ParkingShiftTemplate) (dashboard.ParkingLocation, error) {
	if _, err := uc.authorizeAdmin(ctx, actorID); err != nil {
		return dashboard.ParkingLocation{}, err
	}
	return uc.repo.SaveLocationShiftTemplates(ctx, actorID, locationID, templates)
}

func (uc *DashboardUseCase) UpdateOfficerStatus(ctx context.Context, actorID int64, officerID string, status string) (dashboard.ParkingOfficerOption, error) {
	if _, err := uc.authorizeAdmin(ctx, actorID); err != nil {
		return dashboard.ParkingOfficerOption{}, err
	}
	return uc.repo.UpdateOfficerStatus(ctx, actorID, officerID, status)
}

func (uc *DashboardUseCase) ApplyOfficerMutation(ctx context.Context, actorID int64, req dashboard.ApplyOfficerMutationRequest) (dashboard.ParkingOfficerOption, error) {
	if _, err := uc.authorizeAdmin(ctx, actorID); err != nil {
		return dashboard.ParkingOfficerOption{}, err
	}
	return uc.repo.ApplyOfficerMutation(ctx, actorID, req)
}

func (uc *DashboardUseCase) CreateDisputeCase(ctx context.Context, actorID int64, req dashboard.CreateDisputeCaseRequest) (dashboard.DisputeCaseSummary, error) {
	if _, err := uc.authorizeAdmin(ctx, actorID); err != nil {
		return dashboard.DisputeCaseSummary{}, err
	}
	return uc.repo.CreateDisputeCase(ctx, actorID, req)
}

func (uc *DashboardUseCase) UpdateDisputeCaseStatus(ctx context.Context, actorID int64, disputeID string, req dashboard.UpdateDisputeCaseStatusRequest) (dashboard.DisputeCaseSummary, error) {
	if _, err := uc.authorizeAdmin(ctx, actorID); err != nil {
		return dashboard.DisputeCaseSummary{}, err
	}
	return uc.repo.UpdateDisputeCaseStatus(ctx, actorID, disputeID, req)
}

func (uc *DashboardUseCase) CreateRefundTransaction(ctx context.Context, actorID int64, req dashboard.CreateRefundTransactionRequest) (dashboard.RefundTransactionSummary, error) {
	if _, err := uc.authorizeAdmin(ctx, actorID); err != nil {
		return dashboard.RefundTransactionSummary{}, err
	}
	return uc.repo.CreateRefundTransaction(ctx, actorID, req)
}

func (uc *DashboardUseCase) UpdateRefundTransactionStatus(ctx context.Context, actorID int64, refundID string, req dashboard.UpdateRefundStatusRequest) (dashboard.RefundTransactionSummary, error) {
	if _, err := uc.authorizeAdmin(ctx, actorID); err != nil {
		return dashboard.RefundTransactionSummary{}, err
	}
	return uc.repo.UpdateRefundTransactionStatus(ctx, actorID, refundID, req)
}

func (uc *DashboardUseCase) CreateClosingBatch(ctx context.Context, actorID int64, req dashboard.CreateClosingBatchRequest) (dashboard.ClosingBatchSummary, error) {
	if _, err := uc.authorizeAdmin(ctx, actorID); err != nil {
		return dashboard.ClosingBatchSummary{}, err
	}
	return uc.repo.CreateClosingBatch(ctx, actorID, req)
}

func (uc *DashboardUseCase) UpdateClosingBatchStatus(ctx context.Context, actorID int64, closingID string, req dashboard.UpdateClosingStatusRequest) (dashboard.ClosingBatchSummary, error) {
	if _, err := uc.authorizeAdmin(ctx, actorID); err != nil {
		return dashboard.ClosingBatchSummary{}, err
	}
	return uc.repo.UpdateClosingBatchStatus(ctx, actorID, closingID, req)
}
