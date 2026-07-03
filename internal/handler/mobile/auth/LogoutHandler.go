package auth

import (
	"encoding/json"
	"net/http"

	"modulegue/core/response"
	dto "modulegue/internal/data/mobile/remote/dto/auth"
	model "modulegue/internal/domain/mobile/model/auth"
	usecase "modulegue/internal/domain/mobile/usecase/auth"
)

type LogoutHandler struct {
	logoutUc *usecase.LogoutUseCase
}

func NewLogoutHandler(
	logoutUc *usecase.LogoutUseCase,
) *LogoutHandler {
	return &LogoutHandler{
		logoutUc: logoutUc,
	}
}

func (h *LogoutHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.LogoutRequestDto
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}

	input := model.LogoutReqModel{
		RefreshToken: req.RefreshToken,
	}

	result, err := h.logoutUc.Execute(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "logout gagal")
		return
	}

	response.Success(w, http.StatusOK, result.Message, nil)
}
