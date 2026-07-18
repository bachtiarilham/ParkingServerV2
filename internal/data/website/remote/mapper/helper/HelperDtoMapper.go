package helper

import (
	dto "modulegue/internal/data/website/remote/dto/helper"
	model "modulegue/internal/domain/web/model/helper"
)

func ToGetLokasiRequestDto(src *model.GetLokasiRequestModel) *dto.GetLokasiRequestDto {
	if src == nil {
		return nil
	}
	return &dto.GetLokasiRequestDto{
		Action: src.Action,
	}
}

func ToGetLokasiRequestModel(src *dto.GetLokasiRequestDto) *model.GetLokasiRequestModel {
	if src == nil {
		return nil
	}
	return &model.GetLokasiRequestModel{
		Action: src.Action,
	}
}

func ToGetLokasiResponseDto(src *model.GetLokasiResponseModel) *dto.GetLokasiResponseDto {
	if src == nil {
		return nil
	}

	out := &dto.GetLokasiResponseDto{}
	if src.LokasiItem == nil {
		empty := []dto.LokasiItemDto{}
		out.LokasiItem = &empty
		return out
	}

	items := make([]dto.LokasiItemDto, 0, len(*src.LokasiItem))
	for _, item := range *src.LokasiItem {
		items = append(items, ToLokasiItemDto(item))
	}
	out.LokasiItem = &items
	return out
}

func ToGetLokasiResponseModel(src *dto.GetLokasiResponseDto) *model.GetLokasiResponseModel {
	if src == nil {
		return nil
	}

	out := &model.GetLokasiResponseModel{}
	if src.LokasiItem == nil {
		empty := []model.LokasiItemModel{}
		out.LokasiItem = &empty
		return out
	}

	items := make([]model.LokasiItemModel, 0, len(*src.LokasiItem))
	for _, item := range *src.LokasiItem {
		items = append(items, ToLokasiItemModel(item))
	}
	out.LokasiItem = &items
	return out
}

func ToLokasiItemDto(src model.LokasiItemModel) dto.LokasiItemDto {
	return dto.LokasiItemDto{
		ID:          src.ID,
		NamaParlok:  src.NamaParlok,
		JalanParlok: src.JalanParlok,
	}
}

func ToLokasiItemModel(src dto.LokasiItemDto) model.LokasiItemModel {
	return model.LokasiItemModel{
		ID:          src.ID,
		NamaParlok:  src.NamaParlok,
		JalanParlok: src.JalanParlok,
	}
}

func ToGetTarifRequestDto(src *model.GetTarifRequestModel) *dto.GetTarifRequestDto {
	if src == nil {
		return nil
	}
	return &dto.GetTarifRequestDto{
		IDLokasi: src.IDLokasi,
	}
}

func ToGetTarifRequestModel(src *dto.GetTarifRequestDto) *model.GetTarifRequestModel {
	if src == nil {
		return nil
	}
	return &model.GetTarifRequestModel{
		IDLokasi: src.IDLokasi,
	}
}

func ToGetTarifResponseDto(src *model.GetTarifResponseModel) *dto.GetTarifResponseDto {
	if src == nil {
		return nil
	}

	if src.Tarif == nil {
		return &dto.GetTarifResponseDto{Tarif: []dto.TarifItemDto{}}
	}

	items := make([]dto.TarifItemDto, 0, len(*src.Tarif))
	for _, item := range *src.Tarif {
		items = append(items, ToTarifItemDto(item))
	}
	return &dto.GetTarifResponseDto{
		Tarif: items,
	}
}

func ToGetTarifResponseModel(src *dto.GetTarifResponseDto) *model.GetTarifResponseModel {
	if src == nil {
		return nil
	}

	items := make([]model.TarifItemModel, 0, len(src.Tarif))
	for _, item := range src.Tarif {
		items = append(items, ToTarifItemModel(item))
	}
	return &model.GetTarifResponseModel{
		Tarif: &items,
	}
}

func ToTarifItemDto(src model.TarifItemModel) dto.TarifItemDto {
	return dto.TarifItemDto{
		KetTarif: src.KetTarif,
		ID:       src.ID,
		Tarif:    src.Tarif,
	}
}

func ToTarifItemModel(src dto.TarifItemDto) model.TarifItemModel {
	return model.TarifItemModel{
		KetTarif: src.KetTarif,
		ID:       src.ID,
		Tarif:    src.Tarif,
	}
}
