package helper

import (
	dto "modulegue/internal/data/mobile/remote/dto/helper"
	model "modulegue/internal/domain/mobile/model/helper"
)

func ToNominalPaymentDto(src *model.NominalPaymentModel) *dto.NominalPaymentResponseDto {
	if src == nil {
		return nil
	}

	out := &dto.NominalPaymentResponseDto{}
	if src.Nominal == nil {
		empty := []dto.NominalItemDto{}
		out.Nominal = &empty
	} else {
		items := make([]dto.NominalItemDto, 0, len(*src.Nominal))
		for _, item := range *src.Nominal {
			items = append(items, ToNominalItemDto(item))
		}
		out.Nominal = &items
	}
	return out
}

func ToNominalPaymentModel(src *dto.NominalPaymentResponseDto) *model.NominalPaymentModel {
	if src == nil {
		return nil
	}

	out := &model.NominalPaymentModel{}
	if src.Nominal == nil {
		empty := []model.NominalItemModel{}
		out.Nominal = &empty
	} else {
		items := make([]model.NominalItemModel, 0, len(*src.Nominal))
		for _, item := range *src.Nominal {
			items = append(items, ToNominalItemModel(item))
		}
		out.Nominal = &items
	}
	return out
}

func ToNominalItemDto(src model.NominalItemModel) dto.NominalItemDto {
	return dto.NominalItemDto{
		OptionID:      src.OptionID,
		NominalAmount: src.NominalAmount,
		Label:         src.Label,
	}
}

func ToNominalItemModel(src dto.NominalItemDto) model.NominalItemModel {
	return model.NominalItemModel{
		OptionID:      src.OptionID,
		NominalAmount: src.NominalAmount,
		Label:         src.Label,
	}
}
