package topup

import (
	dto "modulegue/internal/data/website/remote/dto/topup"
	model "modulegue/internal/domain/web/model/topup"
)

func ToTopUpRequestDto(src *model.TopUpRequestModel) *dto.TopUpRequestDto {
	if src == nil {
		return nil
	}
	return &dto.TopUpRequestDto{
		IDUser:       src.IDUser,
		NominalTopUp: int(src.NominalTopUp),
	}
}

func ToTopUpRequestModel(src *dto.TopUpRequestDto) *model.TopUpRequestModel {
	if src == nil {
		return nil
	}
	return &model.TopUpRequestModel{
		IDUser:       src.IDUser,
		NominalTopUp: src.NominalTopUp,
	}
}

func ToTopUpResponseDto(src *model.TopUpResponseModel) *dto.TopUpResponseDto {
	if src == nil {
		return nil
	}
	return &dto.TopUpResponseDto{
		TopUpTransactionID: src.TopUpTransactionID,
		TopUpCode:          src.TopUpCode,
		ExternalReference:  src.ExternalReference,
		BalanceBefore:      src.BalanceBefore,
		BalanceAfter:       src.BalanceAfter,
	}
}

func ToTopUpResponseModel(src *dto.TopUpResponseDto) *model.TopUpResponseModel {
	if src == nil {
		return nil
	}
	return &model.TopUpResponseModel{
		TopUpTransactionID: src.TopUpTransactionID,
		TopUpCode:          src.TopUpCode,
		ExternalReference:  src.ExternalReference,
		BalanceBefore:      src.BalanceBefore,
		BalanceAfter:       src.BalanceAfter,
	}
}
