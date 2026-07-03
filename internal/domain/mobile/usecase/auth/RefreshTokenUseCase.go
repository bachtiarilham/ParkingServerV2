package auth

import (
	"context"
	"fmt"
	"modulegue/core/errorstring"
	"modulegue/core/jwt"
	model "modulegue/internal/domain/mobile/model/auth"
	"modulegue/internal/domain/mobile/repository"
	"time"
)

type RefreshTokenUseCase struct {
	repository repository.AuthRepository
	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewRefreshTokenUseCase(
	repository repository.AuthRepository,
	jwtSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		repository: repository,
		jwtSecret:  jwtSecret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (uc *RefreshTokenUseCase) Execute(ctx context.Context, reqModel model.RefreshTokenReqModel) (*model.RefreshTokenRespModel, error) {
	if reqModel.RefreshToken == "" {
		return nil, errorstring.ErrInvalidRefreshToken
	}

	// 1. Cari session berdasarkan refresh token
	session, err := uc.repository.FindSessionByRefreshToken(ctx, reqModel.RefreshToken)
	if err != nil {
		return nil, errorstring.ErrInvalidRefreshToken
	}

	// 2. Cek apakah refresh token sudah kadaluarsa
	if time.Now().After(session.ExpiresAt) {
		// Hapus session lama jika kadaluarsa
		uc.repository.DeleteSession(ctx, reqModel.RefreshToken) // Log error jika perlu, tapi jangan hentikan proses
		return nil, errorstring.ErrExpiredRefreshToken
	}

	claims, err := jwt.ParseClaimsHS256(reqModel.RefreshToken, uc.jwtSecret)
	if err != nil {
		return nil, errorstring.ErrInvalidRefreshToken
	}
	if claims.Type != "refresh" && claims.TokenType != "refresh" {
		return nil, errorstring.ErrInvalidRefreshToken
	}

	// 3. Generate token baru
	now := time.Now()
	accessExp := now.Add(uc.accessTTL)
	newRefreshExp := now.Add(uc.refreshTTL)

	newAccessToken, err := jwt.SignHS256(jwt.Claims{
		Subject:    fmt.Sprintf("%d", session.UserID),
		UserID:     session.UserID,
		RoleID:     claims.RoleID,
		Expiration: accessExp.Unix(),
		IssuedAt:   now.Unix(),
		Type:       "access",
	}, uc.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("generate new access token: %w", err)
	}

	newRefreshToken, err := jwt.SignHS256(jwt.Claims{
		Subject:    fmt.Sprintf("%d", session.UserID),
		UserID:     session.UserID,
		RoleID:     claims.RoleID,
		TokenType:  "JWT",
		Expiration: newRefreshExp.Unix(),
		IssuedAt:   now.Unix(),
		Type:       "refresh",
	}, uc.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("generate new refresh token: %w", err)
	}

	newSession := model.SessionModel{
		UserID:       session.UserID,
		TokenType:    "JWT",
		RefreshToken: newRefreshToken,
		ExpiresAt:    newRefreshExp,
		UpdatedAt:    now,
		CreatedAt:    now,
	}
	if err := uc.repository.RotateSession(ctx, reqModel.RefreshToken, newSession); err != nil {
		return nil, fmt.Errorf("rotate session: %w", err)
	}

	return &model.RefreshTokenRespModel{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    accessExp.Unix(),
	}, nil
}
