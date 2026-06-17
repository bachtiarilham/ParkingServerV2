package dto

import (
	"modulegue/internal/domain/web/metrics"
	"modulegue/internal/domain/web/settings"
	"modulegue/internal/domain/web/transaction"
)

type TransactionsOverview struct {
	TransactionStats      []metrics.StatCard                 `json:"transactionStats"`
	TransactionRows       []metrics.RowItem                  `json:"transactionRows"`
	PaymentBreakdownItems []transaction.PaymentBreakdownItem `json:"paymentBreakdownItems"`
	TransactionIssueItems []transaction.TransactionIssueItem `json:"transactionIssueItems"`
	ExportQueueItems      []settings.ExportQueueItem         `json:"exportQueueItems"`
}
