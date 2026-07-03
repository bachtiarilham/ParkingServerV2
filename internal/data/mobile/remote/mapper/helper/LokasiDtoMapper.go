package mapper

import (
	dto "modulegue/internal/data/mobile/remote/dto/helper"
	model "modulegue/internal/domain/mobile/model/helper"
)

func ToLokasiDto(src *model.LokasiModel) *dto.LokasiDto {
	if src == nil {
		return nil
	}
	return &dto.LokasiDto{Lokasi: append([]string(nil), src.Lokasi...)}
}
