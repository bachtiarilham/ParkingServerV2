package repository

import (
	"context"
	model "modulegue/internal/domain/web/model/monitoring"
)

type MonitoringRepository interface {
	GetMonitoring(context.Context, model.MonitoringRequestModel) (*model.MonitoringResponseModel, error)
}
