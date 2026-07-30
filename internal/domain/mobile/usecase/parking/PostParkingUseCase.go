package parking

import (
	"context"
	"fmt"
	"strings"
	"time"

	model "modulegue/internal/domain/mobile/model/parking"
	repository "modulegue/internal/domain/mobile/repository"
)

type PostParkingUseCase struct {
	postParkingRepo repository.ParkingRepository
}

func NewPostParkingUseCase(
	postParkingRepo repository.ParkingRepository,
) *PostParkingUseCase {
	return &PostParkingUseCase{
		postParkingRepo: postParkingRepo,
	}
}

func (uc *PostParkingUseCase) Execute(ctx context.Context, reqModel model.PostParkingRequestModel) (*model.PostParkingResponseModel, error) {
	strings.ToUpper(strings.TrimSpace(reqModel.PlateNumber))

	businessModel, err := uc.postParkingRepo.GetParkingMetadata(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("post parking lookup: %w", err)
	}
	if businessModel == nil {
		return nil, fmt.Errorf("post parking lookup returned empty data")
	}

	result, err := uc.postParkingRepo.CreateParkingTransaction(ctx, reqModel, *businessModel)
	if err != nil {
		return nil, fmt.Errorf("create parking: %w", err)
	}

	// result.QrExpired = time.Duration(15 * time.Minute)
	result.QrExpired = time.Now().Add(15 * time.Minute)
	result.BiayaParkir = reqModel.BiayaParkir
	return result, nil
}
