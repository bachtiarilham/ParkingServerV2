package repository

import (
	"context"

	model "modulegue/internal/domain/mobile/model/payment_parking"
)

type PaymentRepository interface {
	PayParking(ctx context.Context, req model.PostPaymentParkingRequestModel) (*model.PaymentBusinessModel, error)
	ProcessPaymentSuccess(ctx context.Context, req model.PaymentBusinessModel) error
	ProcessPaymentFailed(ctx context.Context, req model.PaymentBusinessModel) error
}
