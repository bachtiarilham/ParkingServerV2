package auth

import (
	dto "modulegue/internal/data/shared/remote/dto/auth"
	model "modulegue/internal/domain/shared/model/auth"
)

func ToLoginReqModel(src *dto.LoginRequestDto) *model.LoginRequestModel {
	if src == nil {
		return nil
	}
	return &model.LoginRequestModel{
		Identity:   src.Identity,
		Password:   src.Password,
		DeviceId:   src.DeviceId,
		DeviceName: src.DeviceName,
		FcmToken:   src.FcmToken,
	}
}
