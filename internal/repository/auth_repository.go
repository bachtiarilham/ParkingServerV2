package repository

import (
	"context"
	// "crypto"
	// "crypto/sha256"
	"database/sql"
	"fmt"

	"modulegue/internal/domain/auth"
	// "golang.org/x/crypto/bcrypt"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) auth.Repository {
	return &AuthRepository{db: db}
}

// func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (*auth.Credential, error) {
// 	// Gunakan tabel yang benar: system_user
// 	query := `
//         SELECT id, email, password_hash, created_at, updated_at
//         FROM system_user
//         WHERE email = ?
//         LIMIT 1
//     `
// 	var credential auth.Credential // Hanya gunakan credential
// 	err := r.db.QueryRowContext(ctx, query, email).Scan(
// 		&credential.ID,           // Sesuaikan field Credential.ID dengan kolom 'id'
// 		&credential.Email,        // Kolom 'email'
// 		&credential.PasswordHash, // Kolom 'password_hash'
// 		&credential.CreatedAt,    // Kolom 'created_at'
// 		&credential.UpdatedAt,    // Kolom 'updated_at'
// 	)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, fmt.Errorf("user not found for email: %s", email)
// 		}
// 		return nil, fmt.Errorf("query user by email: %w", err)
// 	}
// 	return &credential, nil
// }

// func (r *AuthRepository) Hash(password string) (string, error) {
// 	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return "", fmt.Errorf("hash password: %w", err)
// 	}
// 	return string(hashedBytes), nil
// }

// func (r *AuthRepository) Compare(hash string, plain string) error {
// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
// 	if err != nil {
// 		return fmt.Errorf("password does not match: %w", err)
// 	}
// 	return nil
// }

func (r *AuthRepository) SaveSession(ctx context.Context, s auth.Session) error {
	query := `
		INSERT INTO system_user_sessions (user_id, refresh_token, expires_at, created_at)
		VALUES (?, ?, ?, NOW())
	`
	_, err := r.db.ExecContext(ctx, query, s.UserID, s.RefreshToken, s.ExpiresAt)
	if err != nil {
		return fmt.Errorf("save session: %w", err)
	}
	return nil
}

func (r *AuthRepository) FindSessionByRefreshToken(ctx context.Context, token string) (auth.Session, error) {
	query := `
		SELECT id, user_id, refresh_token, expires_at, created_at
		FROM system_user_sessions
		WHERE refresh_token = ?
		LIMIT 1
	`
	var s auth.Session
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&s.ID,
		&s.UserID,
		&s.RefreshToken,
		&s.ExpiresAt,
		&s.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return auth.Session{}, fmt.Errorf("session not found for refresh token")
		}
		return auth.Session{}, fmt.Errorf("find session: %w", err)
	}
	return s, nil
}

func (r *AuthRepository) DeleteSession(ctx context.Context, refreshToken string) error {
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM system_user_sessions WHERE refresh_token = ?`,
		refreshToken,
	)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no session found with refresh token to delete")
	}

	return nil
}

func (r *AuthRepository) DeleteAllSessions(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM system_user_sessions WHERE user_id = ?`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("delete all sessions: %w", err)
	}
	return nil
}
