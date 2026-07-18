package helper

import (
	dto "modulegue/internal/data/mobile/remote/dto/helper"
	model "modulegue/internal/domain/mobile/model/helper"
)

func ToTopupResponseDto(src *model.TopupResponseModel) *dto.TopupResponseDto {
	if src == nil {
		return nil
	}

	out := &dto.TopupResponseDto{}
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

	if src.MetodeBayar == nil {
		empty := []dto.MetodeItemDto{}
		out.MetodeBayar = &empty
		return out
	}

	methods := make([]dto.MetodeItemDto, 0, len(*src.MetodeBayar))
	for _, item := range *src.MetodeBayar {
		methods = append(methods, ToMetodeItemDto(item))
	}
	out.MetodeBayar = &methods
	return out
}

func ToTopupResponseModel(src *dto.TopupResponseDto) *model.TopupResponseModel {
	if src == nil {
		return nil
	}

	out := &model.TopupResponseModel{}
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

	if src.MetodeBayar == nil {
		empty := []model.MetodeItemModel{}
		out.MetodeBayar = &empty
		return out
	}

	methods := make([]model.MetodeItemModel, 0, len(*src.MetodeBayar))
	for _, item := range *src.MetodeBayar {
		methods = append(methods, ToMetodeItemModel(item))
	}
	out.MetodeBayar = &methods
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

func ToMetodeItemDto(src model.MetodeItemModel) dto.MetodeItemDto {
	return dto.MetodeItemDto{
		PaymentMethodId: src.PaymentMethodId,
		NamaPayment:     src.NamaPayment,
		CodePayment:     src.CodePayment,
		LogoPayment:     src.LogoPayment,
	}
}

func ToMetodeItemModel(src dto.MetodeItemDto) model.MetodeItemModel {
	return model.MetodeItemModel{
		PaymentMethodId: src.PaymentMethodId,
		NamaPayment:     src.NamaPayment,
		CodePayment:     src.CodePayment,
		LogoPayment:     src.LogoPayment,
	}
}
