package repository

import (
	"context"

	model "modulegue/internal/domain/mobile/model/invoice"
)

type InvoiceRepository interface {
	GetInvoice(ctx context.Context, req string) (*model.UniversalInvoiceResponseModel, error)
}
