package mapper

import (
	"modulegue/internal/data/mobile/remote/dto"
	"modulegue/internal/domain/mobile/model"
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
		Role:       src.Role,
		IsVerified: src.IsVerified,
		Lokasi:     src.Lokasi,
		Zona:       src.Zona,
		Tarif:      ToTarifDtos(src.Tarif),
	}
}
