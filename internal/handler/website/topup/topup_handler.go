package topup

import (
	"encoding/json"
	"net/http"

	"modulegue/core/response"
	dto "modulegue/internal/data/website/remote/dto/topup"
	mapper "modulegue/internal/data/website/remote/mapper/topup"

	// model "modulegue/internal/domain/web/model/topup"
	uc "modulegue/internal/domain/web/usecase/topup"
)

type TopUpHandler struct {
	topupUc *uc.TopUpUseCase
}

func NewTopUpHandler(topupUc *uc.TopUpUseCase) *TopUpHandler {
	return &TopUpHandler{topupUc: topupUc}
}

func (h *TopUpHandler) Execute(w http.ResponseWriter, r *http.Request) {
	var req dto.TopUpRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	input := mapper.ToTopUpRequestModel(&req)
	if input == nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	resp, err := h.topupUc.Execute(r.Context(), *input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "gagal top up")
		return
	}

	response.Success(w, http.StatusOK, "top up berhasil", mapper.ToTopUpResponseDto(resp))
}

// func (h *TopUpHandler) Create(w http.ResponseWriter, r *http.Request) {
// 	h.Execute(w, r)
// }

// func (h *TopUpHandler) handleResponse(resp *model.TopUpResponseModel) any {
// 	return mapper.ToTopUpResponseDto(resp)
// }
