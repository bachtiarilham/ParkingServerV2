package repository

import (
	"context"

	"modulegue/internal/domain/mobile/model"
)

type PostQrRepository interface {
	Submit(ctx context.Context, req model.SubmitQrRequestModel) (*model.PembayaranModel, error)
}
