package topup

import (
	"encoding/json"
	"net/http"
	"strings"

	"modulegue/core/response"
	dto "modulegue/internal/data/mobile/remote/dto/topup"
	mapper "modulegue/internal/data/mobile/remote/mapper/topup"
	usecase "modulegue/internal/domain/mobile/usecase/topup"
	middleware "modulegue/internal/middleware"
)

type TopUpHandler struct {
	topUpUc *usecase.TopUpUseCase
}

func NewSubscriptionHandler(topUpUc *usecase.TopUpUseCase) *TopUpHandler {
	return &TopUpHandler{topUpUc: topUpUc}
}

// Endpoint: POST /api/v2/linespot/topup/create
func (h *TopUpHandler) Execute(w http.ResponseWriter, r *http.Request) {
	userID, okUserId := middleware.UserIDFromContext(r.Context())
	if !okUserId {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.TopupCreateRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	req.PaymentMethodCode = strings.TrimSpace(req.PaymentMethodCode)

	input := mapper.ToTopupCreateRequestModel(&req)
	if input == nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}
	input.UserID = userID

	result, err := h.topUpUc.Execute(r.Context(), *input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memproses topup")
		return
	}

	resp := mapper.ToTopupCreateResponseDto(result)
	if resp == nil {
		response.Error(w, http.StatusInternalServerError, "gagal memproses topup")
		return
	}

	response.Success(w, http.StatusOK, "topup diproses", resp)
}
