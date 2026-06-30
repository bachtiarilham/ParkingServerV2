package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"modulegue/internal/delivery/web/dto"
	"modulegue/internal/delivery/web/mapper"
	"modulegue/internal/domain/dashboard"
	"modulegue/internal/usecase/web"
	"modulegue/pkg/response"
)

type AuthHandler struct {
	authUC *web.AuthUseCase
}

func NewAuthHandler(authUC *web.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}
	if req.Identity == "" || req.Password == "" {
		response.Error(w, http.StatusBadRequest, "identity dan password wajib diisi")
		return
	}

	result, err := h.authUC.Login(r.Context(), dashboard.LoginRequest{
		Identity: req.Identity,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, web.ErrInvalidCredentials):
			response.Error(w, http.StatusUnauthorized, "kredensial tidak valid")
		case errors.Is(err, web.ErrForbiddenRole):
			response.Error(w, http.StatusForbidden, "role tidak diizinkan")
		default:
			response.Error(w, http.StatusInternalServerError, "gagal login dashboard")
		}
		return
	}

	response.Success(w, http.StatusOK, "login berhasil", mapper.FromAuthEnvelope(result))
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	result, err := h.authUC.Refresh(r.Context(), dashboard.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		switch {
		case errors.Is(err, web.ErrInvalidRefreshToken):
			response.Error(w, http.StatusUnauthorized, "refresh token tidak valid")
		case errors.Is(err, web.ErrExpiredRefreshToken):
			response.Error(w, http.StatusUnauthorized, "refresh token sudah kadaluarsa")
		case errors.Is(err, web.ErrForbiddenRole):
			response.Error(w, http.StatusForbidden, "role tidak diizinkan")
		default:
			response.Error(w, http.StatusInternalServerError, "gagal memperbarui token")
		}
		return
	}

	response.Success(w, http.StatusOK, "token diperbarui", mapper.FromAuthEnvelope(result))
}
