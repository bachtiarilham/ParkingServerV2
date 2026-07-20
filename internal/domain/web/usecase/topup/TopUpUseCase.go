package topup

import (
	"context"
	"errors"
	"fmt"
	"time"

	model "modulegue/internal/domain/web/model/topup"
	repository "modulegue/internal/domain/web/repository"
)

type TopUpUseCase struct {
	topupRepo repository.TopUpRepository
}

func NewTopUpUseCase(topupRepo repository.TopUpRepository) *TopUpUseCase {
	return &TopUpUseCase{
		topupRepo: topupRepo,
	}
}

func (uc *TopUpUseCase) Execute(ctx context.Context, reqModel model.TopUpRequestModel) (*model.TopUpResponseModel, error) {
	if reqModel.IDUser <= 0 {
		return nil, errors.New("id_user harus lebih dari 0")
	}
	if reqModel.NominalTopUp <= 0 {
		return nil, errors.New("nominal_topup harus lebih dari 0")
	}

	if reqModel.TopUpCode == "" {
		reqModel.TopUpCode = generateTopUpCode()
	}
	if reqModel.ExternalReference == "" {
		reqModel.ExternalReference = fmt.Sprintf("MANUAL-%s", reqModel.TopUpCode)
	}

	result, err := uc.topupRepo.TopUp(ctx, reqModel)
	if err != nil {
		return nil, err
	}

	if result != nil {
		result.BalanceAfter = result.BalanceBefore + float64(reqModel.NominalTopUp)
	}

	return result, nil
}

func generateTopUpCode() string {
	return fmt.Sprintf("TOPUP-%s", time.Now().Format("20060102-150405"))
}
