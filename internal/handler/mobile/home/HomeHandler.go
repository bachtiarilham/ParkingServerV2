package home

import (
	"net/http"

	"modulegue/core/response"
	mapper "modulegue/internal/data/mobile/remote/mapper/home"
	model "modulegue/internal/domain/mobile/model/home"
	usecase "modulegue/internal/domain/mobile/usecase/home"
	middleware "modulegue/internal/middleware"
)

type HomeHandler struct {
	getDashboardUC *usecase.GetHomeUseCase
}

func NewHomeHandler(getDashboardUC *usecase.GetHomeUseCase) *HomeHandler {
	return &HomeHandler{getDashboardUC: getDashboardUC}
}

// Endpoint: GET /api/v2/linespot/home
func (h *HomeHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	userID, okUserId := middleware.UserIDFromContext(r.Context())
	roleID, okRoleId := middleware.RoleIDFromContext(r.Context())
	if !okUserId || !okRoleId {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	input := model.GetHomeReqModel{
		UserID: userID,
		RoleID: roleID,
	}

	switch roleID {
	case 1:
		result, err := h.getDashboardUC.ExecuteCustomerHome(r.Context(), input)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "gagal memuat dashboard")
			return
		}
		response.Success(w, http.StatusOK, "Dashboard dimuat", mapper.ToCustomerHomeResponse(result))
	case 2:
		result, err := h.getDashboardUC.ExecuteJukirHome(r.Context(), input)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "gagal memuat dashboard")
			return
		}
		response.Success(w, http.StatusOK, "Dashboard dimuat", mapper.ToJukirHomeResponse(result))
	}

}
