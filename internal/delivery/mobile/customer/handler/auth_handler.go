package handler

import (
	"encoding/json"
	"errors"

	// "fmt"
	"net/http"
	"time"

	"modulegue/internal/delivery/mobile/customer/dto"
	auth "modulegue/internal/domain/auth"
	"modulegue/internal/middleware" // Import middleware untuk userID
	authuc "modulegue/internal/usecase/auth"

	"modulegue/pkg/response"
)

type AuthHandler struct {
	registerUC       *authuc.RegisterUseCase
	loginUC          *authuc.LoginUseCase
	refreshUC        *authuc.RefreshTokenUseCase
	logoutUC         *authuc.LogoutUseCase
	changePasswordUC *authuc.ChangePasswordUseCase
}

func NewAuthHandler(
	registerUC *authuc.RegisterUseCase,
	loginUC *authuc.LoginUseCase,
	refreshUC *authuc.RefreshTokenUseCase,
	logoutUC *authuc.LogoutUseCase, // <-- Tambahkan LogoutUseCase jika diperlukan
	changePasswordUC *authuc.ChangePasswordUseCase, // <-- Tambahkan ChangePasswordUseCase
) *AuthHandler {
	return &AuthHandler{
		registerUC:       registerUC,
		loginUC:          loginUC,
		refreshUC:        refreshUC,
		logoutUC:         logoutUC,         // Simpan LogoutUseCase di struct jika diperlukan
		changePasswordUC: changePasswordUC, // Simpan ChangePasswordUseCase di struct
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

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"` // Refresh token opsional untuk revoke spesifik
	}
	// Decode body hanya untuk mengambil refresh token jika diperlukan
	_ = json.NewDecoder(r.Body).Decode(&req) // Abaikan error jika body kosong

	// Ambil userID dari context (harus sudah disisipkan oleh middleware JWT)
	// userID, ok := r.Context().Value("user_id").(int64)
	// if !ok {
	// 	response.Error(w, http.StatusUnauthorized, "Unauthorized")
	// 	return
	// }

	// Buat input untuk usecase
	// Kita bisa mengambil access token dari header, tapi karena sudah diotentikasi,
	// kita hanya perlu userID dari context dan refresh token dari body (jika ada)
	// Untuk sederhananya, kita bisa kirim refresh token kosong untuk logout semua session
	input := authuc.LogoutInput{
		// AccessToken: r.Header.Get("Authorization")[7:], // Ambil dari header, hilangkan "Bearer "
		// Lebih baik tidak mengandalkan access token di body atau header untuk logout,
		// karena tujuan logout adalah menghapus state otentikasi.
		RefreshToken: req.RefreshToken, // Kirim refresh token jika ingin revoke spesifik
	}

	result, err := h.logoutUC.Execute(r.Context(), input) // Asumsikan logoutUC disimpan di struct AuthHandler
	if err != nil {
		// Log error jika perlu
		response.Error(w, http.StatusInternalServerError, "logout gagal")
		return
	}

	response.Success(w, http.StatusOK, result.Message, nil) // Tidak ada data yang dikembalikan
}

// ... (import lainnya)

// ChangePassword handler
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "request tidak valid")
		return
	}

	// Ambil userID dari JWT context
	// userID, ok := r.Context().Value("user_id").(int64) // Sesuaikan dengan key di middleware kamu
	// if !ok {
	// 	response.Error(w, http.StatusUnauthorized, "Unauthorized")
	// 	return
	// }

	userID, ok := middleware.UserIDFromContext(r.Context()) // <-- Gunakan fungsi helper
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Validasi input sederhana
	if req.OldPassword == "" || req.NewPassword == "" {
		response.Error(w, http.StatusBadRequest, "password lama dan password baru wajib diisi")
		return
	}

	// Buat input untuk usecase
	input := authuc.ChangePasswordInput{
		UserID:      userID, // Ambil dari context
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}

	result, err := h.changePasswordUC.Execute(r.Context(), input) // Asumsikan changePasswordUC disimpan di struct AuthHandler
	if err != nil {
		// Log error jika perlu
		switch {
		case errors.Is(err, authuc.ErrOldPasswordMismatch):
			response.Error(w, http.StatusBadRequest, "password lama tidak cocok")
		case errors.Is(err, authuc.ErrNewPasswordSameAsOld):
			response.Error(w, http.StatusBadRequest, "password baru tidak boleh sama dengan password lama")
		default:
			response.Error(w, http.StatusInternalServerError, "gagal mengubah password")
		}
		return
	}

	response.Success(w, http.StatusOK, result.Message, nil) // Tidak ada data yang dikembalikan
}

// Jangan lupa tambahkan changePasswordUC ke struct AuthHandler
// type AuthHandler struct {
// 	// ... field lainnya
// 	changePasswordUC *authuc.ChangePasswordUseCase
// }

// Dan tambahkan ke constructor NewAuthHandler
// func NewAuthHandler(..., changePasswordUC *authuc.ChangePasswordUseCase) *AuthHandler {
// 	return &AuthHandler{
// 	    // ... field lainnya
// 	    changePasswordUC: changePasswordUC,
// 	}
// }
