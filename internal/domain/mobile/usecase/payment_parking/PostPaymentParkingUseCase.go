package payment

import (
	"context"
	"fmt"

	model "modulegue/internal/domain/mobile/model/payment_parking"
	repository "modulegue/internal/domain/mobile/repository"
)

type PostPaymentParkingUseCase struct {
	postPaymentParkingRepo repository.PaymentRepository
}

func NewPostPaymentParkingUseCase(
	postPaymentParkingRepo repository.PaymentRepository,
) *PostPaymentParkingUseCase {
	return &PostPaymentParkingUseCase{
		postPaymentParkingRepo: postPaymentParkingRepo,
	}
}

func (uc *PostPaymentParkingUseCase) Execute(ctx context.Context, reqModel model.PostPaymentParkingRequestModel) (*model.PostPaymentParkingResponseModel, error) {
	prepared, err := uc.postPaymentParkingRepo.PostPaymentParking(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("post payment parking: %w", err)
	}
	if prepared == nil {
		return nil, fmt.Errorf("post payment parking returned empty data")
	}

	businessModel := model.PaymentBusinessModel{
		SessionId:       prepared.SessionId,
		CustomerUserId:  reqModel.CustomerUserId,
		SessionCode:     prepared.SessionCode,
		TransactionCode: prepared.TransactionCode,
		PaymentCode: func() string {
			if prepared.PaymentCode != "" {
				return prepared.PaymentCode
			}
			return reqModel.SessionCode
		}(),
		FailedReason: "",
	}

	if err := uc.postPaymentParkingRepo.BindCustomerToSessionAndTransaction(ctx, businessModel); err != nil {
		return nil, fmt.Errorf("bind customer to session and transaction: %w", err)
	}

	if err := uc.postPaymentParkingRepo.UpdatePaymentTransactionSuccess(ctx, businessModel); err != nil {
		return nil, fmt.Errorf("update payment transaction success: %w", err)
	}

	if err := uc.postPaymentParkingRepo.UpdateParkingSessionSuccess(ctx, businessModel); err != nil {
		return nil, fmt.Errorf("update parking session success: %w", err)
	}

	if err := uc.postPaymentParkingRepo.BuatParkingReceipt(ctx, businessModel); err != nil {
		return nil, fmt.Errorf("buat parking receipt: %w", err)
	}

	if err := uc.postPaymentParkingRepo.BuatFinancialParkingTransaction(ctx, businessModel); err != nil {
		return nil, fmt.Errorf("buat financial parking transaction: %w", err)
	}

	if err := uc.postPaymentParkingRepo.InsertNotifikasiSuccess(ctx, businessModel); err != nil {
		return nil, fmt.Errorf("insert notifikasi success: %w", err)
	}

	return &model.PostPaymentParkingResponseModel{
		SessionId:   prepared.SessionId,
		PaymentCode: businessModel.PaymentCode,
	}, nil
}
