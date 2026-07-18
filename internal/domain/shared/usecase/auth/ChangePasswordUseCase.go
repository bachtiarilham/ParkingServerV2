package auth

import (
	"context"
	"fmt"
	"modulegue/core/errorstring"
	"modulegue/core/hash"
	model "modulegue/internal/domain/shared/model/auth"
	"modulegue/internal/domain/shared/repository"
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
	currentUser, err := uc.authRepo.FindByID(ctx, reqModel.UserId)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	if reqModel.OldPassword == reqModel.NewPassword {
		return nil, errorstring.ErrNewPasswordSameAsOld
	}

	if err := hash.Compare(currentUser.PasswordHash, reqModel.OldPassword); err != nil {
		return nil, errorstring.ErrOldPasswordMismatch
	}

	// 4. Hash password baru
	newHashedPassword, err := hash.Hash(reqModel.NewPassword)
	if err != nil {
		return nil, fmt.Errorf("gagal hash password baru: %w", err)
	}

	// 5. Update password user di database
	err = uc.authRepo.ChangePasswordUser(ctx, reqModel.UserId, newHashedPassword) // Asumsikan method UpdatePassword ada
	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan password baru: %w", err)
	}

	// 6. force logout
	err = uc.sessionRepo.DeleteAllSessions(ctx, reqModel.UserId) // <-- Gunakan uc.sessionRepo
	if err != nil {
		return nil, fmt.Errorf("Gagal hapus session setelah ganti password user: %w", err)
	}

	return &model.ChangePassRespModel{
		Message: "password berhasil diubah",
	}, nil
}
