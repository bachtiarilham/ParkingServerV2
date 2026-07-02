package usecase

import (
	"context"
	"fmt"

	model "modulegue/internal/domain/mobile/model/laporan"
	"modulegue/internal/domain/mobile/repository"
	"modulegue/internal/middleware"
)

type GetLaporanUseCase struct {
	laporanRepo repository.LaporanRepository
}

func NewGetLaporanUseCase(
	laporanRepo repository.LaporanRepository,
) *GetLaporanUseCase {
	return &GetLaporanUseCase{
		laporanRepo: laporanRepo,
	}
}

func (uc *GetLaporanUseCase) Execute(ctx context.Context, reqModel model.LaporanRequestModel) (*model.LaporanModel, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok && reqModel.UserID == 0 {
		return nil, fmt.Errorf("user not authenticated")
	}

	roleID, ok := middleware.RoleIDFromContext(ctx)
	if !ok && reqModel.RoleID == 0 {
		return nil, fmt.Errorf("role not found in context")
	}

	if reqModel.UserID == 0 {
		reqModel.UserID = userID
	}
	if reqModel.RoleID == 0 {
		reqModel.RoleID = roleID
	}

	result, err := uc.laporanRepo.GetLaporan(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("get laporan: %w", err)
	}

	if result == nil {
		return &model.LaporanModel{
			ChartBars:          &model.LaporanChartBarModel{ChartBars: []model.LaporanChartBarItemModel{}},
			PaymentSummaries:   &model.LaporanPaymentSummaryModel{PaymentSummaries: []model.LaporanPaymentSummaryItemModel{}},
			RecentTransactions: &model.LaporanRecentTransactionModel{RecentTransactions: []model.LaporanRecentTransactionItemModel{}},
		}, nil
	}

	if result.ChartBars == nil {
		result.ChartBars = &model.LaporanChartBarModel{ChartBars: []model.LaporanChartBarItemModel{}}
	}
	if result.PaymentSummaries == nil {
		result.PaymentSummaries = &model.LaporanPaymentSummaryModel{PaymentSummaries: []model.LaporanPaymentSummaryItemModel{}}
	}
	if result.RecentTransactions == nil {
		result.RecentTransactions = &model.LaporanRecentTransactionModel{RecentTransactions: []model.LaporanRecentTransactionItemModel{}}
	}

	return result, nil
}
