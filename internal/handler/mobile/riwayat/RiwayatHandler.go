package riwayat

import (
	"net/http"

	"modulegue/core/response"
	mapper "modulegue/internal/data/mobile/remote/mapper/riwayat"
	model "modulegue/internal/domain/mobile/model/riwayat"
	usecase "modulegue/internal/domain/mobile/usecase/riwayat"
	middleware "modulegue/internal/middleware"
)

type RiwayatHandler struct {
	getRiwayatUc *usecase.GetRiwayatUseCase
}

func NewRiwayatHandler(getRiwayatUc *usecase.GetRiwayatUseCase) *RiwayatHandler {
	return &RiwayatHandler{getRiwayatUc: getRiwayatUc}
}

// Endpoint: GET /api/v2/linespot/home
func (h *RiwayatHandler) Execute(w http.ResponseWriter, r *http.Request) {
	userID, okUserId := middleware.UserIDFromContext(r.Context())
	roleID, okRoleId := middleware.RoleIDFromContext(r.Context())
	if !okUserId || !okRoleId {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	input := model.RiwayatRequestModel{
		UserID:    userID,
		RoleID:    roleID,
		Username:  r.URL.Query().Get("username"),
		StartDate: r.URL.Query().Get("startDate"),
		EndDate:   r.URL.Query().Get("endDate"),
		Payment:   r.URL.Query().Get("payment"),
		Vehicle:   r.URL.Query().Get("vehicle"),
		Lokasi:    r.URL.Query().Get("lokasi"),
	}

	result, err := h.getRiwayatUc.Execute(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat dashboard")
		return
	}

	response.Success(w, http.StatusOK, "Dashboard dimuat", mapper.ToRiwayatDto(result))
}
