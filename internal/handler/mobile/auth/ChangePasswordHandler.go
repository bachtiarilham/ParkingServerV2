package auth

import (
	"encoding/json"
	"net/http"

	"modulegue/core/response"
	dto "modulegue/internal/data/mobile/remote/dto/auth"
	model "modulegue/internal/domain/mobile/model/auth"
	usecase "modulegue/internal/domain/mobile/usecase/auth"
	middleware "modulegue/internal/middleware"
)

type ChangePasswordHandler struct {
	changePasswordUc *usecase.ChangePasswordUseCase
}

func NewChangePasswordHandler(
	changePasswordUc *usecase.ChangePasswordUseCase,
) *ChangePasswordHandler {
	return &ChangePasswordHandler{
		changePasswordUc: changePasswordUc,
	}
}

func (h *ChangePasswordHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ChangePasswordRequestDto

	userID, okUserId := middleware.UserIDFromContext(r.Context())
	if !okUserId {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}

	input := model.ChangePassReqModel{
		UserId:      userID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}

	result, err := h.changePasswordUc.Execute(r.Context(), input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "logout gagal")
		return
	}

	response.Success(w, http.StatusOK, result.Message, nil)
}
