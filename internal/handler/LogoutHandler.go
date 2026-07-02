package handler

import (
	"encoding/json"
	"net/http"

	"modulegue/core/response"
	"modulegue/internal/domain/mobile/usecase"
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
	input := usecase.LogoutInput{
		// AccessToken: r.Header.Get("Authorization")[7:], // Ambil dari header, hilangkan "Bearer "
		// Lebih baik tidak mengandalkan access token di body atau header untuk logout,
		// karena tujuan logout adalah menghapus state otentikasi.
		RefreshToken: req.RefreshToken, // Kirim refresh token jika ingin revoke spesifik
	}

	result, err := h.logoutUc.Execute(r.Context(), input) // Asumsikan logoutUC disimpan di struct AuthHandler
	if err != nil {
		// Log error jika perlu
		response.Error(w, http.StatusInternalServerError, "logout gagal")
		return
	}

	response.Success(w, http.StatusOK, result.Message, nil) // Tidak ada data yang dikembalikan
}
