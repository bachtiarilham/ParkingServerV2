package auth

import (
	dto "modulegue/internal/data/mobile/remote/dto/auth"
	mapper "modulegue/internal/data/mobile/remote/mapper/helper"
	model "modulegue/internal/domain/mobile/model/auth"
)

func ToUserDto(src *model.UserModel) *dto.UserDto {
	if src == nil {
		return nil
	}
	return &dto.UserDto{
		UserId:     src.UserId,
		Nik:        src.Nik,
		FullName:   src.FullName,
		Phone:      src.Phone,
		Email:      src.Email,
		Username:   src.Username,
		Password:   src.Password,
		RoleId:     src.RoleId,
		IsVerified: src.IsVerified,
		Lokasi:     src.Lokasi,
		Zona:       src.Zona,
		Tarif:      mapper.ToTarifDtos(src.Tarif),
	}
}
