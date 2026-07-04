package payment

import (
	"context"
	"fmt"
	model "modulegue/internal/domain/mobile/model/payment"
	repository "modulegue/internal/domain/mobile/repository"
)

type PostParkingUseCase struct {
	postParkingRepo repository.PaymentRepository
}

func NewPostParkingUseCase(
	postParkingRepo repository.PaymentRepository,
) *PostParkingUseCase {
	return &PostParkingUseCase{
		postParkingRepo: postParkingRepo,
	}
}

func (uc *PostParkingUseCase) Execute(ctx context.Context, reqModel model.PostParkingRequestModel) (*model.PostParkingResponseModel, error) {
	result, err := uc.postParkingRepo.PostParking(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("post parking: %w", err)
	}

	if result == nil {
		return &model.PostParkingResponseModel{
			PaymentOptions: []model.PembayaranOptionModel{},
		}, nil
	}

	if result.PaymentOptions == nil {
		result.PaymentOptions = []model.PembayaranOptionModel{}
	}

	return result, nil
}
