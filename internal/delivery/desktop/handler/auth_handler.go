package handler

import (
	"encoding/json"
	"net/http"

	"modulegue/internal/delivery/mobile/customer/dto"
	authuc "modulegue/internal/usecase/auth"
	"modulegue/pkg/response"
)

type AuthHandler struct {
	registerUC *authuc.RegisterUseCase
	loginUC    *authuc.LoginUseCase
}

func NewAuthHandler(
	registerUC *authuc.RegisterUseCase,
	loginUC *authuc.LoginUseCase,
) *AuthHandler {
	return &AuthHandler{
		registerUC: registerUC,
		loginUC:    loginUC,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	err := h.registerUC.Execute(r.Context(), authuc.RegisterRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, http.StatusCreated, "registrasi berhasil", nil)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	result, err := h.loginUC.Execute(r.Context(), authuc.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(w, http.StatusOK, "login berhasil", dto.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    result.ExpiresAt,
	})
}
