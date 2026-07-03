package auth

import (
	dto "modulegue/internal/data/mobile/remote/dto/auth"
	model "modulegue/internal/domain/mobile/model/auth"
)

func ToTokenSetDto(result *model.TokenSetModel) *dto.TokenSetDto {
	return &dto.TokenSetDto{
		AccessToken:      result.AccessToken,
		RefreshToken:     result.RefreshToken,
		TokenType:        "Bearer",
		ExpiresInSeconds: result.ExpiresInSeconds,
	}
}

// func ToSessionModel(result *model.SessionModel) *model.SessionModel {
// 	return &model.SessionModel{
// 		UserID:       result.UserID,
// 		RefreshToken: result.RefreshToken,
// 		TokenType:    result.TokenType,
// 		ExpiresAt:    result.ExpiresAt,
// 	}
// }
