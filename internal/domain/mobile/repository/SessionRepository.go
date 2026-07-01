package repository

import (
	"context"
	"modulegue/internal/domain/mobile/model"
)

type SessionRepository interface {
	// FindByEmail(
	// 	ctx context.Context,
	// 	email string,
	// ) (*Credential, error)
	// Hash(password string) (string, error)
	// Compare(hash string, plain string) error
	// SaveSession menyimpan refresh token session baru
	SaveSession(ctx context.Context, s model.SessionModel) error

	// FindSessionByRefreshToken mencari session aktif berdasarkan token
	FindSessionByRefreshToken(ctx context.Context, token string) (model.SessionModel, error)

	// DeleteSession menghapus session (logout)
	DeleteSession(ctx context.Context, refreshToken string) error

	// DeleteAllSessions menghapus semua session user (logout semua device)
	DeleteAllSessions(ctx context.Context, userID int64) error
}
