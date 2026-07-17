package payment

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"modulegue/core/response"
	dto "modulegue/internal/data/mobile/remote/dto/payment_parking"
	mapper "modulegue/internal/data/mobile/remote/mapper/payment_parking"
	model "modulegue/internal/domain/mobile/model/payment_parking"
	usecase "modulegue/internal/domain/mobile/usecase/payment_parking"
	"modulegue/internal/middleware"
)

type PostPaymentParkingHandler struct {
	postPaymentParkingUc *usecase.PostPaymentParkingUseCase
}

func NewPostPaymentParkingHandler(postPaymentParkingUc *usecase.PostPaymentParkingUseCase) *PostPaymentParkingHandler {
	return &PostPaymentParkingHandler{postPaymentParkingUc: postPaymentParkingUc}
}

func (h *PostPaymentParkingHandler) Execute(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.PostPaymentParkingRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	req.SessionCode = strings.TrimSpace(req.SessionCode)
	if req.SessionCode == "" {
		response.Error(w, http.StatusBadRequest, "session_code wajib diisi")
		return
	}

	input := mapper.ToPostPaymentParkingRequestModel(&req)
	if input == nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}
	input.CustomerUserId = userID

	result, err := h.postPaymentParkingUc.Execute(r.Context(), *input)
	if err != nil {
		if errors.Is(err, model.ErrPaymentAlreadyCompleted) {
			response.Error(w, http.StatusConflict, "pembayaran sudah dilakukan")
			return
		}
		response.Error(w, http.StatusInternalServerError, "gagal memproses pembayaran")
		return
	}

	resp := mapper.ToPostPaymentParkingResponseDto(result)
	if resp == nil {
		response.Error(w, http.StatusInternalServerError, "gagal memproses pembayaran")
		return
	}

	response.Success(w, http.StatusOK, "pembayaran diproses", resp)
}
