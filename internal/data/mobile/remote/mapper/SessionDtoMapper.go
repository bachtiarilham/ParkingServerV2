package mapper

import (
	"time"

	"modulegue/internal/data/mobile/remote/dto"
	authentity "modulegue/internal/domain/auth"
)

func ToTokenSetDto(accessToken, refreshToken string, expiresInSeconds int64) dto.TokenSetDto {
	return dto.TokenSetDto{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		ExpiresInSeconds: expiresInSeconds,
	}
}

func ToSessionModel(userID int64, refreshToken string, tokenType string, expiresAt time.Time) authentity.Session {
	return authentity.Session{
		UserID:       userID,
		RefreshToken: refreshToken,
		TokenType:    tokenType,
		ExpiresAt:    expiresAt,
	}
}
