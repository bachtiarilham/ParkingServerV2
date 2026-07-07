package riwayat

import (
	"encoding/json"
	"net/http"
	"strings"

	"modulegue/core/response"
	dto "modulegue/internal/data/mobile/remote/dto/riwayat"
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

// Endpoint: POST /api/v2/linespot/riwayat
func (h *RiwayatHandler) Execute(w http.ResponseWriter, r *http.Request) {
	userID, okUserId := middleware.UserIDFromContext(r.Context())
	roleID, okRoleId := middleware.RoleIDFromContext(r.Context())
	if !okUserId || !okRoleId {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.RiwayatRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	input := model.RiwayatRequestModel{
		UserID:    userID,
		RoleID:    roleID,
		Username:  strings.TrimSpace(req.Username),
		StartDate: strings.TrimSpace(req.StartDate),
		EndDate:   strings.TrimSpace(req.EndDate),
		Payment:   strings.TrimSpace(req.Payment),
		Vehicle:   strings.TrimSpace(req.Vehicle),
		Lokasi:    strings.TrimSpace(req.Lokasi),
	}

	result, err := h.getRiwayatUc.Execute(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat riwayat")
		return
	}

	response.Success(w, http.StatusOK, "Riwayat dimuat", mapper.ToRiwayatDto(result))
}
