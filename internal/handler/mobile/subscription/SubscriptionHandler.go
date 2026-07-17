package subscription

import (
	"net/http"

	"modulegue/core/response"
	mapper "modulegue/internal/data/mobile/remote/mapper/subscription"
	usecase "modulegue/internal/domain/mobile/usecase/subscription"
	middleware "modulegue/internal/middleware"
)

type SubscriptionHandler struct {
	getSubscriptionUc *usecase.SubscriptionUseCase
}

func NewSubscriptionHandler(getSubscriptionUc *usecase.SubscriptionUseCase) *SubscriptionHandler {
	return &SubscriptionHandler{getSubscriptionUc: getSubscriptionUc}
}

// Endpoint: GET /api/v2/linespot/home
func (h *SubscriptionHandler) Execute(w http.ResponseWriter, r *http.Request) {
	userID, okUserId := middleware.UserIDFromContext(r.Context())
	if !okUserId {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	result, err := h.getSubscriptionUc.Execute(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat dashboard")
		return
	}

	response.Success(w, http.StatusOK, "Dashboard dimuat", mapper.ToSubscribeDto(result))
}
