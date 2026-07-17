package topup

import (
	"net/http"

	"modulegue/core/response"
	mapper "modulegue/internal/data/mobile/remote/mapper/topup"
	usecase "modulegue/internal/domain/mobile/usecase/topup"
	middleware "modulegue/internal/middleware"
)

type GetTopUpStatusHandler struct {
	getTopUpStatusUc *usecase.GetTopUpStatusUseCase
}

func NewGetTopUpStatusHandler(getTopUpStatusUc *usecase.GetTopUpStatusUseCase) *GetTopUpStatusHandler {
	return &GetTopUpStatusHandler{getTopUpStatusUc: getTopUpStatusUc}
}

func (h *GetTopUpStatusHandler) Execute(w http.ResponseWriter, r *http.Request) {
	_, okUserId := middleware.UserIDFromContext(r.Context())
	if !okUserId {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	topUpCode := r.PathValue("topupCode")
	if topUpCode == "" {
		response.Error(w, http.StatusBadRequest, "topupcode wajib diisi")
		return
	}

	result, err := h.getTopUpStatusUc.Execute(r.Context(), topUpCode)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat status pembayaran")
		return
	}

	response.Success(w, http.StatusOK, "hasil status", mapper.ToTopupStatusResponseDto(result))
}
