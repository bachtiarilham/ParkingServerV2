package repository

import (
	"context"

	"modulegue/internal/domain/mobile/model"
)

type PembayaranRepository interface {
	GetPembayaran(ctx context.Context, sessionId int64) (*model.PembayaranModel, error)
}
