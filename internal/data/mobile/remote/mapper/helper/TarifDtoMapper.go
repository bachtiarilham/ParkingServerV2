package helper

import (
	dto "modulegue/internal/data/mobile/remote/dto/helper"
	model "modulegue/internal/domain/mobile/model/helper"
)

func ToTarifDto(src *model.TarifModel) *dto.TarifResponseDto {
	if src == nil {
		return nil
	}

	items := make([]dto.TarifResponseItemDto, 0)
	if src.TarifItem != nil {
		for _, item := range *src.TarifItem {
			items = append(items, dto.TarifResponseItemDto{
				KendaraanId:   item.KendaraanId,
				KendaraanKode: item.KendaraanKode,
				KendaraanNama: item.KendaraanNama,
				Nominal:       item.Nominal,
			})
		}
	}

	return &dto.TarifResponseDto{
		TarifResponseItemDto: items,
	}
}
