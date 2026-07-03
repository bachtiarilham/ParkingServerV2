package laporan

import (
	"net/http"

	"modulegue/core/response"
	mapper "modulegue/internal/data/mobile/remote/mapper/laporan"
	model "modulegue/internal/domain/mobile/model/laporan"
	usecase "modulegue/internal/domain/mobile/usecase/laporan"
	middleware "modulegue/internal/middleware"
)

type LaporanHandler struct {
	getLaporanUc *usecase.GetLaporanUseCase
}

func NewLaporanHandler(getLaporanUc *usecase.GetLaporanUseCase) *LaporanHandler {
	return &LaporanHandler{getLaporanUc: getLaporanUc}
}

// Endpoint: GET /api/v2/linespot/home
func (h *LaporanHandler) Execute(w http.ResponseWriter, r *http.Request) {
	userID, okUserId := middleware.UserIDFromContext(r.Context())
	roleID, okRoleId := middleware.RoleIDFromContext(r.Context())
	if !okUserId || !okRoleId {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	input := model.LaporanRequestModel{
		UserID:    userID,
		RoleID:    roleID,
		Username:  r.URL.Query().Get("username"),
		StartDate: r.URL.Query().Get("startDate"),
		EndDate:   r.URL.Query().Get("endDate"),
		Lokasi:    r.URL.Query().Get("lokasi"),
	}

	result, err := h.getLaporanUc.Execute(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat dashboard")
		return
	}

	response.Success(w, http.StatusOK, "Dashboard dimuat", mapper.ToLaporanDto(result))
}
