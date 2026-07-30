package repository

import (
	"context"
	model "modulegue/internal/domain/mobile/model/parking"
)

type ParkingRepository interface {
	GetParkingMetadata(ctx context.Context, req model.PostParkingRequestModel) (*model.ParkingBusinessModel, error)
	CreateParkingTransaction(ctx context.Context, req1 model.PostParkingRequestModel, req model.ParkingBusinessModel) (*model.PostParkingResponseModel, error)
}
