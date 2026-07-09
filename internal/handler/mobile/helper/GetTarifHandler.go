package helper

import (
	"net/http"

	"modulegue/core/response"
	mapper "modulegue/internal/data/mobile/remote/mapper/helper"
	usecase "modulegue/internal/domain/mobile/usecase/helper"
	"modulegue/internal/middleware"
)

type GetTarifHandler struct {
	getTarifUc *usecase.GetTarifUseCase
}

func NewGetTarifHandler(
	getTarifUc *usecase.GetTarifUseCase,
) *GetTarifHandler {
	return &GetTarifHandler{
		getTarifUc: getTarifUc,
	}
}

func (h *GetTarifHandler) Execute(w http.ResponseWriter, r *http.Request) {
	userId, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	result, err := h.getTarifUc.Execute(r.Context(), userId)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat tarif")
		return
	}

	resp := mapper.ToTarifDto(result)
	if resp == nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat tarif")
		return
	}

	response.Success(w, http.StatusOK, "tarif dimuat", resp)
}
