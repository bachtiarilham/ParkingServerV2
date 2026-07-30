package helper

import (
	dto "modulegue/internal/data/mobile/remote/dto/helper"
	model "modulegue/internal/domain/mobile/model/helper"
)

func ToPaymentMethodDto(src *model.PaymentMethodModel) *dto.PaymentMethodResponseDto {
	if src == nil {
		return nil
	}

	out := &dto.PaymentMethodResponseDto{}
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

func ToPaymentMethodModel(src *dto.PaymentMethodResponseDto) *model.PaymentMethodModel {
	if src == nil {
		return nil
	}

	out := &model.PaymentMethodModel{}

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
