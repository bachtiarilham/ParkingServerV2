package repository

import (
	"context"
	"modulegue/internal/delivery/web/dto/request"
	"modulegue/internal/domain/web/finance"
	"modulegue/internal/domain/web/transaction"
)

type transactionRepo interface {
	CreateDisputeCase(ctx context.Context, adminID int64, req request.CreateDisputeCaseRequest) (transaction.DisputeCaseSummary, error)
	UpdateDisputeCaseStatus(ctx context.Context, adminID int64, disputeID string, req request.UpdateDisputeCaseStatusRequest) (transaction.DisputeCaseSummary, error)
	CreateRefundTransaction(ctx context.Context, adminID int64, req request.CreateRefundTransactionRequest) (finance.RefundTransactionSummary, error)
	UpdateRefundTransactionStatus(ctx context.Context, adminID int64, refundID string, req request.UpdateRefundStatusRequest) (finance.RefundTransactionSummary, error)
	CreateClosingBatch(ctx context.Context, adminID int64, req request.CreateClosingBatchRequest) (finance.ClosingBatchSummary, error)
	UpdateClosingBatchStatus(ctx context.Context, adminID int64, closingID string, req request.UpdateClosingStatusRequest) (finance.ClosingBatchSummary, error)
}
