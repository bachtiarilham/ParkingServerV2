package mapper

import (
	"modulegue/internal/data/mobile/remote/dto"
	"modulegue/internal/domain/mobile/model"
)

func ToRiwayatDto(src *model.RiwayatModel) *dto.RiwayatDto {
	if src == nil {
		return nil
	}
	out := &dto.RiwayatDto{Sections: make([]dto.RiwayatSectionDto, 0, len(src.Sections))}
	for _, item := range src.Sections {
		out.Sections = append(out.Sections, *ToRiwayatSectionDto(&item))
	}
	return out
}

func ToRiwayatSectionDto(src *model.RiwayatSectionModel) *dto.RiwayatSectionDto {
	if src == nil {
		return nil
	}
	out := &dto.RiwayatSectionDto{Date: src.Date, Items: make([]dto.RiwayatItemDto, 0, len(src.Items))}
	for _, item := range src.Items {
		out.Items = append(out.Items, *ToRiwayatItemDto(&item))
	}
	return out
}

func ToRiwayatItemDto(src *model.RiwayatItemModel) *dto.RiwayatItemDto {
	if src == nil {
		return nil
	}
	return &dto.RiwayatItemDto{Code: src.Code, PlateNumber: src.PlateNumber, VehicleType: src.VehicleType, Time: src.Time, Amount: src.Amount, IsEntry: src.IsEntry}
}
