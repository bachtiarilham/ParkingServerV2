package monitoring

import (
	"context"
	"fmt"
	model "modulegue/internal/domain/web/model/monitoring"
	"modulegue/internal/domain/web/repository"
)

type MonitoringUseCase struct {
	monitoringRepo repository.MonitoringRepository
}

func NewMonitoringUseCase(
	monitoringRepo repository.MonitoringRepository,
) *MonitoringUseCase {
	return &MonitoringUseCase{
		monitoringRepo: monitoringRepo,
	}
}

func (uc *MonitoringUseCase) Execute(ctx context.Context, reqModel model.MonitoringRequestModel) (*model.MonitoringResponseModel, error) {
	resp, err := uc.monitoringRepo.GetMonitoring(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data monitoring : %w", err)
	}
	return resp, nil
}
