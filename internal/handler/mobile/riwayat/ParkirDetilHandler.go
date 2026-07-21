package riwayat

import (
	"net/http"

	"modulegue/core/response"
	dto "modulegue/internal/data/mobile/remote/dto/riwayat"
	mapper "modulegue/internal/data/mobile/remote/mapper/riwayat"
	usecase "modulegue/internal/domain/mobile/usecase/riwayat"
	middleware "modulegue/internal/middleware"
)

type ParkirDetilHandler struct {
	getParkirDetilUc *usecase.GetParkirDetilUseCase
}

func NewParkirDetilHandler(getParkirDetilUc *usecase.GetParkirDetilUseCase) *ParkirDetilHandler {
	return &ParkirDetilHandler{getParkirDetilUc: getParkirDetilUc}
}

func (h *ParkirDetilHandler) Execute(w http.ResponseWriter, r *http.Request) {
	_, okUserId := middleware.UserIDFromContext(r.Context())
	if !okUserId {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	transaction_code := r.PathValue("transaction_code")
	if transaction_code == "" {
		response.Error(w, http.StatusBadRequest, "transaction_code wajib diisi")
		return
	}

	dtoReq := &dto.DetilParkirRequestDto{
		TransactionCode: transaction_code,
	}

	input := mapper.ToDetilParkirRequestModel(dtoReq)
	if input == nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	result, err := h.getParkirDetilUc.Execute(r.Context(), *input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat status pembayaran")
		return
	}

	response.Success(w, http.StatusOK, "hasil status", mapper.ToDetilParkirDto(result))
}
