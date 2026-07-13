package payment

import (
	dto "modulegue/internal/data/mobile/remote/dto/payment"
	model "modulegue/internal/domain/mobile/model/payment"
)

func ToPostPaymentParkingRequestModel(src *dto.PostPaymentParkingRequestDto) *model.PostPaymentParkingRequestModel {
	if src == nil {
		return nil
	}
	return &model.PostPaymentParkingRequestModel{
		SessionCode: src.SessionCode,
	}
}
