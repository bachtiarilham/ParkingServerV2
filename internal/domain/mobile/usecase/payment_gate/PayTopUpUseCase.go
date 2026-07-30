package paymentgate

import (
	"context"
	"fmt"

	model "modulegue/internal/domain/mobile/model/payment_gate"
	"modulegue/internal/domain/mobile/repository"
	midtransService "modulegue/internal/service/payment_gateway/midtrans"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type PayTopUpUseCase struct {
	repo           repository.PaymentGateRepository
	midtransClient *midtransService.MidtransClient
}

func NewPayTopUpUseCase(
	repo repository.PaymentGateRepository,
	midtransClient *midtransService.MidtransClient,
) *PayTopUpUseCase {
	return &PayTopUpUseCase{
		repo:           repo,
		midtransClient: midtransClient,
	}
}

func (uc *PayTopUpUseCase) Execute(ctx context.Context, reqmodel model.PayRequestModel) (*model.PayResponseModel, error) {
	result, err := uc.repo.PayTopUp(ctx, reqmodel)
	if err != nil {
		return nil, fmt.Errorf("pay topup repo: %w", err)
	}
	if result == nil {
		return nil, fmt.Errorf("pay topup returned empty result")
	}

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  result.OrderID,
			GrossAmt: result.GrossAmount,
		},
	}

	snapResp, err := uc.midtransClient.CreateSnapToken(ctx, snapReq)
	if err != nil {
		return nil, fmt.Errorf("midtrans snap request: %w", err)
	}

	result.SnapToken = snapResp.Token
	result.RedirectURL = snapResp.RedirectURL
	result.Status = "PENDING"

	return result, nil
}
