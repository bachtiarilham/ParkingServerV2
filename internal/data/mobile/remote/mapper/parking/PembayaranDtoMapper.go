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
		SessionId:         src.SessionId,
		SessionCode:       src.SessionCode,
		TransactionCode:   src.TransactionCode,
		PlateNumber:       src.PlateNumber,
		VehicleTypeCode:   src.VehicleTypeCode,
		VehicleTypeName:   src.VehicleTypeName,
		ZoneId:            src.ZoneId,
		ZoneName:          src.ZoneName,
		LocationId:        src.LocationId,
		LocationName:      src.LocationName,
		Address:           src.Address,
		AreaId:            src.AreaId,
		AreaName:          src.AreaName,
		Amount:            src.Amount,
		QrString:          src.QrString,
		QrExpiredAt:       src.QrExpiredAt,
		PaymentCode:       src.PaymentCode,
		PaymentStatusCode: src.PaymentStatusCode,
		PaymentStatusName: src.PaymentStatusName,
	}
	return out
}
