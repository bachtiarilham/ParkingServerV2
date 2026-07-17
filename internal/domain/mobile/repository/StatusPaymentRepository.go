package repository

import (
	"context"
	model "modulegue/internal/domain/mobile/model/payment_parking"
)

type StatusPaymentRepository interface {
	GetPembayaranStatus(ctx context.Context, sessionId string) (*model.PostPaymentParkingResponseModel, error)
}
