package auth

import (
	dto "modulegue/internal/data/mobile/remote/dto/auth"
	model "modulegue/internal/domain/mobile/model/auth"
)

func ToLoginRespDto(src *model.TokenSetModel, roleId int64) *dto.LoginRespDto {
	if src == nil || roleId == 0 {
		return nil
	}
	return &dto.LoginRespDto{
		TokenSetDto: ToTokenSetDto(src),
		RoleId:      roleId,
	}
}
