package riwayat

import (
	dto "modulegue/internal/data/mobile/remote/dto/riwayat"
	model "modulegue/internal/domain/mobile/model/riwayat"
)

func ToRiwayatDto(src *model.RiwayatModel) *dto.RiwayatResponseDto {
	if src == nil {
		return nil
	}
	out := &dto.RiwayatResponseDto{
		ParkingSections: make([]dto.ParkingSectionDto, 0, len(src.ParkirSections)),
		TopUpSections:   make([]dto.TopUpSectionDto, 0, len(src.TopUpSections)),
	}
	for _, item := range src.ParkirSections {
		out.ParkingSections = append(out.ParkingSections, *ToRiwayatSectionDto(&item))
	}
	for _, item := range src.TopUpSections {
		out.TopUpSections = append(out.TopUpSections, *ToTopUpSectionDto(&item))
	}
	return out
}

func ToRiwayatSectionDto(src *model.RiwayatSectionModel) *dto.ParkingSectionDto {
	if src == nil {
		return nil
	}
	out := &dto.ParkingSectionDto{Items: make([]dto.ParkingItemDto, 0, len(src.Items))}
	out.Date = src.Date
	for _, item := range src.Items {
		out.Items = append(out.Items, *ToRiwayatItemDto(&item))
	}
	return out
}

func ToRiwayatItemDto(src *model.RiwayatItemModel) *dto.ParkingItemDto {
	if src == nil {
		return nil
	}
	return &dto.ParkingItemDto{
		Code:        src.Code,
		PlateNumber: src.PlateNumber,
		VehicleType: src.VehicleType,
		Time:        src.Time,
		Amount:      src.Amount,
		IsEntry:     src.IsEntry,
	}
}

func ToTopUpSectionDto(src *model.TopUpSectionModel) *dto.TopUpSectionDto {
	if src == nil {
		return nil
	}
	out := &dto.TopUpSectionDto{Items: make([]dto.TopUpItemDto, 0, len(src.Items))}
	out.Date = src.Date
	for _, item := range src.Items {
		out.Items = append(out.Items, *ToTopUpItemDto(&item))
	}
	return out
}

func ToTopUpItemDto(src *model.TopUpItemModel) *dto.TopUpItemDto {
	if src == nil {
		return nil
	}
	return &dto.TopUpItemDto{
		Code:              src.Code,
		PaymentMethodName: src.PaymentMethodName,
		TransactionStatus: src.TransactionStatus,
		ProviderName:      src.ProviderName,
		Time:              src.Time,
		Amount:            src.Amount,
	}
}
