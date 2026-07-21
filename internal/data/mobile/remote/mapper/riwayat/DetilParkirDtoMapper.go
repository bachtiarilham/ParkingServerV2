package riwayat

import (
	dto "modulegue/internal/data/mobile/remote/dto/riwayat"
	model "modulegue/internal/domain/mobile/model/riwayat"
)

func ToDetilParkirRequestModel(src *dto.DetilParkirRequestDto) *model.DetilParkirRequestModel {
	if src == nil {
		return nil
	}
	return &model.DetilParkirRequestModel{
		TransactionCode: src.TransactionCode,
	}
}

func ToDetilParkirDto(src *model.DetilParkirModel) *dto.DetilParkirDto {
	if src == nil {
		return nil
	}
	return &dto.DetilParkirDto{
		Tanggal:           src.Tanggal,
		TransactionID:     src.TransactionID,
		TransactionCode:   src.TransactionCode,
		SessionID:         src.SessionID,
		PlateNumber:       src.PlateNumber,
		VehicleTypeID:     src.VehicleTypeID,
		VehicleTypeCode:   src.VehicleTypeCode,
		VehicleTypeName:   src.VehicleTypeName,
		PaymentMethodID:   src.PaymentMethodID,
		PaymentMethodCode: src.PaymentMethodCode,
		PaymentMethodName: src.PaymentMethodName,
		LocationID:        src.LocationID,
		LocationName:      src.LocationName,
		LocationAddress:   src.LocationAddress,
		AreaID:            src.AreaID,
		AreaName:          src.AreaName,
		ZoneID:            src.ZoneID,
		ZoneName:          src.ZoneName,
		BaseAmount:        src.BaseAmount,
		DiscountAmount:    src.DiscountAmount,
		FinalAmount:       src.FinalAmount,
		CompanyShare:      src.CompanyShare,
		JukirShare:        src.JukirShare,
		TaxAmount:         src.TaxAmount,
		FeeAmount:         src.FeeAmount,
		TransactionStatus: src.TransactionStatus,
		OperationType:     src.OperationType,
		OccurredAt:        src.OccurredAt,
		PaidAt:            src.PaidAt,
		CreatedAt:         src.CreatedAt,
	}
}

func ToDetilTransaksiRequestModel(src *dto.DetilTransaksiRequestDto) *model.DetilTransaksiRequestModel {
	if src == nil {
		return nil
	}
	return &model.DetilTransaksiRequestModel{
		TopUpCode: src.TopUpCode,
	}
}

func ToDetilTransaksiDto(src *model.DetilTransaksiModel) *dto.DetilTransaksiDto {
	if src == nil {
		return nil
	}
	return &dto.DetilTransaksiDto{
		Tanggal:            src.Tanggal,
		TopUpTransactionID: src.TopUpTransactionID,
		TopUpCode:          src.TopUpCode,
		UserID:             src.UserID,
		WalletID:           src.WalletID,
		PaymentMethodID:    src.PaymentMethodID,
		PaymentMethodCode:  src.PaymentMethodCode,
		PaymentMethodName:  src.PaymentMethodName,
		Amount:             src.Amount,
		AdminFee:           src.AdminFee,
		TotalAmount:        src.TotalAmount,
		TransactionStatus:  src.TransactionStatus,
		ExternalReference:  src.ExternalReference,
		ProviderName:       src.ProviderName,
		CreatedAt:          src.CreatedAt,
		ExpiredAt:          src.ExpiredAt,
		PaidAt:             src.PaidAt,
		CompletedAt:        src.CompletedAt,
		FailedReason:       src.FailedReason,
	}
}
