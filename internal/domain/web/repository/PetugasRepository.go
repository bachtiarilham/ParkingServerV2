package repository

import (
	"context"
	model "modulegue/internal/domain/web/model/petugas"
)

type PetugasRepository interface {
	GetPetugas(context.Context, model.PetugasRequestModel) (*model.PetugasResponseModel, error)
}
