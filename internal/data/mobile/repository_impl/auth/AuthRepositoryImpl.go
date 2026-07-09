package auth

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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

func (r *AuthRepositoryImpl) LoginUser(ctx context.Context, reqModel model.LoginRequestModel) (*model.TokenSetModel, *model.LoginRespModel, error) {
	identifier := strings.TrimSpace(reqModel.Identity)
	password := strings.TrimSpace(reqModel.Password)
	if identifier == "" || password == "" {
		return nil, nil, fmt.Errorf("identity and password are required")
	}

	query := `
	SELECT
		ui.id AS userId,
		ui.role_id AS roleId,
		ua.password_hash AS passwordHash
	FROM user_identity ui
	JOIN user_auth ua
		ON ua.user_id = ui.id
	JOIN master_role mr
		ON mr.id = ui.role_id
	WHERE
		(
			ui.username = ?
			OR ui.email = ?
			OR ui.phone_number = ?
		)
		AND ui.status = 'ACTIVE'
		AND mr.is_active = 1
		AND (
			ua.locked_until IS NULL
			OR ua.locked_until <= NOW()
		)
	LIMIT 1;
	`
	var (
		userID       int64
		roleID       int64
		passwordHash string
	)

	if err := r.db.QueryRowContext(
		ctx,
		query,
		identifier,
		identifier,
		identifier,
	).Scan(&userID, &roleID, &passwordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, fmt.Errorf("invalid credentials")
		}
		return nil, nil, fmt.Errorf("login user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	return &model.TokenSetModel{}, &model.LoginRespModel{
		UserId:       userID,
		RoleId:       roleID,
		PasswordHash: passwordHash,
	}, nil
}

func (r *AuthRepositoryImpl) LogoutUser(ctx context.Context, reqModel model.LogoutReqModel) error {
	if reqModel.RefreshToken == "" {
		return nil
	}
	return nil
}

func (r *AuthRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.LoginRespModel, error) {
	query := `
	SELECT
		ua.password_hash AS passwordHash
	FROM user_auth ua
	JOIN user_identity ui
		ON ui.id = ua.user_id
	WHERE ua.user_id = ?
	AND ui.status = 'ACTIVE'
	LIMIT 1;
	`
	var (
		passwordHash string
	)

	if err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(&passwordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("login user: %w", err)
	}

	return &model.LoginRespModel{
		PasswordHash: passwordHash,
	}, nil
}

func (r *AuthRepositoryImpl) RegisterUser(ctx context.Context, req model.RegisterRequestModel) error {
	req.FullName = strings.TrimSpace(req.FullName)
	req.NIK = strings.TrimSpace(req.NIK)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	hashedPassword, err := r.Hash(req.Password)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queryInsertUser := `
	INSERT INTO user_identity (
		full_name,
		nik,
		phone_number,
		email,
		username,
		role_id,
		is_verified,
		status,
		created_at,
		updated_at
	)
	VALUES (
		?, ?, ?, ?, ?,
		(SELECT id FROM master_role WHERE role_code = 'CUSTOMER' LIMIT 1),
		0,
		'ACTIVE',
		NOW(),
		NOW()
	);
	`

	result, err := tx.ExecContext(
		ctx,
		queryInsertUser,
		req.FullName,
		req.NIK,
		req.Phone,
		req.Email,
		req.Username,
	)
	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	queryInsertAuth := `
	INSERT INTO user_auth (
		user_id,
		password_hash,
		failed_login_count,
		created_at,
		updated_at
	)
	VALUES (?, ?, 0, NOW(), NOW());
	`

	_, err = tx.ExecContext(ctx, queryInsertAuth, userId, hashedPassword)
	if err != nil {
		return err
	}

	queryInsertProfile := `
	INSERT INTO user_profile (
		user_id,
		language,
		timezone,
		created_at,
		updated_at
	)
	VALUES (?, 'id', 'Asia/Jakarta', NOW(), NOW());
	`

	_, err = tx.ExecContext(ctx, queryInsertProfile, userId)
	if err != nil {
		return err
	}

	queryInsertWallet := `
	INSERT INTO wallet_account (
		user_id,
		wallet_type_id,
		current_balance_amount,
		is_active,
		created_at,
		updated_at
	)
	SELECT
		?,
		id,
		0,
		1,
		NOW(),
		NOW()
	FROM master_wallet_type
	WHERE wallet_type_code = 'BALANCE';
	`

	_, err = tx.ExecContext(ctx, queryInsertWallet, userId)
	if err != nil {
		return err
	}

	queryInsertSummary := `
	INSERT INTO summary_user_home (
		user_id,
		role_id,
		full_name,
		username,
		email,
		phone_number,
		photo_url,
		saldo,
		updated_at
	)
	SELECT
		ui.id,
		ui.role_id,
		ui.full_name,
		ui.username,
		ui.email,
		ui.phone_number,
		ui.photo_url,
		0,
		NOW()
	FROM user_identity ui
	WHERE ui.id = ?;
	`

	_, err = tx.ExecContext(ctx, queryInsertSummary, userId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	_, _ = result.RowsAffected()
	return nil

}

func (r *AuthRepositoryImpl) ExistsByIdentity(ctx context.Context, nik, email, username, phone string) (*model.UserExistRespModel, error) {
	query := `
	SELECT
    CASE WHEN EXISTS (
        SELECT 1 FROM user_identity WHERE nik = ? LIMIT 1
    ) THEN 1 ELSE 0 END AS nikUsed,

    CASE WHEN EXISTS (
        SELECT 1 FROM user_identity WHERE phone_number = ? LIMIT 1
    ) THEN 1 ELSE 0 END AS phoneUsed,

    CASE WHEN EXISTS (
        SELECT 1 FROM user_identity WHERE email = ? LIMIT 1
    ) THEN 1 ELSE 0 END AS emailUsed,

    CASE WHEN EXISTS (
        SELECT 1 FROM user_identity WHERE username = ? LIMIT 1
    ) THEN 1 ELSE 0 END AS usernameUsed;
	`
	var nikExists, emailExists, usernameExists, phoneExists int
	err := r.db.QueryRowContext(
		ctx,
		query,
		nik,
		phone,
		email,
		username,
	).Scan(&nikExists, &emailExists, &usernameExists, &phoneExists)
	if err != nil {
		return nil, fmt.Errorf("check user existence: %w", err)
	}

	return &model.UserExistRespModel{
		NikExists:      nikExists == 1,
		EmailExists:    emailExists == 1,
		UsernameExists: usernameExists == 1,
		PhoneExists:    phoneExists == 1,
	}, nil
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
