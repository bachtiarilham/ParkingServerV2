package repository

import (
	"context"
	model "modulegue/internal/domain/mobile/model/payment_gate"
)

type PaymentGateRepository interface {
	PayParking(ctx context.Context, req model.PayRequestModel) (*model.PayResponseModel, error)
	PayTransfer(ctx context.Context, req model.PayRequestModel) (*model.PayResponseModel, error)
	PayMembership(ctx context.Context, req model.PayRequestModel) (*model.PayResponseModel, error)
	PayTopUp(ctx context.Context, req model.PayRequestModel) (*model.PayResponseModel, error)
	PayCashParking(ctx context.Context, req model.PayRequestModel) (string, *model.PayResponseModel, error)
	IsMember(ctx context.Context, sessioncode string) (*model.MemberCheckResult, error)
}
