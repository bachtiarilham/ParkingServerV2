package auth

import (
	dto "modulegue/internal/data/shared/remote/dto/auth"
	model "modulegue/internal/domain/shared/model/auth"
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
