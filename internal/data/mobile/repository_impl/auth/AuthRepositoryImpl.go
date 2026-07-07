package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	model "modulegue/internal/domain/mobile/model/auth"
	"modulegue/internal/domain/mobile/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthRepositoryImpl struct {
	db *sql.DB
}

func NewAuthRepositoryImpl(db *sql.DB) repository.AuthRepository {
	return &AuthRepositoryImpl{db: db}
}

func (r *AuthRepositoryImpl) LoginUser(ctx context.Context, identifier, password string) (*model.UserModel, error) {
	identifier = strings.TrimSpace(identifier)
	lookupColumn := detectIdentifierColumn(identifier)

	query := `
		SELECT
			su.id,
			su.role_id,
			su.full_name,
			su.nik,
			su.phone_number,
			su.email,
			su.username,
			su.password_hash,
			su.is_verified,
			su.registered_at,
			su.created_at,
			su.updated_at
		FROM system_user su
		WHERE ` + lookupColumn + ` = ?
		LIMIT 1
	`

	var user model.UserModel
	var nik sql.NullString
	var phone sql.NullString
	var email sql.NullString
	var passwordHash sql.NullString
	var registeredAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, identifier).Scan(
		&user.UserId,
		&user.RoleId,
		&user.FullName,
		&nik,
		&phone,
		&email,
		&user.Username,
		&passwordHash,
		&user.IsVerified,
		&registeredAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("login user query: %w", err)
	}

	user.Nik = nik.String
	user.Phone = phone.String
	user.Email = email.String
	user.PasswordHash = passwordHash.String
	if registeredAt.Valid {
		user.RegisteredAt = registeredAt.Time
	}

	return &user, nil
}

func (r *AuthRepositoryImpl) LogoutUser(ctx context.Context, reqModel model.LogoutReqModel) error {
	if reqModel.RefreshToken == "" {
		return nil
	}
	return nil
}

func (r *AuthRepositoryImpl) RegisterUser(ctx context.Context, req model.RegisterRequestModel) (*model.RegisterResponseModel, error) {
	req.FullName = strings.TrimSpace(req.FullName)
	req.NIK = strings.TrimSpace(req.NIK)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	hashedPassword, err := r.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	insertQuery := `
		INSERT INTO system_user (
			role_id,
			full_name,
			nik,
			phone_number,
			email,
			username,
			password_hash,
			employment_status,
			is_verified,
			registered_at,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, 'active', ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		insertQuery,
		2,
		req.FullName,
		nullIfEmpty(req.NIK),
		nullIfEmpty(req.Phone),
		nullIfEmpty(req.Email),
		req.Username,
		hashedPassword,
		false,
		now,
		now,
		now,
	)
	if err != nil {
		return nil, fmt.Errorf("insert system_user: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("get inserted user id: %w", err)
	}

	return &model.RegisterResponseModel{
		Message: fmt.Sprintf("registrasi berhasil untuk user %d", userID),
	}, nil
}

func (r *AuthRepositoryImpl) ExistsByEmailOrUsernameOrPhone(ctx context.Context, email, username, phone string) (*model.UserExistRespModel, error) {
	query := `
		SELECT
			MAX(CASE WHEN email = ? THEN 1 ELSE 0 END) AS email_exists,
			MAX(CASE WHEN username = ? THEN 1 ELSE 0 END) AS username_exists,
			MAX(CASE WHEN phone_number = ? THEN 1 ELSE 0 END) AS phone_exists
		FROM system_user
		WHERE email = ? OR username = ? OR phone_number = ?
	`

	var emailExists, usernameExists, phoneExists int
	err := r.db.QueryRowContext(
		ctx,
		query,
		email, username, phone,
		email, username, phone,
	).Scan(&emailExists, &usernameExists, &phoneExists)
	if err != nil {
		return nil, fmt.Errorf("check user existence: %w", err)
	}

	return &model.UserExistRespModel{
		EmailExists:    emailExists == 1,
		UsernameExists: usernameExists == 1,
		PhoneExists:    phoneExists == 1,
	}, nil

}

func (r *AuthRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.UserModel, error) {
	query := `
		SELECT
			su.id,
			su.role_id,
			su.full_name,
			su.nik,
			su.phone_number,
			su.email,
			su.username,
			su.password_hash,
			su.is_verified,
			su.registered_at,
			su.created_at,
			su.updated_at
		FROM system_user su
		WHERE su.id = ?
		LIMIT 1
	`

	var user model.UserModel
	var nik sql.NullString
	var phone sql.NullString
	var email sql.NullString
	var passwordHash sql.NullString
	var registeredAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.UserId,
		&user.RoleId,
		&user.FullName,
		&nik,
		&phone,
		&email,
		&user.Username,
		&passwordHash,
		&user.IsVerified,
		&registeredAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	user.Nik = nik.String
	user.Phone = phone.String
	user.Email = email.String
	user.PasswordHash = passwordHash.String
	if registeredAt.Valid {
		user.RegisteredAt = registeredAt.Time
	}

	return &user, nil
}

func (r *AuthRepositoryImpl) ChangePasswordUser(ctx context.Context, userID int64, newPasswordHash string) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE system_user SET password_hash = ?, updated_at = NOW() WHERE id = ?`,
		newPasswordHash,
		userID,
	)
	if err != nil {
		return fmt.Errorf("update password_hash: %w", err)
	}
	return nil
}

func (r *AuthRepositoryImpl) Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(hashedBytes), nil
}

func (r *AuthRepositoryImpl) Compare(hash string, plain string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	if err != nil {
		return fmt.Errorf("password does not match: %w", err)
	}
	return nil
}

func nullIfEmpty(value string) any {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return value
}

func detectIdentifierColumn(identifier string) string {
	switch {
	case strings.Contains(identifier, "@"):
		return "su.email"
	case isPhoneIdentifier(identifier):
		return "su.phone_number"
	default:
		return "su.username"
	}
}

func isPhoneIdentifier(identifier string) bool {
	if identifier == "" {
		return false
	}

	digitCount := 0
	for i, r := range identifier {
		switch {
		case r >= '0' && r <= '9':
			digitCount++
		case r == '+' && i == 0:
		case r == ' ' || r == '-' || r == '(' || r == ')':
		default:
			return false
		}
	}

	return digitCount >= 8
}
