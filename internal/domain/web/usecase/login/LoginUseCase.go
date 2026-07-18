package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	// "modulegue/internal/delivery/mobile/customer/dto"
	"modulegue/core/errorstring"
	"modulegue/core/hash"
	"modulegue/core/jwt"
	modelauth "modulegue/internal/domain/shared/model/auth"
	repoauth "modulegue/internal/domain/shared/repository"
	modelresp "modulegue/internal/domain/web/model/home"
	repohome "modulegue/internal/domain/web/repository"
)

type LoginUseCase struct {
	loginRepo    repohome.LoginRepository
	sessionRepo  repoauth.SessionRepository
	sessionModel modelauth.SessionModel
	jwtSecret    string
	accessTTL    time.Duration
	refreshTTL   time.Duration
}

func NewLoginUseCase(
	loginRepo repohome.LoginRepository,
	sessionRepo repoauth.SessionRepository,
	sessionModel modelauth.SessionModel,
	jwtSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *LoginUseCase {
	return &LoginUseCase{
		loginRepo:    loginRepo,
		sessionRepo:  sessionRepo,
		sessionModel: sessionModel,
		jwtSecret:    jwtSecret,
		accessTTL:    accessTTL,
		refreshTTL:   refreshTTL,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, reqModel modelauth.LoginRequestModel) (*modelauth.TokenSetModel, *modelresp.HomeResponseModel, error) {
	if (reqModel.Identity == "") || reqModel.Password == "" {
		return nil, nil, errors.New("tolong isi field yang wajib")
	}

	_, hasilLogin, err := uc.loginRepo.Login(ctx, reqModel)
	if err != nil {
		return nil, nil, errorstring.ErrInvalidCredentials
	}

	// 2. Verifikasi password
	if err := hash.Compare(hasilLogin.PasswordHash, reqModel.Password); err != nil {
		return nil, nil, errorstring.ErrInvalidCredentials
	}

	now := time.Now()
	accessExp := now.Add(uc.accessTTL)

	// 3. Buat access token
	accessToken, err := jwt.SignHS256(jwt.Claims{
		Subject:    fmt.Sprintf("%d", hasilLogin.UserId),
		UserID:     hasilLogin.UserId,
		RoleID:     hasilLogin.RoleId,
		Expiration: accessExp.Unix(),
		IssuedAt:   now.Unix(),
		Type:       "access",
	}, uc.jwtSecret)
	if err != nil {
		return nil, nil, fmt.Errorf("generate access token: %w", err)
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
		return nil, nil, fmt.Errorf("generate refresh token: %w", err)
	}

	// 5. Simpan session ke DB
	if err := uc.sessionRepo.SaveSession(ctx, modelauth.SessionModel{
		UserID:       hasilLogin.UserId,
		RefreshToken: refreshToken,
		DeviceId:     reqModel.DeviceId,
		DeviceName:   reqModel.DeviceName,
		FcmToken:     reqModel.FcmToken,
		IpAddress:    reqModel.IpAdrress,
		UserAgent:    reqModel.UserAgent,
		ExpiresAt:    refreshExp,
	}); err != nil {
		return nil, nil, fmt.Errorf("save session: %w", err)
	}

	// 6. simpan last login
	if err := uc.sessionRepo.SaveLastLogin(ctx, hasilLogin.UserId); err != nil {
		return nil, nil, fmt.Errorf("save last login: %w", err)
	}

	hasilLogin.Token = accessToken
	return nil, hasilLogin, nil
}
