package laporan

import (
	"encoding/json"
	"net/http"

	"modulegue/core/response"
	"modulegue/core/utils"
	dto "modulegue/internal/data/mobile/remote/dto/laporan"
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

func (h *LaporanHandler) Execute(w http.ResponseWriter, r *http.Request) {
	userID, okUserId := middleware.UserIDFromContext(r.Context())
	if !okUserId {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var inputDto dto.LaporanFilterRequestDto
	if err := json.NewDecoder(r.Body).Decode(&inputDto); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	startDate, err := utils.ParseISODate(inputDto.StartDate)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "startDate tidak valid")
		return
	}

	endDate, err := utils.ParseISODate(inputDto.EndDate)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "endDate tidak valid")
		return
	}

	input := model.LaporanRequestModel{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	result, err := h.getLaporanUc.Execute(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat laporan")
		return
	}

	response.Success(w, http.StatusOK, "Laporan dimuat", mapper.ToLaporanDto(result))
}
