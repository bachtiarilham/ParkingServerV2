package paymentgate

import (
	"context"
	"fmt"
	"strings"

	model "modulegue/internal/domain/mobile/model/payment_gate"
	"modulegue/internal/domain/mobile/repository"
)

type SignatureVerifier interface {
	VerifySignature(orderID, statusCode, grossAmount, signature string) bool
}

type PaymentCallbackUseCase struct {
	repo     repository.PaymentCallbackRepository
	reposync repository.SyncRepo
	verifier SignatureVerifier
}

func NewPaymentCallbackUseCase(
	repo repository.PaymentCallbackRepository,
	reposync repository.SyncRepo,
	verifier SignatureVerifier,
) *PaymentCallbackUseCase {
	return &PaymentCallbackUseCase{
		repo:     repo,
		reposync: reposync,
		verifier: verifier,
	}
}

func (uc *PaymentCallbackUseCase) Execute(ctx context.Context, reqModel model.CallbackRequestModel) error {
	// 1. Verify Midtrans signature to prevent spoofing
	// Gross amount may have decimal parts, e.g. "2500.00", we need to strip ".00" if present or normalize
	grossAmount := reqModel.GrossAmount
	if idx := strings.Index(grossAmount, "."); idx != -1 {
		grossAmount = grossAmount[:idx]
	}

	isValid := uc.verifier.VerifySignature(reqModel.OrderID, reqModel.StatusCode, reqModel.GrossAmount, reqModel.SignatureKey)
	// isValid := true
	if !isValid {
		return fmt.Errorf("invalid signature key: verification failed")
	}

	// 2. We only process status changes on SUCCESS / SETTLEMENT / CAPTURE
	// DENY, EXPIRE, CANCEL can be handled accordingly, but only settlement/capture means lunas
	status := strings.ToLower(reqModel.TransactionStatus)
	if status != "settlement" && status != "capture" {
		// Log or return, it is PENDING or FAILED but we return nil so Midtrans doesn't retry
		return nil
	}

	// 3. Lookup the transaction info from database
	txInfo, err := uc.repo.GetPaymentTransaction(ctx, reqModel.OrderID)
	if err != nil {
		return fmt.Errorf("get transaction details: %w", err)
	}

	// 4. Route to specific domain callback processing
	switch txInfo.PaymentType {
	case "PARKING":
		if err := uc.repo.ProcessParkingCallback(ctx, reqModel.OrderID, reqModel.TransactionID, txInfo.ReferenceID, model.PaymentTransactionModel{
			PaymentType: txInfo.PaymentType,
			UserID:      txInfo.UserID,
			ReferenceID: txInfo.ReferenceID,
			Amount:      txInfo.Amount,
		}); err != nil {
			return fmt.Errorf("process parking callback: %w", err)
		}

		if err := uc.reposync.SyncAfterParkir(ctx, txInfo.UserID, txInfo.ReferenceID, txInfo.Amount); err != nil {
			return fmt.Errorf("process sync parkir: %w", err)
		}

	case "TOPUP":
		if err := uc.repo.ProcessTopupCallback(ctx, reqModel.OrderID, reqModel.TransactionID, txInfo.UserID, txInfo.Amount); err != nil {
			return fmt.Errorf("process topup callback: %w", err)
		}
	case "MEMBERSHIP":
		if err := uc.repo.ProcessMembershipCallback(ctx, reqModel.OrderID, reqModel.TransactionID, txInfo.UserID, txInfo.ReferenceID); err != nil {
			return fmt.Errorf("process membership callback: %w", err)
		}

		if err := uc.reposync.SyncAfterMembership(ctx, txInfo.UserID, txInfo.ReferenceID); err != nil {
			return fmt.Errorf("process sync membership: %w", err)
		}
	default:
		return fmt.Errorf("unsupported payment type in callback: %s", txInfo.PaymentType)
	}

	return nil
}
