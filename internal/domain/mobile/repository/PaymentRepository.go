package repository

import (
	"context"

	model "modulegue/internal/domain/mobile/model/payment_parking"
)

type PaymentRepository interface {
	PostPaymentParking(ctx context.Context, req model.PostPaymentParkingRequestModel) (*model.PaymentBusinessModel, error)
	BindCustomerToSessionAndTransaction(ctx context.Context, req model.PaymentBusinessModel) error
	//if success
	UpdatePaymentTransactionSuccess(ctx context.Context, req model.PaymentBusinessModel) error
	UpdateParkingSessionSuccess(ctx context.Context, req model.PaymentBusinessModel) error
	BuatParkingReceipt(ctx context.Context, req model.PaymentBusinessModel) error
	BuatFinancialParkingTransaction(ctx context.Context, req model.PaymentBusinessModel) error
	InsertNotifikasiSuccess(ctx context.Context, req model.PaymentBusinessModel) error
	//if gagal
	UpdatePaymentTransactionFailed(ctx context.Context, req model.PaymentBusinessModel) error
	UpdateParkingSessionFailed(ctx context.Context, req model.PaymentBusinessModel) error
	InsertNotifikasiFailed(ctx context.Context, req model.PaymentBusinessModel) error
}
