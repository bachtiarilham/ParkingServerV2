package repository

import (
	"context"
	model "modulegue/internal/domain/mobile/model/payment_gate"
)

type PaymentCallbackRepository interface {
	GetPaymentTransaction(ctx context.Context, txCode string) (*model.PaymentTransactionModel, error)
	GetPaymentStatus(ctx context.Context, txCode string) (status string, err error)
	GetPaymentStatusCash(ctx context.Context, sessionCode string) (status string, err error)
	ProcessParkingCallback(ctx context.Context, txCode string, extRef string, sessionID int64, req model.PaymentTransactionModel) error
	ProcessTopupCallback(ctx context.Context, txCode string, extRef string, userID int64, amount int64) error
	ProcessMembershipCallback(ctx context.Context, txCode string, extRef string, userID int64, packageID int64) error
}
