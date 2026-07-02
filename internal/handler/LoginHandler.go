package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"modulegue/core/errorstring"
	"modulegue/core/response"
	"modulegue/internal/data/mobile/remote/dto"
	"modulegue/internal/domain/mobile/usecase"
)

type LoginHandler struct {
	loginUc *usecase.LoginUseCase
}

func NewLoginHandler(
	loginUc *usecase.LoginUseCase,
) *LoginHandler {
	return &LoginHandler{
		loginUc: loginUc,
	}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	request := usecase.LoginInput{
		Identity: req.Identity,
		Password: req.Password,
	}

	result, err := h.loginUc.Execute(r.Context(), request)
	if err != nil {
		// Log error jika perlu
		if errors.Is(err, errorstring.ErrInvalidCredentials) {
			response.Error(w, http.StatusUnauthorized, "email atau password salah")
			return
		}
		// Error lainnya
		response.Error(w, http.StatusInternalServerError, "terjadi kesalahan internal")
		return
	}

	resp := dto.LoginResponseDto{
		UserDto: dto.UserDto{
			UserId:   result.UserID,
			FullName: result.FullName,
			RoleId:   result.RoleID,
		},
		TokenSetDto: dto.TokenSetDto{
			AccessToken:      result.AccessToken,
			RefreshToken:     result.RefreshToken,
			ExpiresInSeconds: result.ExpiresAt - time.Now().Unix(),
		},
	}
	response.Success(w, http.StatusOK, "login berhasil", resp)
}
