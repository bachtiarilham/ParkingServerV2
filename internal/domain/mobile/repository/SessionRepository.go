package repository

import (
	"context"

	model "modulegue/internal/domain/mobile/model/auth"
)

type SessionRepository interface {
	SaveSession(ctx context.Context, s model.SessionModel) error
	SaveLastLogin(ctx context.Context, userId int64) error
	IsSessionActive(ctx context.Context, s model.RefreshTokenReqModel) (model.SessionModel, error)
	RotateSession(ctx context.Context, oldRefreshToken string, newSession model.SessionModel) error
	DeleteSession(ctx context.Context, refreshToken string) error
	DeleteAllSessions(ctx context.Context, userID int64) error
}
