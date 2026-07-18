package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"modulegue/core/errorstring"
	"modulegue/core/response"

	dto "modulegue/internal/data/shared/remote/dto/auth"
	model "modulegue/internal/domain/shared/model/auth"
	usecase "modulegue/internal/domain/shared/usecase/auth"
)

type RefreshTokenHandler struct {
	refreshTokenUc *usecase.RefreshTokenUseCase
}

func NewRefreshTokenHandler(
	refreshTokenUc *usecase.RefreshTokenUseCase,
) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		refreshTokenUc: refreshTokenUc,
	}
}

func (h *RefreshTokenHandler) Execute(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	requestModel := model.RefreshTokenReqModel{
		RefreshToken: req.RefreshToken,
	}

	result, err := h.refreshTokenUc.Execute(r.Context(), requestModel)
	if err != nil {
		switch {
		case errors.Is(err, errorstring.ErrInvalidRefreshToken):
			response.Error(w, http.StatusUnauthorized, "refresh token tidak valid")
		case errors.Is(err, errorstring.ErrExpiredRefreshToken):
			response.Error(w, http.StatusUnauthorized, "refresh token sudah kadaluarsa")
		case errors.Is(err, errorstring.ErrUserNotFound):
			response.Error(w, http.StatusUnauthorized, "akun tidak ditemukan atau tidak aktif")
		default:
			response.Error(w, http.StatusInternalServerError, "terjadi kesalahan internal")
		}
		return
	}

	resp := dto.TokenSetDto{
		AccessToken:      result.AccessToken,
		RefreshToken:     result.RefreshToken,
		ExpiresInSeconds: result.ExpiresAt,
	}

	response.Success(w, http.StatusOK, "token diperbarui", resp)
}
