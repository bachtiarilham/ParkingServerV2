package helper

import (
	dto "modulegue/internal/data/mobile/remote/dto/helper"
	model "modulegue/internal/domain/mobile/model/helper"
)

func ToLokasiDto(src *model.LokasiModel) *dto.LokasiResponseDto {
	if src == nil {
		return nil
	}
	return &dto.LokasiResponseDto{Lokasi: append([]string(nil), src.Lokasi...)}
}
