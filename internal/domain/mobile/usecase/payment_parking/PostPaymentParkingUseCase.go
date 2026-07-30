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
	prepared, err := uc.postPaymentParkingRepo.PayParking(ctx, reqModel)
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

	if err := uc.postPaymentParkingRepo.ProcessPaymentSuccess(ctx, businessModel); err != nil {
		return nil, fmt.Errorf("process payment success: %w", err)
	}

	return &model.PostPaymentParkingResponseModel{
		SessionId:   prepared.SessionId,
		PaymentCode: businessModel.PaymentCode,
	}, nil
}
