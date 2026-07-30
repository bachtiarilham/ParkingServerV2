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

type PayParkingUseCase struct {
	payParkingRepo repository.PaymentGateRepository
	midtransClient *midtransService.MidtransClient
}

func NewPayParkingUseCase(
	payParkingRepo repository.PaymentGateRepository,
	midtransClient *midtransService.MidtransClient,
) *PayParkingUseCase {
	return &PayParkingUseCase{
		payParkingRepo: payParkingRepo,
		midtransClient: midtransClient,
	}
}

func (uc *PayParkingUseCase) Execute(ctx context.Context, reqModel model.PayRequestModel) (*model.PayResponseModel, error) {

	ismember, err := uc.payParkingRepo.IsMember(ctx, *reqModel.TargetID)
	if err != nil {
		return nil, fmt.Errorf("ismember repo: %w", err)
	}

	if ismember.IsMember {
		reqModel.Amount = (70 * reqModel.Amount) / 100
	}

	reqModel.JukirShare = (reqModel.Amount * 30) / 100
	reqModel.CompanyShare = (reqModel.Amount * 49) / 100
	reqModel.GovShare = (reqModel.Amount * 20) / 100
	reqModel.MidtransShare = int64(reqModel.Amount * 1 / 100)

	result, err := uc.payParkingRepo.PayParking(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("pay parking repo: %w", err)
	}
	if result == nil {
		return nil, fmt.Errorf("pay parking returned empty result")
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
