package petugas

import (
	"encoding/json"
	"net/http"

	"modulegue/core/response"
	dto "modulegue/internal/data/website/remote/dto/petugas"
	mapper "modulegue/internal/data/website/remote/mapper/petugas"
	model "modulegue/internal/domain/web/model/petugas"
	uc "modulegue/internal/domain/web/usecase/petugas"
)

type PetugasHandler struct {
	petugasUc *uc.PetugasUseCase
}

func NewPetugasHandler(petugasUc *uc.PetugasUseCase) *PetugasHandler {
	return &PetugasHandler{petugasUc: petugasUc}
}

func (h *PetugasHandler) GetPetugas(w http.ResponseWriter, r *http.Request) {
	var req dto.PetugasRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	input := model.PetugasRequestModel{
		IDLokasi: req.IDLokasi,
	}

	resp, err := h.petugasUc.Execute(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal mengambil data petugas")
		return
	}

	response.Success(w, http.StatusOK, "data petugas berhasil dimuat", mapper.ToPetugasResponseDto(resp))
}
