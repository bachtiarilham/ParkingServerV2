package paymentgate

import (
	dto "modulegue/internal/data/mobile/remote/dto/payment_gate"
	model "modulegue/internal/domain/mobile/model/payment_gate"
)

func ToCallbackRequestModel(src *dto.CallbackRequestDto) *model.CallbackRequestModel {
	return &model.CallbackRequestModel{
		TransactionTime:   src.TransactionTime,
		TransactionStatus: src.TransactionStatus,
		TransactionID:     src.TransactionID,
		StatusMessage:     src.StatusMessage,
		StatusCode:        src.StatusCode,
		SignatureKey:      src.SignatureKey,
		PaymentType:       src.PaymentType,
		OrderID:           src.OrderID,
		MerchantID:        src.MerchantID,
		GrossAmount:       src.GrossAmount,
		Currency:          src.Currency,
	}
}
