package helper

import (
	"encoding/json"
	"modulegue/core/response"
	dto "modulegue/internal/data/website/remote/dto/helper"
	mapper "modulegue/internal/data/website/remote/mapper/helper"
	uc "modulegue/internal/domain/web/usecase/helper"
	"net/http"
)

type GetLokasiHandler struct {
	getLokasiUc *uc.GetLokasiUseCase
}

func NewGetLokasiHandler(
	getLokasiUc *uc.GetLokasiUseCase,
) *GetLokasiHandler {
	return &GetLokasiHandler{
		getLokasiUc: getLokasiUc,
	}
}

func (h *GetLokasiHandler) Execute(w http.ResponseWriter, r *http.Request) {

	var req dto.GetLokasiRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "json jelek")
		return
	}

	input := mapper.ToGetLokasiRequestModel(&req)
	if input == nil {
		response.Error(w, http.StatusBadRequest, "json jelek")
		return
	}

	resp, err := h.getLokasiUc.Execute(r.Context(), *input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal mengambil lokasi")
		return
	}

	response.Success(w, http.StatusOK, "Dashboard dimuat", mapper.ToGetLokasiResponseDto(resp))
}
