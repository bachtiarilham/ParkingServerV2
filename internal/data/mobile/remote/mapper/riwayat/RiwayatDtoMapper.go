package riwayat

import (
	dto "modulegue/internal/data/mobile/remote/dto/riwayat"
	model "modulegue/internal/domain/mobile/model/riwayat"
)

func ToRiwayatDto(src *model.RiwayatModel) *dto.RiwayatResponseDto {
	if src == nil {
		return nil
	}
	out := &dto.RiwayatResponseDto{Sections: make([]dto.RiwayatSectionDto, 0, len(src.Sections))}
	for _, item := range src.Sections {
		out.Sections = append(out.Sections, *ToRiwayatSectionDto(&item))
	}
	return out
}

func ToRiwayatSectionDto(src *model.RiwayatSectionModel) *dto.RiwayatSectionDto {
	if src == nil {
		return nil
	}
	out := &dto.RiwayatSectionDto{Items: make([]dto.RiwayatItemDto, 0, len(src.Items))}
	out.Date = src.Date
	for _, item := range src.Items {
		out.Items = append(out.Items, *ToRiwayatItemDto(&item))
	}
	return out
}

func ToRiwayatItemDto(src *model.RiwayatItemModel) *dto.RiwayatItemDto {
	if src == nil {
		return nil
	}
	return &dto.RiwayatItemDto{
		Code:        src.Code,
		PlateNumber: src.PlateNumber,
		VehicleType: src.VehicleType,
		Time:        src.Time,
		Amount:      src.Amount,
		IsEntry:     src.IsEntry,
	}
}
