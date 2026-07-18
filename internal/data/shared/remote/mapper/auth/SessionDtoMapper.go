package auth

import (
	dto "modulegue/internal/data/shared/remote/dto/auth"
	model "modulegue/internal/domain/shared/model/auth"
)

func ToTokenSetDto(result *model.TokenSetModel) *dto.TokenSetDto {
	return &dto.TokenSetDto{
		AccessToken:      result.AccessToken,
		RefreshToken:     result.RefreshToken,
		ExpiresInSeconds: result.ExpiresIn,
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
