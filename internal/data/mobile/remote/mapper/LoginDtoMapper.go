package mapper

import (
	"modulegue/internal/data/mobile/remote/dto"
	authentity "modulegue/internal/domain/auth"
)

func ToLoginRequestModel(src *dto.LoginRequestDto) *authentity.Credential {
	if src == nil {
		return nil
	}
	return &authentity.Credential{Email: src.Identity, PasswordHash: src.Password}
}

func ToLoginResponseDto(user dto.UserDto, token dto.TokenSetDto) *dto.LoginResponseDto {
	return &dto.LoginResponseDto{
		UserDto:     user,
		TokenSetDto: token,
	}
}
