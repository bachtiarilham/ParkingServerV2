package helper

import (
	"net/http"

	"modulegue/core/response"
	mapper "modulegue/internal/data/mobile/remote/mapper/helper"
	usecase "modulegue/internal/domain/mobile/usecase/helper"
	middleware "modulegue/internal/middleware"
)

type GetPaymentMethodHandler struct {
	getPaymentMethodUc *usecase.GetPaymentMethodUseCase
}

func NewGetPaymentMethodHandler(
	getPaymentMethodUc *usecase.GetPaymentMethodUseCase,
) *GetPaymentMethodHandler {
	return &GetPaymentMethodHandler{
		getPaymentMethodUc: getPaymentMethodUc,
	}
}

func (h *GetPaymentMethodHandler) Execute(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	result, err := h.getPaymentMethodUc.Execute(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat nominal top up")
		return
	}

	resp := mapper.ToPaymentMethodDto(result)
	if resp == nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat nominal top up")
		return
	}

	response.Success(w, http.StatusOK, "nominal top up dimuat", resp)
}
