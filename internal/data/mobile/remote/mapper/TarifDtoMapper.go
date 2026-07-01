package mapper

import (
	"modulegue/internal/data/mobile/remote/dto"
	"modulegue/internal/domain/mobile/model"
)

func ToTarifDto(src *model.TarifModel) *dto.TarifDto {
	if src == nil {
		return nil
	}
	return &dto.TarifDto{Kendaraan: src.Kendaraan, Nominal: src.Nominal}
}

func ToTarifDtos(src []model.TarifModel) []dto.TarifDto {
	out := make([]dto.TarifDto, 0, len(src))
	for _, item := range src {
		out = append(out, *ToTarifDto(&item))
	}
	return out
}
