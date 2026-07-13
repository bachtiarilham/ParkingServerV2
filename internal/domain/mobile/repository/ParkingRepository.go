package repository

import (
	"context"
	model "modulegue/internal/domain/mobile/model/parking"
)

type ParkingRepository interface {
	PostParking(ctx context.Context, req model.PostParkingRequestModel) (*model.ParkingBusinessModel, error)
	InsertParkingSession(ctx context.Context, req model.ParkingBusinessModel) (sessionId int64, err error)
	InsertQrisString(ctx context.Context, req model.ParkingBusinessModel) error
	ReturnToApp(ctx context.Context, req model.ParkingBusinessModel) (*model.PostParkingResponseModel, error)
}
