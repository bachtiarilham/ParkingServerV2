package auth

import (
	"context"
	"errors"
	"fmt"
	domain_auth "modulegue/internal/domain/auth" // Alias untuk menghindari konflik
	"modulegue/internal/domain/user"
	"modulegue/pkg/hash" // Gunakan pkg/hash kamu
)

var (
	ErrOldPasswordMismatch  = errors.New("password lama tidak cocok")
	ErrNewPasswordSameAsOld = errors.New("password baru tidak boleh sama dengan password lama")
	ErrInvalidUserID        = errors.New("user tidak valid")
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
	userRepo user.Repository
	authRepo domain_auth.Repository // Tidak digunakan secara langsung untuk password user di sini, karena password ada di user_repo
}

func NewChangePasswordUseCase(userRepo user.Repository, authRepo domain_auth.Repository) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

func (uc *ChangePasswordUseCase) Execute(ctx context.Context, input ChangePasswordInput) (ChangePasswordOutput, error) {
	// 1. Ambil user dari repository berdasarkan userID dari context
	currentUser, err := uc.userRepo.FindByID(ctx, input.UserID) // Asumsikan FindByID ada di userRepo
	if err != nil {
		// Log error jika perlu
		return ChangePasswordOutput{}, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	// 2. Verifikasi apakah OldPassword cocok dengan hash yang disimpan
	if err := hash.Compare(currentUser.PasswordHash, input.OldPassword); err != nil {
		// Password lama salah
		return ChangePasswordOutput{}, ErrOldPasswordMismatch
	}

	// 3. (Opsional) Cegah user mengganti password dengan password lama yang sama
	if hash.Compare(currentUser.PasswordHash, input.NewPassword) == nil {
		return ChangePasswordOutput{}, ErrNewPasswordSameAsOld
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
	err = uc.authRepo.DeleteAllSessions(ctx, currentUser.ID) // <-- Gunakan uc.authRepo
	if err != nil {
		// Log error, tapi jangan hentikan proses change password krn ini opsional
		// log.Printf("Gagal hapus session setelah ganti password user %d: %v", currentUser.ID, err)
	}

	return ChangePasswordOutput{
		Message: "password berhasil diubah",
	}, nil
}
