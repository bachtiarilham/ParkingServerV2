package auth

import (
	"context"
	"fmt"
	"modulegue/core/errorstring"
	"modulegue/core/hash"
	model "modulegue/internal/domain/mobile/model/auth"
	"modulegue/internal/domain/mobile/repository"
	"modulegue/internal/middleware"
)

type ChangePasswordUseCase struct {
	authRepo    repository.AuthRepository
	sessionRepo repository.SessionRepository
}

func NewChangePasswordUseCase(
	authRepo repository.AuthRepository,
	sessionRepo repository.SessionRepository,
) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		authRepo:    authRepo,
		sessionRepo: sessionRepo,
	}
}

func (uc *ChangePasswordUseCase) Execute(ctx context.Context, reqModel model.ChangePassReqModel) (*model.ChangePassRespModel, error) {
	// 1. Ambil user dari repository berdasarkan userID dari context
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		// Jika userID tidak ditemukan di context, berarti otentikasi gagal sebelumnya
		return nil, fmt.Errorf("user not authenticated: %w", ErrLogoutFailed)
	}
	currentUser, err := uc.authRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	// 2. Verifikasi apakah OldPassword cocok dengan hash yang disimpan
	if err := hash.Compare(currentUser.PasswordHash, reqModel.OldPassword); err != nil {
		return nil, errorstring.ErrOldPasswordMismatch
	}

	// 3. (Opsional) Cegah user mengganti password dengan password lama yang sama
	if hash.Compare(currentUser.PasswordHash, reqModel.NewPassword) == nil {
		return nil, errorstring.ErrNewPasswordSameAsOld
	}

	// 4. Hash password baru
	newHashedPassword, err := hash.Hash(reqModel.NewPassword)
	if err != nil {
		return nil, fmt.Errorf("gagal hash password baru: %w", err)
	}

	// 5. Update password user di database
	currentUser.PasswordHash = newHashedPassword
	// Jika FindByID mengembalikan pointer, kamu bisa langsung update
	// Jika tidak, kamu mungkin perlu method Update(*user.User)
	err = uc.authRepo.ChangePasswordUser(ctx, currentUser.UserId, newHashedPassword) // Asumsikan method UpdatePassword ada
	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan password baru: %w", err)
	}

	// 6. (Opsional) Hapus semua session user (force logout)
	// Ini adalah langkah keamanan: karena password berubah, session lama mungkin tidak lagi valid.
	// Tapi ini bisa merusak UX jika user sedang aktif di banyak tab/perangkat.
	// Kita bisa menambahkan opsi ini nanti jika diperlukan.
	// domain_auth.Repository.DeleteAllSessions(ctx, currentUser.ID)
	err = uc.sessionRepo.DeleteAllSessions(ctx, currentUser.UserId) // <-- Gunakan uc.sessionRepo
	if err != nil {
		// Log error, tapi jangan hentikan proses change password krn ini opsional
		// log.Printf("Gagal hapus session setelah ganti password user %d: %v", currentUser.ID, err)
	}

	return &model.ChangePassRespModel{
		Message: "password berhasil diubah",
	}, nil
}
