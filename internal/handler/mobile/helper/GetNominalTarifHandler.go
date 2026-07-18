package helper

import (
	"net/http"

	"modulegue/core/response"
	mapper "modulegue/internal/data/mobile/remote/mapper/helper"
	usecase "modulegue/internal/domain/mobile/usecase/helper"
	middleware "modulegue/internal/middleware"
)

type GetNominalTopUpHandler struct {
	getNominalTopUpUc *usecase.GetNominalTopUpUseCase
}

func NewGetNominalTopUpHandler(
	getNominalTopUpUc *usecase.GetNominalTopUpUseCase,
) *GetNominalTopUpHandler {
	return &GetNominalTopUpHandler{
		getNominalTopUpUc: getNominalTopUpUc,
	}
}

func (h *GetNominalTopUpHandler) Execute(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	result, err := h.getNominalTopUpUc.Execute(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat nominal top up")
		return
	}

	resp := mapper.ToTopupResponseDto(result)
	if resp == nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat nominal top up")
		return
	}

	response.Success(w, http.StatusOK, "nominal top up dimuat", resp)
}
