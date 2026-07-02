package repository

import (
	"context"

	model "modulegue/internal/domain/mobile/model/payment"
)

type PaymentRepository interface {
	Submit(ctx context.Context, req model.SubmitQrRequestModel) (*model.PembayaranModel, error)
	GetPembayaran(ctx context.Context, sessionId int64) (*model.PembayaranModel, error)
}
