package payment

import (
	"net/http"

	"modulegue/core/response"
	usecase "modulegue/internal/domain/mobile/usecase/payment_parking"
)

type GetPembayaranStatusHandler struct {
	getPembayaranStatusUc *usecase.GetPembayaranStatusUseCase
}

func NewGetPembayaranStatusHandler(getPembayaranStatusUc *usecase.GetPembayaranStatusUseCase) *GetPembayaranStatusHandler {
	return &GetPembayaranStatusHandler{getPembayaranStatusUc: getPembayaranStatusUc}
}

func (h *GetPembayaranStatusHandler) Execute(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionId")
	if sessionID == "" {
		response.Error(w, http.StatusBadRequest, "sessionId wajib diisi")
		return
	}

	result, err := h.getPembayaranStatusUc.Execute(r.Context(), sessionID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat status pembayaran")
		return
	}
	if result == "" {
		response.Error(w, http.StatusNotFound, "status pembayaran tidak ditemukan")
		return
	}

	response.Success(w, http.StatusOK, result, nil)
}
