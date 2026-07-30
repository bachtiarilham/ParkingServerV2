package paymentgate

import (
	"context"
	"fmt"
	"strings"

	model "modulegue/internal/domain/mobile/model/payment_gate"
	"modulegue/internal/domain/mobile/repository"
)

type GetPaymentStatusUseCase struct {
	repo repository.PaymentCallbackRepository
}

func NewGetPaymentStatusUseCase(
	repo repository.PaymentCallbackRepository,
) *GetPaymentStatusUseCase {
	return &GetPaymentStatusUseCase{
		repo: repo,
	}
}

func (uc *GetPaymentStatusUseCase) Execute(ctx context.Context, txCode string) (*model.PaymentStatusModel, error) {
	var status string
	var err error

	// If transaction code starts with SESS (case-insensitive), retrieve status for Cash parking session.
	// Otherwise, retrieve payment status using transaction code (TRX-).
	if strings.HasPrefix(strings.ToUpper(txCode), "SESS") {
		status, err = uc.repo.GetPaymentStatusCash(ctx, txCode)
		if err != nil {
			return nil, fmt.Errorf("get cash payment status: %w", err)
		}
	} else {
		status, err = uc.repo.GetPaymentStatus(ctx, txCode)
		if err != nil {
			return nil, fmt.Errorf("get midtrans payment status: %w", err)
		}
	}

	var statusText string
	switch status {
	case "SUCCESS", "3":
		statusText = "Lunas"
	case "PENDING", "1":
		statusText = "Menunggu Pembayaran"
	default:
		statusText = "Gagal"
	}

	return &model.PaymentStatusModel{
		OrderID:    txCode,
		Status:     status,
		StatusText: statusText,
	}, nil
}
