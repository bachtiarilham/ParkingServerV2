package payment

import (
	"context"
	// "errors"
	"fmt"
	"modulegue/internal/domain/payment"
)

type SubmitQrInput struct {
	ScannedQrString string
	CustomerID      int64 // Didapat dari context JWT
}

type SubmitQrOutput struct {
	SessionID int64
	Message   string
}

type SubmitQrUseCase struct {
	repo payment.Repository
}

func NewSubmitQrUseCase(repo payment.Repository) *SubmitQrUseCase {
	return &SubmitQrUseCase{repo: repo}
}

func (uc *SubmitQrUseCase) Execute(ctx context.Context, input SubmitQrInput) (SubmitQrOutput, error) {
	// 1. Cari sesi parkir aktif berdasarkan QR code
	session, err := uc.repo.GetActiveSessionByCode(ctx, input.ScannedQrString)
	if err != nil {
		return SubmitQrOutput{}, fmt.Errorf("failed to find active session: %w", err)
	}

	// 2. Validasi apakah sesi ini milik customer (opsional, tergantung kebijakan)
	if session.CustomerID != input.CustomerID {
		// Jika strict, return error. Jika tidak, lanjutkan.
		// return SubmitQrOutput{}, errors.New("session does not belong to this customer")
	}

	return SubmitQrOutput{
		SessionID: session.ID,
		Message:   "QR scanned successfully, session details retrieved.",
	}, nil
}
