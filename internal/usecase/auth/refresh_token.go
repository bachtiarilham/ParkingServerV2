package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	domain_auth "modulegue/internal/domain/auth"
	"modulegue/internal/domain/user"
	"modulegue/pkg/jwt" // Pastikan kamu punya fungsi untuk generate token
	"time"
	// "modulegue/pkg/hash" // Jika refresh token di-hash sebelum disimpan
)

var (
	ErrInvalidRefreshToken = errors.New("refresh token tidak valid")
	ErrExpiredRefreshToken = errors.New("refresh token sudah kadaluarsa")
	// ErrUserNotFound        = errors.New("user tidak ditemukan")
)

type RefreshTokenInput struct {
	RefreshToken string
}

type RefreshTokenOutput struct {
	AccessToken  string
	RefreshToken string // Token refresh baru
	ExpiresAt    int64  // Timestamp unix untuk access token
}

type RefreshTokenUseCase struct {
	authRepo   domain_auth.Repository
	userRepo   user.Repository // Butuh untuk verifikasi user aktif
	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewRefreshTokenUseCase(
	authRepo domain_auth.Repository,
	userRepo user.Repository, // Tambahkan repo user
	jwtSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		authRepo:   authRepo,
		userRepo:   userRepo, // Simpan repo user
		jwtSecret:  jwtSecret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (uc *RefreshTokenUseCase) Execute(ctx context.Context, input RefreshTokenInput) (RefreshTokenOutput, error) {
	// 1. Cari session berdasarkan refresh token
	session, err := uc.authRepo.FindSessionByRefreshToken(ctx, input.RefreshToken)
	if err != nil {
		return RefreshTokenOutput{}, ErrInvalidRefreshToken
	}

	// 2. Cek apakah refresh token sudah kadaluarsa
	if time.Now().After(session.ExpiresAt) {
		// Hapus session lama jika kadaluarsa
		uc.authRepo.DeleteSession(ctx, input.RefreshToken) // Log error jika perlu, tapi jangan hentikan proses
		return RefreshTokenOutput{}, ErrExpiredRefreshToken
	}

	// 3. Cek apakah user masih aktif (opsional, tapi bagus untuk security)
	// user, err := uc.userRepo.GetByID(ctx, session.UserID) // Tambahkan method GetByID ke user.Repository
	// if err != nil || user.EmploymentStatus != "active" {
	// 	// Jika user tidak ditemukan atau tidak aktif, hapus session dan tolak
	// 	uc.authRepo.DeleteSession(ctx, input.RefreshToken)
	// 	return RefreshTokenOutput{}, ErrUserNotFound
	// }

	// 4. Hapus session lama (refresh token lama tidak bisa digunakan lagi)
	err = uc.authRepo.DeleteSession(ctx, input.RefreshToken)
	if err != nil {
		// Log error, tapi lanjutkan
		// fmt.Printf("Warning: could not delete old session: %v\n", err)
	}

	// 5. Generate token baru
	now := time.Now()
	accessExp := now.Add(uc.accessTTL)
	newRefreshExp := now.Add(uc.refreshTTL)

	// a. Buat access token baru
	newAccessToken, err := jwt.SignHS256(jwt.Claims{
		Subject:    fmt.Sprintf("%d", session.UserID),
		UserID:     session.UserID,
		Expiration: accessExp.Unix(),
		IssuedAt:   now.Unix(),
		Type:       "access",
	}, uc.jwtSecret)
	if err != nil {
		return RefreshTokenOutput{}, fmt.Errorf("generate new access token: %w", err)
	}

	// b. Buat refresh token baru (bisa gunakan string acak, atau hash dari string acak)
	newRefreshToken, err := generateSecureRandomString(32) // Implementasi generateSecureRandomString
	if err != nil {
		return RefreshTokenOutput{}, fmt.Errorf("generate new refresh token: %w", err)
	}

	// c. Simpan session baru ke database
	newSession := domain_auth.Session{
		UserID:       session.UserID,
		RefreshToken: newRefreshToken, // Gunakan refresh token baru
		ExpiresAt:    newRefreshExp,
	}
	err = uc.authRepo.SaveSession(ctx, newSession)
	if err != nil {
		return RefreshTokenOutput{}, fmt.Errorf("save new session: %w", err)
	}

	return RefreshTokenOutput{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    accessExp.Unix(),
	}, nil
}

// Helper function untuk generate refresh token baru (contoh sederhana, gunakan crypto/rand)
func generateSecureRandomString(length int) (string, error) {
	// Gunakan crypto/rand untuk keamanan
	// Contoh: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	// Kita gunakan contoh sederhana di sini, ganti dengan implementasi yang benar-benar aman
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
