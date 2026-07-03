package helper

import (
	"net/http"

	"modulegue/core/response"
	mapper "modulegue/internal/data/mobile/remote/mapper/helper"
	usecase "modulegue/internal/domain/mobile/usecase/helper"
)

type GetLocationHandler struct {
	getLocationUc *usecase.GetLokasiUseCase
}

func NewGetLocationHandler(
	getLocationUc *usecase.GetLokasiUseCase,
) *GetLocationHandler {
	return &GetLocationHandler{
		getLocationUc: getLocationUc,
	}
}

func (h *GetLocationHandler) Execute(w http.ResponseWriter, r *http.Request) {
	result, err := h.getLocationUc.Execute(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat lokasi")
		return
	}

	resp := mapper.ToLokasiDto(result)
	if resp == nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat lokasi")
		return
	}

	response.Success(w, http.StatusOK, "lokasi dimuat", resp)
}
