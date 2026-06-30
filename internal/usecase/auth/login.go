package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	// "modulegue/internal/delivery/mobile/customer/dto"
	"modulegue/internal/domain/auth"
	"modulegue/internal/domain/user"
	"modulegue/pkg/hash"
	"modulegue/pkg/jwt"
)

var (
	ErrInvalidCredentials = errors.New("email atau password salah")
	ErrUserNotFound       = errors.New("user tidak ditemukan")
)

type LoginInput struct {
	Identity string
	Password string
}

type LoginOuput struct {
	UserID       int64
	FullName     string
	RoleID       int64
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type LoginUseCase struct {
	userRepo   user.Repository
	authRepo   auth.Repository
	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewLoginUseCase(
	userRepo user.Repository,
	authRepo auth.Repository,
	jwtSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:   userRepo,
		authRepo:   authRepo,
		jwtSecret:  jwtSecret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, req LoginInput) (LoginOuput, error) {
	// 1. Cari user berdasarkan identity
	u, err := uc.userRepo.FindByEmail(ctx, req.Identity)
	if err != nil {
		return LoginOuput{}, ErrInvalidCredentials
	}

	// 2. Verifikasi password
	if err := hash.Compare(u.PasswordHash, req.Password); err != nil {
		return LoginOuput{}, ErrInvalidCredentials
	}

	now := time.Now()
	accessExp := now.Add(uc.accessTTL)

	// 3. Buat access token
	accessToken, err := jwt.SignHS256(jwt.Claims{
		Subject:    fmt.Sprintf("%d", u.ID),
		UserID:     u.ID,
		Expiration: accessExp.Unix(),
		IssuedAt:   now.Unix(),
		Type:       "access",
	}, uc.jwtSecret)
	if err != nil {
		return LoginOuput{}, fmt.Errorf("generate access token: %w", err)
	}

	// 4. Buat refresh token (TTL lebih panjang, tanpa expiry di payload — dicek dari DB)
	refreshExp := now.Add(uc.refreshTTL)
	refreshToken, err := jwt.SignHS256(jwt.Claims{
		Subject:    fmt.Sprintf("%d", u.ID),
		UserID:     u.ID,
		RoleID:     u.RoleID,
		TokenType:  "refresh",
		Expiration: refreshExp.Unix(),
		IssuedAt:   now.Unix(),
		Type:       "refresh",
	}, uc.jwtSecret)
	if err != nil {
		return LoginOuput{}, fmt.Errorf("generate refresh token: %w", err)
	}

	// 5. Simpan session ke DB
	if err := uc.authRepo.SaveSession(ctx, auth.Session{
		UserID:       u.ID,
		TokenType:    "JWT",
		RefreshToken: refreshToken,
		ExpiresAt:    refreshExp,
	}); err != nil {
		return LoginOuput{}, fmt.Errorf("save session: %w", err)
	}

	return LoginOuput{
		UserID:       u.ID,
		FullName:     u.FullName,
		RoleID:       u.RoleID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExp.Unix(),
	}, nil
}
