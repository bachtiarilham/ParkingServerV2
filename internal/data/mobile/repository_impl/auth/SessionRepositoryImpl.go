package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	model "modulegue/internal/domain/mobile/model/auth"
	"modulegue/internal/domain/mobile/repository"
)

type SessionRepositoryImpl struct {
	db *sql.DB
}

func NewSessionRepositoryImpl(db *sql.DB) repository.SessionRepository {
	return &SessionRepositoryImpl{db: db}
}

func (r *SessionRepositoryImpl) SaveSession(ctx context.Context, s model.SessionModel) error {
	_, err := r.db.ExecContext(
		ctx,
		`
		INSERT INTO system_user_sessions (
			user_id,
			refresh_token,
			token_type,
			expires_at,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, NOW(), NOW())
		`,
		s.UserID,
		s.RefreshToken,
		s.TokenType,
		s.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("save session: %w", err)
	}
	return nil
}

func (r *SessionRepositoryImpl) FindSessionByRefreshToken(ctx context.Context, token string) (model.SessionModel, error) {
	var session model.SessionModel
	var sessionID int64

	err := r.db.QueryRowContext(
		ctx,
		`
		SELECT id, user_id, refresh_token, token_type, expires_at, created_at, updated_at
		FROM system_user_sessions
		WHERE refresh_token = ?
		LIMIT 1
		`,
		token,
	).Scan(
		&sessionID,
		&session.UserID,
		&session.RefreshToken,
		&session.TokenType,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.SessionModel{}, fmt.Errorf("session not found")
		}
		return model.SessionModel{}, fmt.Errorf("find session by refresh token: %w", err)
	}

	return session, nil
}

func (r *SessionRepositoryImpl) RotateSession(ctx context.Context, oldRefreshToken string, newSession model.SessionModel) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM system_user_sessions WHERE refresh_token = ?`, oldRefreshToken); err != nil {
		return fmt.Errorf("delete old session: %w", err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`
		INSERT INTO system_user_sessions (
			user_id,
			refresh_token,
			token_type,
			expires_at,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?)
		`,
		newSession.UserID,
		newSession.RefreshToken,
		newSession.TokenType,
		newSession.ExpiresAt,
		newSession.CreatedAt,
		newSession.UpdatedAt,
	); err != nil {
		return fmt.Errorf("insert new session: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (r *SessionRepositoryImpl) DeleteSession(ctx context.Context, refreshToken string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM system_user_sessions WHERE refresh_token = ?`, refreshToken)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func (r *SessionRepositoryImpl) DeleteAllSessions(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM system_user_sessions WHERE user_id = ?`, userID)
	if err != nil {
		return fmt.Errorf("delete all sessions: %w", err)
	}
	return nil
}
