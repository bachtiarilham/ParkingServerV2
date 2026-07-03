package auth

import (
	dto "modulegue/internal/data/mobile/remote/dto/auth"
	model "modulegue/internal/domain/mobile/model/auth"
)

func ToLoginRequestModel(src *dto.LoginRequestDto) *model.LoginRequestModel {
	if src == nil {
		return nil
	}
	return &model.LoginRequestModel{
		Identity: src.Identity,
		Password: src.Password,
	}
}

func ToLoginResponseDto(src *model.LoginResponseModel) *dto.LoginResponseDto {
	if src == nil {
		return nil
	}

	return &dto.LoginResponseDto{
		UserDto: dto.UserDto{
			UserId:     src.UserModel.UserId,
			Nik:        src.UserModel.Nik,
			FullName:   src.UserModel.FullName,
			Phone:      src.UserModel.Phone,
			Email:      src.UserModel.Email,
			Username:   src.UserModel.Username,
			Password:   src.UserModel.Password,
			RoleId:     src.UserModel.RoleId,
			IsVerified: src.UserModel.IsVerified,
			Lokasi:     src.UserModel.Lokasi,
			Zona:       src.UserModel.Zona,
			Tarif:      nil,
		},
		TokenSetDto: dto.TokenSetDto{
			AccessToken:      src.TokenSetModel.AccessToken,
			RefreshToken:     src.TokenSetModel.RefreshToken,
			TokenType:        src.TokenSetModel.TokenType,
			ExpiresInSeconds: src.TokenSetModel.ExpiresInSeconds,
		},
	}
}
