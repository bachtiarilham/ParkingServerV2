package repository

import (
	"context"

	model "modulegue/internal/domain/mobile/model/auth"
)

type SessionRepository interface {
	SaveSession(ctx context.Context, s model.SessionModel) error
	FindSessionByRefreshToken(ctx context.Context, token string) (model.SessionModel, error)
	RotateSession(ctx context.Context, oldRefreshToken string, newSession model.SessionModel) error
	DeleteSession(ctx context.Context, refreshToken string) error
	DeleteAllSessions(ctx context.Context, userID int64) error
}
