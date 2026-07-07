package auth

import (
	dto "modulegue/internal/data/mobile/remote/dto/auth"
	model "modulegue/internal/domain/mobile/model/auth"
)

func ToRegisterRequestModel(src *dto.RegisterRequestDto) *model.RegisterRequestModel {
	if src == nil {
		return nil
	}
	return &model.RegisterRequestModel{FullName: src.FullName, NIK: src.NIK, Phone: src.Phone, Email: src.Email, Username: src.Username, Password: src.Password}
}

func ToRegisterResponseDto(src *model.RegisterResponseModel) *dto.RegisterResponseDto {
	if src == nil {
		return nil
	}
	return &dto.RegisterResponseDto{Message: src.Message}
}
