package riwayat

import (
	"net/http"

	"modulegue/core/response"
	dto "modulegue/internal/data/mobile/remote/dto/riwayat"
	mapper "modulegue/internal/data/mobile/remote/mapper/riwayat"
	usecase "modulegue/internal/domain/mobile/usecase/riwayat"
	middleware "modulegue/internal/middleware"
)

type TransaksiDetilHandler struct {
	getTransaksiDetilUc *usecase.GetTransaksiDetilUseCase
}

func NewTransaksiDetilHandler(getTransaksiDetilUc *usecase.GetTransaksiDetilUseCase) *TransaksiDetilHandler {
	return &TransaksiDetilHandler{getTransaksiDetilUc: getTransaksiDetilUc}
}

func (h *TransaksiDetilHandler) Execute(w http.ResponseWriter, r *http.Request) {
	_, okUserId := middleware.UserIDFromContext(r.Context())
	if !okUserId {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	topup_code := r.PathValue("topup_code")
	if topup_code == "" {
		response.Error(w, http.StatusBadRequest, "topup_code wajib diisi")
		return
	}

	dtoReq := &dto.DetilTransaksiRequestDto{
		TopUpCode: topup_code,
	}

	input := mapper.ToDetilTransaksiRequestModel(dtoReq)
	if input == nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	result, err := h.getTransaksiDetilUc.Execute(r.Context(), *input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal memuat status pembayaran")
		return
	}

	response.Success(w, http.StatusOK, "hasil status", mapper.ToDetilTransaksiDto(result))
}
