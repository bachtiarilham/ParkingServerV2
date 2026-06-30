package handler

import (
	"net/http"

	"modulegue/internal/data/mobile/remote/mapper"
	middleware "modulegue/internal/middleware"
	home "modulegue/internal/usecase/home"
	"modulegue/pkg/response"
)

type HomeHandler struct {
	getDashboardUC *home.GetDashboardUseCase
}

func NewHomeHandler(getDashboardUC *home.GetDashboardUseCase) *HomeHandler {
	return &HomeHandler{getDashboardUC: getDashboardUC}
}

// Endpoint: GET /api/v2/linespot/home
func (h *HomeHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context()) // <-- Gunakan fungsi helper
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Jika customerID ≠ userID, maka pastikan user punya akses (misalnya via region scope)
	// Untuk sementara, kita asumsikan customerID = userID
	// if customerID != userID {
	// 	response.Error(w, http.StatusForbidden, "tidak memiliki akses")
	// 	return
	// }

	input := home.GetDashboardInput{
		UserID: userID,
		Limit:  5,
		Offset: 0,
	}
	result, err := h.getDashboardUC.Execute(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat dashboard")
		return
	}

	response.Success(w, http.StatusOK, "Dashboard dimuat", mapper.ToHomeResponse(&result))
}
