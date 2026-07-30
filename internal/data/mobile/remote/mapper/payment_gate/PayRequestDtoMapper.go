package paymentgate

import (
	dto "modulegue/internal/data/mobile/remote/dto/payment_gate"
	model "modulegue/internal/domain/mobile/model/payment_gate"
)

func ToPayRequestModel(src *dto.PayRequestDto) *model.PayRequestModel {
	return &model.PayRequestModel{
		PaymentType:       src.PaymentType,
		TargetID:          src.TargetID,
		PaymentMethodCode: src.PaymentMethodCode,
		Amount:            src.Amount,
		PromoCode:         src.PromoCode,
	}
}
