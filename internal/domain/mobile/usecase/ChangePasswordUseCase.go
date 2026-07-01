package usecase

import (
	"context"
	"fmt"
	"modulegue/core/errorstring"
	"modulegue/core/hash"
	"modulegue/internal/domain/mobile/repository"
)

type ChangePasswordInput struct {
	UserID      int64 // Didapat dari context JWT
	OldPassword string
	NewPassword string
}

type ChangePasswordOutput struct {
	Message string
}

type ChangePasswordUseCase struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository // Tidak digunakan secara langsung untuk password user di sini, karena password ada di user_repo
}

func NewChangePasswordUseCase(userRepo repository.UserRepository, sessionRepo repository.SessionRepository) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (uc *ChangePasswordUseCase) Execute(ctx context.Context, input ChangePasswordInput) (ChangePasswordOutput, error) {
	// 1. Ambil user dari repository berdasarkan userID dari context
	currentUser, err := uc.userRepo.GetUser(ctx, input.UserID) // Asumsikan FindByID ada di userRepo
	if err != nil {
		// Log error jika perlu
		return ChangePasswordOutput{}, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	// 2. Verifikasi apakah OldPassword cocok dengan hash yang disimpan
	if err := hash.Compare(currentUser.PasswordHash, input.OldPassword); err != nil {
		// Password lama salah
		return ChangePasswordOutput{}, errorstring.ErrOldPasswordMismatch
	}

	// 3. (Opsional) Cegah user mengganti password dengan password lama yang sama
	if hash.Compare(currentUser.PasswordHash, input.NewPassword) == nil {
		return ChangePasswordOutput{}, errorstring.ErrNewPasswordSameAsOld
	}

	// 4. Hash password baru
	newHashedPassword, err := hash.Hash(input.NewPassword)
	if err != nil {
		return ChangePasswordOutput{}, fmt.Errorf("gagal hash password baru: %w", err)
	}

	// 5. Update password user di database
	currentUser.PasswordHash = newHashedPassword
	// Jika FindByID mengembalikan pointer, kamu bisa langsung update
	// Jika tidak, kamu mungkin perlu method Update(*user.User)
	err = uc.userRepo.UpdatePassword(ctx, currentUser.ID, newHashedPassword) // Asumsikan method UpdatePassword ada
	if err != nil {
		return ChangePasswordOutput{}, fmt.Errorf("gagal menyimpan password baru: %w", err)
	}

	// 6. (Opsional) Hapus semua session user (force logout)
	// Ini adalah langkah keamanan: karena password berubah, session lama mungkin tidak lagi valid.
	// Tapi ini bisa merusak UX jika user sedang aktif di banyak tab/perangkat.
	// Kita bisa menambahkan opsi ini nanti jika diperlukan.
	// domain_auth.Repository.DeleteAllSessions(ctx, currentUser.ID)
	err = uc.sessionRepo.DeleteAllSessions(ctx, currentUser.ID) // <-- Gunakan uc.sessionRepo
	if err != nil {
		// Log error, tapi jangan hentikan proses change password krn ini opsional
		// log.Printf("Gagal hapus session setelah ganti password user %d: %v", currentUser.ID, err)
	}

	return ChangePasswordOutput{
		Message: "password berhasil diubah",
	}, nil
}
