package payment

import (
	dto "modulegue/internal/data/mobile/remote/dto/payment_parking"
	model "modulegue/internal/domain/mobile/model/payment_parking"
)

func ToPostPaymentParkingResponseDto(src *model.PostPaymentParkingResponseModel) *dto.PostPaymentParkingResponseDto {
	if src == nil {
		return nil
	}

	return &dto.PostPaymentParkingResponseDto{
		SessionId:         src.SessionId,
		SessionCode:       src.SessionCode,
		TransactionCode:   src.TransactionCode,
		PlateNumber:       src.PlateNumber,
		VehicleTypeCode:   src.VehicleTypeCode,
		VehicleTypeName:   src.VehicleTypeName,
		LocationId:        src.LocationId,
		LocationName:      src.LocationName,
		AreaId:            src.AreaId,
		AreaName:          src.AreaName,
		Amount:            src.Amount,
		ParkingStatusCode: src.ParkingStatusCode,
		ParkingStatusName: src.ParkingStatusName,
		PaymentStatusCode: src.PaymentStatusCode,
		PaymentStatusName: src.PaymentStatusName,
		PaymentCode:       src.PaymentCode,
		FailedReason:      src.FailedReason,
		ReceiptNumber:     src.ReceiptNumber,
		StartedAt:         src.StartedAt,
		PaidAt:            src.PaidAt,
		QrExpiredAt:       src.QrExpiredAt,
	}
}
