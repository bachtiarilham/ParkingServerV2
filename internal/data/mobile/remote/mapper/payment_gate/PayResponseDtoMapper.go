package paymentgate

import (
	dto "modulegue/internal/data/mobile/remote/dto/payment_gate"
	model "modulegue/internal/domain/mobile/model/payment_gate"
)

func ToPayResponseDto(src *model.PayResponseModel) *dto.PayResponseDto {
	if src == nil {
		return nil
	}

	var paymentActionDto *dto.PaymentActionDto
	if src.PaymentAction != nil {
		paymentActionDto = ToPaymentActionDto(src.PaymentAction)
	}

	return &dto.PayResponseDto{
		OrderID:       src.OrderID,
		GrossAmount:   src.GrossAmount,
		PaymentMethod: src.PaymentMethod,
		Status:        src.Status,
		ExpiryTime:    src.ExpiryTime,
		SnapToken:     src.SnapToken,
		RedirectURL:   src.RedirectURL,
		PaymentAction: paymentActionDto,
	}
}

func ToPaymentActionDto(src *model.PaymentActionModel) *dto.PaymentActionDto {
	if src == nil {
		return nil
	}
	return &dto.PaymentActionDto{
		DeepLinkURL:  src.DeepLinkURL,
		StoreName:    src.StoreName,
		PaymentCode:  src.PaymentCode,
		QRCodeString: src.QRCodeString,
		QRCodeURL:    src.QRCodeURL,
		Instruction:  src.Instruction,
	}
}
