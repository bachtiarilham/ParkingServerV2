package payment

import (
	dto "modulegue/internal/data/mobile/remote/dto/parking"
	model "modulegue/internal/domain/mobile/model/parking"
)

func ToParkingResponseDto(src *model.PostParkingResponseModel) *dto.PostParkingResponseDto {
	if src == nil {
		return nil
	}
	out := &dto.PostParkingResponseDto{
		SessionCode: src.SessionCode,
		PlateNumber: src.PlateNumber,
		Waktu:       src.Waktu,
		QrExpired:   src.QrExpired,
		BiayaParkir: src.BiayaParkir,
	}
	return out
}
