package repository

import (
	"context"

	model "modulegue/internal/domain/mobile/model/helper"
)

type HelperRepository interface {
	GetLokasi(ctx context.Context, userId int64) (*model.LokasiModel, error)
	GetTarif(ctx context.Context, userId int64) (*model.TarifModel, error)
	GetNominalPayment(ctx context.Context) (*model.NominalPaymentModel, error)
	GetPaymentMethod(ctx context.Context) (*model.PaymentMethodModel, error)
}
