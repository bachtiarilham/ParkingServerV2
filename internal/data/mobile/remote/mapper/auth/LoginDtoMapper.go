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
		Identity:   src.Identity,
		Password:   src.Password,
		DeviceId:   src.DeviceId,
		DeviceName: src.DeviceName,
		FcmToken:   src.FcmToken,
	}
}
