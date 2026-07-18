package auth

import (
	"context"
	"errors"
	"fmt"
	model "modulegue/internal/domain/shared/model/auth"
	"modulegue/internal/domain/shared/repository"
)

var (
	ErrInvalidAccessToken = errors.New("access token tidak valid")
	ErrLogoutFailed       = errors.New("logout gagal")
)

type LogoutUseCase struct {
	authRepo    repository.AuthRepository
	sessionRepo repository.SessionRepository
}

func NewLogoutUseCase(
	authRepo repository.AuthRepository,
	sessionRepo repository.SessionRepository,
) *LogoutUseCase {
	return &LogoutUseCase{
		authRepo:    authRepo,
		sessionRepo: sessionRepo,
	}
}

func (uc *LogoutUseCase) Execute(ctx context.Context, reqModel model.LogoutReqModel) (*model.LogoutRespModel, error) {

	err := uc.authRepo.LogoutUser(ctx, reqModel)
	if err != nil {
		return nil, ErrLogoutFailed
	}

	// 3. Hapus semua session yang terkait dengan user ini (atau hanya session tertentu jika refresh token disediakan)
	if reqModel.RefreshToken != "" {
		// Jika refresh token disediakan, hapus session spesifik
		err = uc.sessionRepo.DeleteSession(ctx, reqModel.RefreshToken)
		if err != nil {
			// Log error jika perlu
			return nil, fmt.Errorf("failed to delete specific session: %w", err)
		}
	}
	// else {
	// 	// Jika tidak ada refresh token, hapus semua session untuk user ini (logout dari semua perangkat)
	// 	err = uc.authRepo.DeleteAllSessions(ctx, userID)
	// 	if err != nil {
	// 		// Log error jika perlu
	// 		return nil, fmt.Errorf("failed to delete all user sessions: %w", err)
	// 	}
	// }

	return &model.LogoutRespModel{
		Message: "logout berhasil",
	}, nil
}
