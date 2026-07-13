package payment

import (
	"context"
	"fmt"
	repository "modulegue/internal/domain/mobile/repository"
)

type GetPembayaranStatusUseCase struct {
	getPembayaranStatusRepo repository.StatusPaymentRepository
}

func NewGetPembayaranStatusUseCase(
	getPembayaranStatusRepo repository.StatusPaymentRepository,
) *GetPembayaranStatusUseCase {
	return &GetPembayaranStatusUseCase{
		getPembayaranStatusRepo: getPembayaranStatusRepo,
	}
}

func (uc *GetPembayaranStatusUseCase) Execute(ctx context.Context, sessionId string) (string, error) {
	result, err := uc.getPembayaranStatusRepo.GetPembayaranStatus(ctx, sessionId)
	if err != nil {
		return "", fmt.Errorf("get pembayaran status: %w", err)
	}
	if result != nil {
		return "pembayaran berhasil", nil
	}

	return "", nil
}
