package usecase

import (
	"context"
	"errors"
	"fmt"
	"modulegue/internal/domain/mobile/repository"
	"modulegue/internal/middleware"
)

var (
	ErrInvalidAccessToken = errors.New("access token tidak valid")
	ErrLogoutFailed       = errors.New("logout gagal")
)

type LogoutInput struct {
	// Kita bisa menerima access token untuk extract user_id dan/atau refresh token
	// Tapi karena kita menggunakan middleware otentikasi, user_id bisa diambil dari context
	// Dan refresh token bisa diambil dari body request (jika ingin revoke refresh token)
	AccessToken  string
	RefreshToken string // Opsional: jika ingin revoke refresh token juga
}

type LogoutOutput struct {
	Message string
}

type LogoutUseCase struct {
	sessionRepo repository.SessionRepository
}

func NewLogoutUseCase(sessionRepo repository.SessionRepository) *LogoutUseCase {
	return &LogoutUseCase{
		sessionRepo: sessionRepo,
	}
}

func (uc *LogoutUseCase) Execute(ctx context.Context, input LogoutInput) (LogoutOutput, error) {
	// 1. (Opsional) Validasi access token untuk memastikan user benar-benar login saat ini
	// Ini bisa dilakukan, tapi karena kita sudah menggunakan middleware otentikasi,
	// maka user_id seharusnya sudah valid di context saat request ini diterima.
	// Jadi langkah ini bisa dilewati atau digunakan untuk verifikasi ekstra.
	// _, err := jwt.ParseClaimsHS256(input.AccessToken, secret) // Secret perlu diinject ke usecase ini
	// if err != nil {
	//     return LogoutOutput{}, ErrInvalidAccessToken
	// }

	// 2. Ambil userID dari context (ini adalah user yang sedang logout)
	// userID, ok := ctx.Value("user_id").(int64) // Sesuaikan dengan key di middleware kamu
	userID, ok := middleware.UserIDFromContext(ctx) // Gunakan helper dari middleware untuk mengambil userID
	if !ok {
		// Jika userID tidak ditemukan di context, berarti otentikasi gagal sebelumnya
		return LogoutOutput{}, fmt.Errorf("user not authenticated: %w", ErrLogoutFailed)
	}

	// 3. Hapus semua session yang terkait dengan user ini (atau hanya session tertentu jika refresh token disediakan)
	var err error
	if input.RefreshToken != "" {
		// Jika refresh token disediakan, hapus session spesifik
		err = uc.sessionRepo.DeleteSession(ctx, input.RefreshToken)
		if err != nil {
			// Log error jika perlu
			return LogoutOutput{}, fmt.Errorf("failed to delete specific session: %w", err)
		}
	} else {
		// Jika tidak ada refresh token, hapus semua session untuk user ini (logout dari semua perangkat)
		err = uc.sessionRepo.DeleteAllSessions(ctx, userID)
		if err != nil {
			// Log error jika perlu
			return LogoutOutput{}, fmt.Errorf("failed to delete all user sessions: %w", err)
		}
	}

	return LogoutOutput{
		Message: "logout berhasil",
	}, nil
}
