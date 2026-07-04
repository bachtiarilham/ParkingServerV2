package payment

import (
	"net/http"

	"modulegue/core/response"
	mapper "modulegue/internal/data/mobile/remote/mapper/payment"
	model "modulegue/internal/domain/mobile/model/payment"
	usecase "modulegue/internal/domain/mobile/usecase/payment"
)

type PostParkingHandler struct {
	postParkingUc *usecase.PostParkingUseCase
}

func NewPostParkingHandler(postParkingUc *usecase.PostParkingUseCase) *PostParkingHandler {
	return &PostParkingHandler{postParkingUc: postParkingUc}
}

func (h *PostParkingHandler) Execute(w http.ResponseWriter, r *http.Request) {

	result, err := h.postParkingUc.Execute(r.Context(), model.PostParkingRequestModel{})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat dashboard")
		return
	}

	response.Success(w, http.StatusOK, "Dashboard dimuat", mapper.ToPembayaranDto(result))
}
