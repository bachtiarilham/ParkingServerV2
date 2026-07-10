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
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(
		ctx,
		`
		UPDATE user_session
		SET
			revoked_at = NOW(),
			updated_at = NOW()
		WHERE user_id = ?
		AND device_id = ?
		AND revoked_at IS NULL;
		`,
		s.UserID,
		s.DeviceId,
	); err != nil {
		return fmt.Errorf("revoke previous session: %w", err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`
		INSERT INTO user_session (
			user_id,
			refresh_token,
			device_id,
			device_name,
			fcm_token,
			ip_address,
			user_agent,
			expired_at,
			revoked_at,
			created_at,
			updated_at
		)
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			NULL,
			NOW(),
			NOW()
		);
		`,
		s.UserID,
		s.RefreshToken,
		s.DeviceId,
		s.DeviceName,
		s.FcmToken,
		s.IpAddress,
		s.UserAgent,
		s.ExpiresAt,
	); err != nil {
		return fmt.Errorf("insert session: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (r *SessionRepositoryImpl) SaveLastLogin(ctx context.Context, userId int64) error {
	_, err := r.db.ExecContext(
		ctx,
		`
		UPDATE user_auth
		SET
			last_login_at = NOW(),
			failed_login_count = 0,
			locked_until = NULL,
			updated_at = NOW()
		WHERE user_id = ?;
		`,
		userId,
	)
	if err != nil {
		return fmt.Errorf("save last login: %w", err)
	}
	return nil
}

func (r *SessionRepositoryImpl) IsSessionActive(ctx context.Context, reqModel model.RefreshTokenReqModel) (model.SessionModel, error) {
	var session model.SessionModel
	var revokedAt sql.NullTime

	err := r.db.QueryRowContext(
		ctx,
		`
		SELECT
			user_id,
			refresh_token,
			expired_at,
			revoked_at
		FROM user_session
		WHERE refresh_token = ?
		AND revoked_at IS NULL
		AND expired_at > NOW()
		LIMIT 1;
		`,
		reqModel.RefreshToken,
	).Scan(
		&session.UserID,
		&session.RefreshToken,
		&session.ExpiresAt,
		&revokedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.SessionModel{}, fmt.Errorf("session not found")
		}
		return model.SessionModel{}, fmt.Errorf("find session by refresh token: %w", err)
	}

	if revokedAt.Valid {
		session.RevokedAt = revokedAt.Time
	}

	return session, nil
}

func (r *SessionRepositoryImpl) RotateSession(ctx context.Context, oldRefreshToken string, newSession model.SessionModel) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(
		ctx,
		`
		UPDATE user_session
		SET
			revoked_at = NOW(),
			updated_at = NOW()
		WHERE refresh_token = ?
		AND revoked_at IS NULL;
		`,
		oldRefreshToken,
	); err != nil {
		return fmt.Errorf("delete old session: %w", err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`
		INSERT INTO user_session (
		user_id,
		refresh_token,
		device_id,
		device_name,
		fcm_token,
		ip_address,
		user_agent,
		expired_at,
		revoked_at,
		created_at,
		updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NULL, NOW(), NOW());
		`,
		newSession.UserID,
		newSession.RefreshToken,
		newSession.DeviceId,
		newSession.DeviceName,
		newSession.FcmToken,
		newSession.IpAddress,
		newSession.UserAgent,
		newSession.ExpiresAt,
	); err != nil {
		return fmt.Errorf("insert new session: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (r *SessionRepositoryImpl) DeleteSession(ctx context.Context, refreshToken string) error {
	_, err := r.db.ExecContext(ctx, `
	UPDATE user_session
	SET
		revoked_at = NOW(),
		updated_at = NOW()
	WHERE refresh_token = ?
	AND revoked_at IS NULL;
	`, refreshToken)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func (r *SessionRepositoryImpl) DeleteAllSessions(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx,
		`
	UPDATE user_session
	SET
		revoked_at = NOW(),
		updated_at = NOW()
	WHERE user_id = ?
	AND revoked_at IS NULL;
	`, userID)
	if err != nil {
		return fmt.Errorf("delete all sessions: %w", err)
	}
	return nil
}
