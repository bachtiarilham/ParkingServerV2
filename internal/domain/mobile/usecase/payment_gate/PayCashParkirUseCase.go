package paymentgate

import (
	"context"
	"fmt"

	model "modulegue/internal/domain/mobile/model/payment_gate"
	"modulegue/internal/domain/mobile/repository"
	midtransService "modulegue/internal/service/payment_gateway/midtrans"
)

type PayCashParkirUseCase struct {
	repo           repository.PaymentGateRepository
	reposync       repository.SyncRepo
	repocall       repository.PaymentCallbackRepository
	midtransClient *midtransService.MidtransClient
}

func NewPayCashParkirUseCase(
	repo repository.PaymentGateRepository,
	reposync repository.SyncRepo,
	repocall repository.PaymentCallbackRepository,
) *PayCashParkirUseCase {
	return &PayCashParkirUseCase{
		repo:     repo,
		reposync: reposync,
		repocall: repocall,
	}
}

func (uc *PayCashParkirUseCase) Execute(ctx context.Context, reqmodel model.PayRequestModel) (*model.PayResponseModel, error) {
	reqmodel.JukirShare = (reqmodel.Amount * 30) / 100
	reqmodel.CompanyShare = (reqmodel.Amount * 49) / 100
	reqmodel.GovShare = (reqmodel.Amount * 20) / 100
	reqmodel.MidtransShare = int64(reqmodel.Amount * 1 / 100)

	txCode, result, err := uc.repo.PayCashParking(ctx, reqmodel)
	if err != nil {
		return nil, fmt.Errorf("pay cash repo: %w", err)
	}

	txInfo, err := uc.repocall.GetPaymentTransaction(ctx, txCode)
	if err != nil {
		return nil, fmt.Errorf("get transaction details: %w", err)
	}

	err = uc.reposync.SyncAfterParkir(ctx, reqmodel.UserID, txInfo.ReferenceID, reqmodel.Amount)
	if err != nil {
		return nil, fmt.Errorf("process sync parkir: %w", err)
	}

	return result, nil
}
