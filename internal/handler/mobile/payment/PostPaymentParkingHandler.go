package payment

import (
	"net/http"

	"modulegue/core/response"
	mapper "modulegue/internal/data/mobile/remote/mapper/payment"
	model "modulegue/internal/domain/mobile/model/payment"
	usecase "modulegue/internal/domain/mobile/usecase/payment"
)

type PostPaymentParkingHandler struct {
	postPaymentParkingUc *usecase.PostPaymentParkingUseCase
}

func NewPostPaymentParkingHandler(postPaymentParkingUc *usecase.PostPaymentParkingUseCase) *PostPaymentParkingHandler {
	return &PostPaymentParkingHandler{postPaymentParkingUc: postPaymentParkingUc}
}

func (h *PostPaymentParkingHandler) Execute(w http.ResponseWriter, r *http.Request) {

	result, err := h.postPaymentParkingUc.Execute(r.Context(), model.PostPaymentParkingRequestModel{})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat dashboard")
		return
	}

	response.Success(w, http.StatusOK, "Dashboard dimuat", mapper.ToPostPaymentParkingResponseDto(result))
}
