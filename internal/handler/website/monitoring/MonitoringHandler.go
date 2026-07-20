package monitoring

import (
	"encoding/json"
	"modulegue/core/response"
	dto "modulegue/internal/data/website/remote/dto/monitoring"
	mapper "modulegue/internal/data/website/remote/mapper/monitoring"
	uc "modulegue/internal/domain/web/usecase/monitoring"
	"net/http"
)

type MonitoringHandler struct {
	monitoringUc *uc.MonitoringUseCase
}

func NewGetLokasiHandler(
	monitoringUc *uc.MonitoringUseCase,
) *MonitoringHandler {
	return &MonitoringHandler{
		monitoringUc: monitoringUc,
	}
}

func (h *MonitoringHandler) Execute(w http.ResponseWriter, r *http.Request) {

	var req dto.MonitoringRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "json jelek")
		return
	}

	input := mapper.ToMonitoringRequestModel(&req)
	if input == nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	resp, err := h.monitoringUc.Execute(r.Context(), *input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal mengambil lokasi")
		return
	}

	response.Success(w, http.StatusOK, "Dashboard dimuat", mapper.ToMonitoringResponseDto(resp))
}
