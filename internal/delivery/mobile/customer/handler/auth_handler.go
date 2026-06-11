package handler

import (
	"encoding/json"
	"errors"

	// "fmt"
	"net/http"
	"time"

	"modulegue/internal/delivery/mobile/customer/dto"
	auth "modulegue/internal/domain/auth"
	authuc "modulegue/internal/usecase/auth"

	"modulegue/pkg/response"
)

type AuthHandler struct {
	registerUC *authuc.RegisterUseCase
	loginUC    *authuc.LoginUseCase
	refreshUC  *authuc.RefreshTokenUseCase
}

func NewAuthHandler(
	registerUC *authuc.RegisterUseCase,
	loginUC *authuc.LoginUseCase,
	refreshUC *authuc.RefreshTokenUseCase,
) *AuthHandler {
	return &AuthHandler{
		registerUC: registerUC,
		loginUC:    loginUC,
		refreshUC:  refreshUC,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	// --- Mapping: DTO -> UseCase Input ---
	input := authuc.RegisterInput{
		FullName: req.FullName,
		Nik:      req.NIK,
		Phone:    req.PhoneNumber,
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	// --- Panggil UseCase ---
	result, err := h.registerUC.Execute(r.Context(), input) // Asumsikan registerUC adalah *auth.RegisterUseCase
	if err != nil {
		// Log error jika perlu
		switch {
		case errors.Is(err, authuc.ErrEmailAlreadyExists):
			response.Error(w, http.StatusConflict, "email sudah terdaftar")
		case errors.Is(err, authuc.ErrUsernameAlreadyExists):
			response.Error(w, http.StatusConflict, "username sudah digunakan")
		default:
			// Error lainnya
			response.Error(w, http.StatusInternalServerError, "terjadi kesalahan internal")
		}
		return
	}

	// --- Mapping: UseCase Output -> DTO ---
	resp := dto.RegisterResponse{
		Message: result.Message,
		UserID:  result.UserID,
	}

	response.Success(w, http.StatusCreated, "registrasi berhasil", resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	request := authuc.LoginRequest{
		Identity: req.Identity,
		Password: req.Password,
	}

	result, err := h.loginUC.Execute(r.Context(), request)
	if err != nil {
		// Log error jika perlu
		if errors.Is(err, authuc.ErrInvalidCredentials) {
			response.Error(w, http.StatusUnauthorized, "email atau password salah")
			return
		}
		// Error lainnya
		response.Error(w, http.StatusInternalServerError, "terjadi kesalahan internal")
		return
	}

	resp := dto.LoginResponse{
		AuthUser: dto.AuthUser{
			UserId:   result.UserID,
			FullName: result.FullName,
			Role:     result.Role,
		},
		TokenSet: dto.TokenSet{
			AccessToken:      result.AccessToken,
			RefreshToken:     result.RefreshToken,
			ExpiresInSeconds: result.ExpiresAt - time.Now().Unix(),
		},
	}
	response.Success(w, http.StatusOK, "login berhasil", resp)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequestDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	input := authuc.RefreshTokenInput{
		RefreshToken: req.RefreshToken,
	}

	result, err := h.refreshUC.Execute(r.Context(), input)
	if err != nil {
		// Log error jika perlu
		switch {
		case errors.Is(err, authuc.ErrInvalidRefreshToken):
			response.Error(w, http.StatusUnauthorized, "refresh token tidak valid")
		case errors.Is(err, authuc.ErrExpiredRefreshToken):
			response.Error(w, http.StatusUnauthorized, "refresh token sudah kadaluarsa")
		case errors.Is(err, authuc.ErrUserNotFound):
			response.Error(w, http.StatusUnauthorized, "akun tidak ditemukan atau tidak aktif")
		default:
			response.Error(w, http.StatusInternalServerError, "terjadi kesalahan internal")
		}
		return
	}

	// Mapping: UseCase Output -> DTO
	var session auth.Session
	resp := dto.RefreshTokenResponseDto{
		AuthUser: dto.AuthUser{
			// Ambil dari user entity jika diperlukan, atau gunakan ID dari session
			UserId: session.UserID, // Gunakan ID dari session lama
			// FullName: user.FullName, // Jika kamu ingin mengembalikan nama
		},
		TokenSet: dto.TokenSet{
			AccessToken:      result.AccessToken,
			RefreshToken:     result.RefreshToken,
			TokenType:        "Bearer",
			ExpiresInSeconds: result.ExpiresAt - time.Now().Unix(),
		},
	}

	response.Success(w, http.StatusOK, "token diperbarui", resp)
}
