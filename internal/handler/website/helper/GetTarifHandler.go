package helper

import (
	"encoding/json"
	"modulegue/core/response"
	dto "modulegue/internal/data/website/remote/dto/helper"
	mapper "modulegue/internal/data/website/remote/mapper/helper"
	model "modulegue/internal/domain/web/model/helper"
	uc "modulegue/internal/domain/web/usecase/helper"
	"net/http"
)

type GetTarifHandler struct {
	getTarifUc *uc.GetTarifUseCase
}

func NewGetTarifHandler(
	getTarifUc *uc.GetTarifUseCase,
) *GetTarifHandler {
	return &GetTarifHandler{
		getTarifUc: getTarifUc,
	}
}

func (h *GetTarifHandler) Execute(w http.ResponseWriter, r *http.Request) {

	var req dto.GetTarifRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "json jelek")
		return
	}

	input := model.GetTarifRequestModel{
		IDLokasi: req.IDLokasi,
	}

	resp, err := h.getTarifUc.Execute(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal mengambil tarif")
		return
	}

	response.Success(w, http.StatusOK, "Dashboard dimuat", mapper.ToGetTarifResponseDto(resp))
}
