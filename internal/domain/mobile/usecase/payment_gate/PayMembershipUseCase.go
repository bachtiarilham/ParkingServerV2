package paymentgate

import (
	"context"
	"fmt"
	"strings"

	model "modulegue/internal/domain/mobile/model/payment_gate"
	"modulegue/internal/domain/mobile/repository"
	midtransService "modulegue/internal/service/payment_gateway/midtrans"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type PayMembershipUseCase struct {
	repo           repository.PaymentGateRepository
	midtransClient *midtransService.MidtransClient
}

func NewPayMembershipUseCase(
	repo repository.PaymentGateRepository,
	midtransClient *midtransService.MidtransClient,
) *PayMembershipUseCase {
	return &PayMembershipUseCase{
		repo:           repo,
		midtransClient: midtransClient,
	}
}

func (uc *PayMembershipUseCase) Execute(ctx context.Context, reqmodel model.PayRequestModel) (*model.PayResponseModel, error) {
	if reqmodel.TargetID != nil {
		*reqmodel.TargetID = strings.ToLower(*reqmodel.TargetID)
	}

	result, err := uc.repo.PayMembership(ctx, reqmodel)
	if err != nil {
		return nil, fmt.Errorf("pay membership repo: %w", err)
	}
	if result == nil {
		return nil, fmt.Errorf("pay membership returned empty result")
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
