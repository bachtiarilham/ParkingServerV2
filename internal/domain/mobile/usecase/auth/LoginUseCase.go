package auth

import (
	"context"
	"fmt"
	"time"

	// "modulegue/internal/delivery/mobile/customer/dto"
	"modulegue/core/errorstring"
	"modulegue/core/hash"
	"modulegue/core/jwt"
	model "modulegue/internal/domain/mobile/model/auth"
	"modulegue/internal/domain/mobile/repository"
)

type LoginUseCase struct {
	authRepo     repository.AuthRepository
	sessionRepo  repository.SessionRepository
	sessionModel model.SessionModel
	jwtSecret    string
	accessTTL    time.Duration
	refreshTTL   time.Duration
}

func NewLoginUseCase(
	authRepo repository.AuthRepository,
	sessionRepo repository.SessionRepository,
	sessionModel model.SessionModel,
	jwtSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *LoginUseCase {
	return &LoginUseCase{
		authRepo:     authRepo,
		sessionRepo:  sessionRepo,
		sessionModel: sessionModel,
		jwtSecret:    jwtSecret,
		accessTTL:    accessTTL,
		refreshTTL:   refreshTTL,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, reqModel model.LoginRequestModel) (*model.LoginResponseModel, error) {
	if (reqModel.Identity == "") || reqModel.Password == "" {
		return nil, ErrInvalidInput
	}

	// 1. Cari user berdasarkan identity
	hasilLogin, err := uc.authRepo.LoginUser(ctx, reqModel.Identity, reqModel.Password)
	if err != nil {
		return nil, errorstring.ErrInvalidCredentials
	}

	// 2. Verifikasi password
	if err := hash.Compare(hasilLogin.PasswordHash, reqModel.Password); err != nil {
		return nil, errorstring.ErrInvalidCredentials
	}

	now := time.Now()
	accessExp := now.Add(uc.accessTTL)

	// 3. Buat access token
	accessToken, err := jwt.SignHS256(jwt.Claims{
		Subject:    fmt.Sprintf("%d", hasilLogin.UserId),
		UserID:     hasilLogin.UserId,
		Expiration: accessExp.Unix(),
		IssuedAt:   now.Unix(),
		Type:       "access",
	}, uc.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	// 4. Buat refresh token (TTL lebih panjang, tanpa expiry di payload — dicek dari DB)
	refreshExp := now.Add(uc.refreshTTL)
	refreshToken, err := jwt.SignHS256(jwt.Claims{
		Subject:    fmt.Sprintf("%d", hasilLogin.UserId),
		UserID:     hasilLogin.UserId,
		RoleID:     hasilLogin.RoleId,
		TokenType:  "JWT",
		Expiration: refreshExp.Unix(),
		IssuedAt:   now.Unix(),
		Type:       "refresh",
	}, uc.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	// 5. Simpan session ke DB
	if err := uc.sessionRepo.SaveSession(ctx, model.SessionModel{
		UserID:       hasilLogin.UserId,
		TokenType:    "JWT",
		RefreshToken: refreshToken,
		ExpiresAt:    refreshExp,
	}); err != nil {
		return nil, fmt.Errorf("save session: %w", err)
	}

	return &model.LoginResponseModel{
		UserModel: model.UserModel{
			UserId:       hasilLogin.UserId,
			Nik:          hasilLogin.Nik,
			FullName:     hasilLogin.FullName,
			Phone:        hasilLogin.Phone,
			Email:        hasilLogin.Email,
			Username:     hasilLogin.Username,
			Password:     "",
			PasswordHash: "",
			RoleId:       hasilLogin.RoleId,
			IsVerified:   hasilLogin.IsVerified,
			Lokasi:       hasilLogin.Lokasi,
			Zona:         hasilLogin.Zona,
			Tarif:        hasilLogin.Tarif,
			RegisteredAt: hasilLogin.RegisteredAt,
			CreatedAt:    hasilLogin.CreatedAt,
			UpdatedAt:    hasilLogin.UpdatedAt,
		},
		TokenSetModel: model.TokenSetModel{
			AccessToken:      accessToken,
			RefreshToken:     refreshToken,
			TokenType:        "Bearer",
			ExpiresInSeconds: accessExp.Unix(),
		},
	}, nil
}
