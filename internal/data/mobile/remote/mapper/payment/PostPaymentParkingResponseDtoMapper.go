package payment

import (
	dto "modulegue/internal/data/mobile/remote/dto/payment"
	model "modulegue/internal/domain/mobile/model/payment"
)

func ToPostPaymentParkingResponseDto(src *model.PostPaymentParkingResponseModel) *dto.PostPaymentParkingResponseDto {
	if src == nil {
		return nil
	}

	result := &dto.PostPaymentParkingResponseDto{
		Title:               src.Title,
		SuccessTitle:        src.SuccessTitle,
		SuccessDescription:  src.SuccessDescription,
		TotalAmount:         src.TotalAmount,
		PaymentStatus:       src.PaymentStatus,
		ReferenceNumber:     src.ReferenceNumber,
		VerificationMessage: src.VerificationMessage,
		ThankYouTitle:       src.ThankYouTitle,
		ThankYouDescription: src.ThankYouDescription,
		DownloadLabel:       src.DownloadLabel,
		BackToHomeLabel:     src.BackToHomeLabel,
		Details:             []dto.PostPaymentParkingDetailItemDto{},
	}

	for _, item := range src.Details {
		result.Details = append(result.Details, dto.PostPaymentParkingDetailItemDto{
			Label: item.Label,
			Value: item.Value,
		})
	}

	return result
}
