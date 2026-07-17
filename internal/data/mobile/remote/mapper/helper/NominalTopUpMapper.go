package helper

import (
	dto "modulegue/internal/data/mobile/remote/dto/helper"
	model "modulegue/internal/domain/mobile/model/helper"
)

func ToTopupOptionsResponseDto(src *model.TopupOptionsResponseModel) *dto.TopupOptionsResponseDto {
	if src == nil {
		return nil
	}

	out := &dto.TopupOptionsResponseDto{}
	if src.Nominal == nil {
		empty := []dto.TopupOptionItemDto{}
		out.Nominal = &empty
		return out
	}

	items := make([]dto.TopupOptionItemDto, 0, len(*src.Nominal))
	for _, item := range *src.Nominal {
		items = append(items, ToTopupOptionItemDto(item))
	}
	out.Nominal = &items
	return out
}

func ToTopupOptionItemDto(src model.TopupOptionItemModel) dto.TopupOptionItemDto {
	return dto.TopupOptionItemDto{
		OptionID:      src.OptionID,
		NominalAmount: src.NominalAmount,
		Label:         src.Label,
	}
}

func ToTopupOptionsResponseModel(src *dto.TopupOptionsResponseDto) *model.TopupOptionsResponseModel {
	if src == nil {
		return nil
	}

	out := &model.TopupOptionsResponseModel{}
	if src.Nominal == nil {
		empty := []model.TopupOptionItemModel{}
		out.Nominal = &empty
		return out
	}

	items := make([]model.TopupOptionItemModel, 0, len(*src.Nominal))
	for _, item := range *src.Nominal {
		items = append(items, ToTopupOptionItemModel(item))
	}
	out.Nominal = &items
	return out
}

func ToTopupOptionItemModel(src dto.TopupOptionItemDto) model.TopupOptionItemModel {
	return model.TopupOptionItemModel{
		OptionID:      src.OptionID,
		NominalAmount: src.NominalAmount,
		Label:         src.Label,
	}
}
