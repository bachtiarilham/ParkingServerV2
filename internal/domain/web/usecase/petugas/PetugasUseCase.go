package petugas

import (
	"context"
	"fmt"
	model "modulegue/internal/domain/web/model/petugas"
	"modulegue/internal/domain/web/repository"
)

type PetugasUseCase struct {
	petugasRepo repository.PetugasRepository
}

func NewPetugasUseCase(
	petugasRepo repository.PetugasRepository,
) *PetugasUseCase {
	return &PetugasUseCase{
		petugasRepo: petugasRepo,
	}
}

func (uc *PetugasUseCase) Execute(ctx context.Context, reqModel model.PetugasRequestModel) (*model.PetugasResponseModel, error) {
	resp, err := uc.petugasRepo.GetPetugas(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data petugas : %w", err)
	}
	return resp, nil
}
