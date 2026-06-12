package repository

import (
	"context"
	"database/sql"
	"fmt"

	"modulegue/internal/domain/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.Repository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	query := `
		SELECT id, role_id, full_name, email, username, password_hash, employment_status, is_verified, registered_at, created_at, updated_at
		FROM system_user
		WHERE username = ?
		LIMIT 1
	`
	var u user.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&u.ID,
		&u.RoleID,
		&u.FullName,
		&u.Email,
		&u.Username,
		&u.PasswordHash, // Ini yang akan digunakan untuk verifikasi
		&u.EmploymentStatus,
		&u.IsVerified,
		&u.RegisteredAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found for username: %s", username)
		}
		return nil, fmt.Errorf("find user by username: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, role_id, full_name, email, username, password_hash, employment_status, is_verified, registered_at, created_at, updated_at
		FROM system_user
		WHERE email = ?
		LIMIT 1
	`
	var u user.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.RoleID,
		&u.FullName,
		&u.Email,
		&u.Username,
		&u.PasswordHash,
		&u.EmploymentStatus,
		&u.IsVerified,
		&u.RegisteredAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found for email: %s", email)
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	query := `
		INSERT INTO system_user (role_id, full_name, nik, phone_number, email, username, password_hash, employment_status, is_verified, registered_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	result, err := r.db.ExecContext(ctx, query,
		u.RoleID,
		u.FullName,
		u.Nik,
		u.Phone,
		u.Email,
		u.Username,
		u.PasswordHash,
		u.EmploymentStatus,
		u.IsVerified,
		u.RegisteredAt,
	)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	// Ambil ID yang di-generate oleh database
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}
	u.ID = id // Update ID di struct user

	// TODO: Jika ingin membuat wallet default, panggil walletRepo.CreateDefaultWallet(u.ID)
	// Contoh pseudocode:
	// walletRepo.CreateDefaultWallet(ctx, id)

	return nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	query := `
		UPDATE system_user
		SET role_id = ?, full_name = ?, email = ?, username = ?, password_hash = ?, employment_status = ?, is_verified = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, u.RoleID, u.FullName, u.Email, u.Username, u.PasswordHash, u.EmploymentStatus, u.IsVerified, u.ID)
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*user.User, error) {
	query := `
		SELECT id, role_id, full_name, nik, phone_number, email, username, password_hash, employment_status, is_verified, registered_at, created_at, updated_at
		FROM system_user
		WHERE id = ?
		LIMIT 1
	`
	var u user.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.RoleID,
		&u.FullName,
		&u.Nik,
		&u.Phone,
		&u.Email,
		&u.Username,
		&u.PasswordHash, // Ini seharusnya tidak dikembalikan ke client
		&u.EmploymentStatus,
		&u.IsVerified,
		&u.RegisteredAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found for id: %d", id)
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID int64, newPasswordHash string) error {
	query := `UPDATE system_user SET password_hash = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, newPasswordHash, userID)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}
