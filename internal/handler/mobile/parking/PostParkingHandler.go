package parking

import (
	"encoding/json"
	"net/http"
	"strings"

	"modulegue/core/response"
	dto "modulegue/internal/data/mobile/remote/dto/parking"
	mapper "modulegue/internal/data/mobile/remote/mapper/parking"
	model "modulegue/internal/domain/mobile/model/parking"
	usecase "modulegue/internal/domain/mobile/usecase/parking"
	"modulegue/internal/middleware"
)

type PostParkingHandler struct {
	postParkingUc *usecase.PostParkingUseCase
}

func NewPostParkingHandler(postParkingUc *usecase.PostParkingUseCase) *PostParkingHandler {
	return &PostParkingHandler{postParkingUc: postParkingUc}
}

func (h *PostParkingHandler) Execute(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.PostParkingRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	req.PlateNumber = strings.TrimSpace(req.PlateNumber)
	req.VehicleTypeCode = strings.TrimSpace(req.VehicleTypeCode)

	if req.PlateNumber == "" || req.VehicleTypeCode == "" {
		response.Error(w, http.StatusBadRequest, "nomor polisi dan jenis kendaraan wajib diisi")
		return
	}

	input := model.PostParkingRequestModel{
		OfficerUserId:   userID,
		PlateNumber:     req.PlateNumber,
		VehicleTypeCode: req.VehicleTypeCode,
		SelectedAreaId:  req.SelectedAreaId,
		BiayaParkir:     req.BiayaParkir,
	}

	result, err := h.postParkingUc.Execute(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memproses parking")
		return
	}

	response.Success(w, http.StatusOK, "parking berhasil", mapper.ToParkingResponseDto(result))
}
