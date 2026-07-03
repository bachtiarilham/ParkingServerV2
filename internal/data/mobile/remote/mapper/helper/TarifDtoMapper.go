package helper

import (
	dto "modulegue/internal/data/mobile/remote/dto/helper"
	model "modulegue/internal/domain/mobile/model/helper"
)

func ToTarifDto(src *model.TarifModel) *dto.TarifDto {
	if src == nil {
		return nil
	}

	out := &dto.TarifDto{
		ItemTarif: make([]dto.TarifItemDto, 0, len(src.ItemTarif)),
	}
	for _, item := range src.ItemTarif {
		out.ItemTarif = append(out.ItemTarif, dto.TarifItemDto{
			Kendaraan: item.Kendaraan,
			Nominal:   item.Nominal,
		})
	}
	return out
}

func ToTarifDtos(src []model.TarifModel) []dto.TarifDto {
	out := make([]dto.TarifDto, 0, len(src))
	for _, item := range src {
		mapped := ToTarifDto(&item)
		if mapped != nil {
			out = append(out, *mapped)
		}
	}
	return out
}
