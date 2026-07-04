package payment

import (
	"context"
	"fmt"
	model "modulegue/internal/domain/mobile/model/payment"
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
	result, err := uc.postPaymentParkingRepo.PostPaymentParking(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("post payment parking: %w", err)
	}

	if result == nil {
		return &model.PostPaymentParkingResponseModel{
			Details: []model.PostPaymentParkingDetailItemModel{},
		}, nil
	}

	if result.Details == nil {
		result.Details = []model.PostPaymentParkingDetailItemModel{}
	}

	return result, nil
}
