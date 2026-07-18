package repository

import (
	"context"
	model "modulegue/internal/domain/web/model/helper"
)

type HelperRepository interface {
	GetLokasi(context.Context, model.GetLokasiRequestModel) (*model.GetLokasiResponseModel, error)
	GetTarif(context.Context, model.GetTarifRequestModel) (*model.GetTarifResponseModel, error)
}
