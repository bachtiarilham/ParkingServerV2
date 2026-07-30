package riwayat

import (
	dto "modulegue/internal/data/mobile/remote/dto/riwayat"
	model "modulegue/internal/domain/mobile/model/riwayat"
)

func ToRiwayatParkirDto(src *model.RiwayatParkirModel) *dto.RiwayatParkirDto {
	if src == nil {
		return nil
	}

	items := make([]dto.RiwayatParkirItemDto, len(src.Items))
	for i, item := range src.Items {
		items[i] = dto.RiwayatParkirItemDto{
			ID:           item.ID,
			TicketNo:     item.TicketNo,
			LocationName: item.LocationName,
			ParkingType:  item.ParkingType,
			LicensePlate: item.LicensePlate,
			VehicleType:  item.VehicleType,
			CheckInTime:  item.CheckInTime,
			CheckOutTime: item.CheckOutTime,
			DurationText: item.DurationText,
			Amount:       item.Amount,
			Status:       item.Status,
		}
	}

	return &dto.RiwayatParkirDto{
		Summary: dto.ParkingSummaryDto{
			TotalDurationMinutes: src.Summary.TotalDurationMinutes,
			TotalDurationText:    src.Summary.TotalDurationText,
			TotalAmount:          src.Summary.TotalAmount,
			CompletedCount:       src.Summary.CompletedCount,
		},
		Items: items,
	}
}
