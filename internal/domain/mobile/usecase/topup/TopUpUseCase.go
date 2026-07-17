package topup

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	model "modulegue/internal/domain/mobile/model/topup"
	"modulegue/internal/domain/mobile/repository"
)

type TopUpUseCase struct {
	topupRepo repository.TopUpRepository
	adminFee  int64
}

func NewTopUpUseCase(
	topupRepo repository.TopUpRepository,
	adminFee int64,
) *TopUpUseCase {
	return &TopUpUseCase{
		topupRepo: topupRepo,
		adminFee:  adminFee,
	}
}

func (uc *TopUpUseCase) Execute(ctx context.Context, reqModel model.TopupCreateRequestModel) (*model.TopupCreateResponseModel, error) {
	if reqModel.TopupCode == "" {
		topupCode, err := generateTopupCode()
		if err != nil {
			return nil, fmt.Errorf("generate topup code: %w", err)
		}
		reqModel.TopupCode = topupCode
	}

	if reqModel.ExternalReference == "" {
		reqModel.ExternalReference = generateExternalReference(reqModel.TopupCode)
	}

	if reqModel.ProviderName == "" {
		reqModel.ProviderName = "DEV_PROVIDER"
	}

	if reqModel.QRString == "" {
		reqModel.QRString = generateDummyQRString(reqModel.TopupCode)
	}

	reqModel.AdminFee = uc.adminFee

	if reqModel.ExpiredAt.IsZero() {
		reqModel.ExpiredAt = time.Now().Add(15 * time.Minute)
	}

	result, err := uc.topupRepo.TopUpCreate(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("get subscription: %w", err)
	}
	return result, nil
}

func generateTopupCode() (string, error) {
	suffix, err := randomHex(4)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("TOPUP-%s-%s", time.Now().UTC().Format("20060102150405"), suffix), nil
}

func generateExternalReference(topupCode string) string {
	return "EXT-" + topupCode
}

func generateDummyQRString(topupCode string) string {
	return "QRIS-DUMMY-" + topupCode
}

func randomHex(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
