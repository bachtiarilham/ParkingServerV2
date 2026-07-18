package petugas

import (
	dto "modulegue/internal/data/website/remote/dto/petugas"
	model "modulegue/internal/domain/web/model/petugas"
)

func ToPetugasRequestDto(src *model.PetugasRequestModel) *dto.PetugasRequestDto {
	if src == nil {
		return nil
	}
	return &dto.PetugasRequestDto{
		IDLokasi: src.IDLokasi,
	}
}

func ToPetugasRequestModel(src *dto.PetugasRequestDto) *model.PetugasRequestModel {
	if src == nil {
		return nil
	}
	return &model.PetugasRequestModel{
		IDLokasi: src.IDLokasi,
	}
}

func ToPetugasResponseDto(src *model.PetugasResponseModel) *dto.PetugasResponseDto {
	if src == nil {
		return nil
	}

	out := &dto.PetugasResponseDto{}
	if src.Petugas == nil {
		empty := []dto.PetugasItemDto{}
		out.Petugas = &empty
		return out
	}

	items := make([]dto.PetugasItemDto, 0, len(*src.Petugas))
	for _, item := range *src.Petugas {
		items = append(items, dto.PetugasItemDto{
			ID:              item.ID,
			Nama:            item.Nama,
			JmlTransaksi:    item.JmlTransaksi,
			TotalPendapatan: item.TotalPendapatan,
			IsAktif:         item.IsAktif,
			Parlok:          item.Parlok,
		})
	}
	out.Petugas = &items
	return out
}

func ToPetugasResponseModel(src *dto.PetugasResponseDto) *model.PetugasResponseModel {
	if src == nil {
		return nil
	}

	out := &model.PetugasResponseModel{}
	if src.Petugas == nil {
		empty := []model.PetugasItemModel{}
		out.Petugas = &empty
		return out
	}

	items := make([]model.PetugasItemModel, 0, len(*src.Petugas))
	for _, item := range *src.Petugas {
		items = append(items, model.PetugasItemModel{
			ID:              item.ID,
			Nama:            item.Nama,
			JmlTransaksi:    item.JmlTransaksi,
			TotalPendapatan: item.TotalPendapatan,
			IsAktif:         item.IsAktif,
			Parlok:          item.Parlok,
		})
	}
	out.Petugas = &items
	return out
}
