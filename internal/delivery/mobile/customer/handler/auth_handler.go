package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

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

// func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
// 	var req dto.RegisterRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		response.Error(w, http.StatusBadRequest, "request tidak valid")
// 		return
// 	}

// 	err := h.registerUC.Execute(r.Context(), dto.RegisterRequest{
// 		FullName: req.FullName,
// 		Nik:      req.Nik,
// 		Email:    req.Email,
// 		Phone:    req.Phone,
// 		Username: req.Username,
// 		Password: req.Password,
// 	})
// 	if err != nil {
// 		response.Error(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	response.Success(w, http.StatusCreated, "registrasi berhasil", nil)
// }

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
		UserID:  fmt.Sprintf("%d", result.UserID),
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

	// response.Success(w, http.StatusOK, "login berhasil", dto.LoginResponse{
	// 	UserID:       result.UserID,
	// 	FullName:     result.FullName,
	// 	Role:         result.Role,
	// 	AccessToken:  result.AccessToken,
	// 	RefreshToken: result.RefreshToken,
	// 	ExpiresAt:    result.ExpiresAt,
	// })

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
	response.Error(w, http.StatusInternalServerError, "lu gagal login")

}
