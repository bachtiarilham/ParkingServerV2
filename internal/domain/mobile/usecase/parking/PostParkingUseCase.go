package parking

import (
	"context"
	"fmt"
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
	businessModel, err := uc.postParkingRepo.PostParking(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("post parking lookup: %w", err)
	}
	if businessModel == nil {
		return nil, fmt.Errorf("post parking lookup returned empty data")
	}

	now := time.Now().UTC()
	businessModel.PlateNumber = reqModel.PlateNumber
	businessModel.SessionCode = buildParkingCode("SESSION", now)
	businessModel.TransactionCode = buildParkingCode("TRX", now)
	businessModel.PaymentCode = buildParkingCode("PAY", now)
	businessModel.QrisString = businessModel.SessionCode
	businessModel.ExternalReference = businessModel.PaymentCode
	businessModel.ProviderName = "DEV_PROVIDER"
	businessModel.ParkingStatusId = 1
	businessModel.PaymentStatusId = 1

	sessionID, err := uc.postParkingRepo.InsertParkingSession(ctx, *businessModel)
	if err != nil {
		return nil, fmt.Errorf("insert parking session: %w", err)
	}
	businessModel.SessionId = sessionID

	if err := uc.postParkingRepo.InsertQrisString(ctx, *businessModel); err != nil {
		return nil, fmt.Errorf("insert qris string: %w", err)
	}

	result, err := uc.postParkingRepo.ReturnToApp(ctx, *businessModel)
	if err != nil {
		return nil, fmt.Errorf("return to app: %w", err)
	}
	if result != nil {
		return result, nil
	}

	return nil, nil
}

func buildParkingCode(prefix string, t time.Time) string {
	return fmt.Sprintf("%s-%s-%06d", prefix, t.Format("20060102"), t.Nanosecond()%1000000)
}
