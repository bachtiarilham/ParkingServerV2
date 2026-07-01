package mapper

import (
	"modulegue/internal/data/mobile/remote/dto"
	"modulegue/internal/domain/mobile/model"
)

func ToLokasiDto(src *model.LokasiModel) *dto.LokasiDto {
	if src == nil {
		return nil
	}
	return &dto.LokasiDto{Lokasi: append([]string(nil), src.Lokasi...)}
}
