package repository

import (
	"context"

	model "modulegue/internal/domain/mobile/model/payment"
)

type PaymentRepository interface {
	PostParking(ctx context.Context, req model.PostParkingRequestModel) (*model.PostParkingResponseModel, error)
	PostPaymentParking(ctx context.Context, req model.PostPaymentParkingRequestModel) (*model.PostPaymentParkingResponseModel, error)
	GetPembayaranStatus(ctx context.Context, sessionId string) (*model.PostPaymentParkingResponseModel, error)
}
