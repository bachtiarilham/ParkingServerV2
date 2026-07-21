package riwayat

import (
	"encoding/json"
	"net/http"

	"modulegue/core/response"
	"modulegue/core/utils"
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
	if !okUserId {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.RiwayatRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	startDate, err := utils.ParseISODate(req.StartDate)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "startDate tidak valid")
		return
	}

	endDate, err := utils.ParseISODate(req.EndDate)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "endDate tidak valid")
		return
	}

	input := model.RiwayatRequestModel{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	result, err := h.getRiwayatUc.Execute(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat riwayat")
		return
	}

	response.Success(w, http.StatusOK, "Riwayat dimuat", mapper.ToRiwayatDto(result))
}
