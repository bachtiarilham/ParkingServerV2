package payment

import (
	dto "modulegue/internal/data/mobile/remote/dto/payment_parking"
	model "modulegue/internal/domain/mobile/model/payment_parking"
)

func ToPostPaymentParkingRequestModel(src *dto.PostPaymentParkingRequestDto) *model.PostPaymentParkingRequestModel {
	if src == nil {
		return nil
	}
	return &model.PostPaymentParkingRequestModel{
		SessionCode: src.SessionCode,
	}
}
