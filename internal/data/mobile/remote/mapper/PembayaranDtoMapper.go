package mapper

import (
	"modulegue/internal/data/mobile/remote/dto"
	"modulegue/internal/domain/mobile/model"
)

func ToPembayaranDto(src *model.PembayaranModel) *dto.PembayaranDto {
	if src == nil {
		return nil
	}
	out := &dto.PembayaranDto{
		Title:               src.Title,
		StatusCard:          ToPembayaranStatusCardDto(src.StatusCard),
		TotalPembayaran:     src.TotalPembayaran,
		DetailLabel:         src.DetailLabel,
		QrisSection:         ToPembayaranQrisSectionDto(src.QrisSection),
		PaymentOptionsTitle: src.PaymentOptionsTitle,
		PaymentOptions:      make([]dto.PembayaranOptionDto, 0, len(src.PaymentOptions)),
		PrintButtonLabel:    src.PrintButtonLabel,
	}
	for _, item := range src.PaymentOptions {
		out.PaymentOptions = append(out.PaymentOptions, *ToPembayaranOptionDto(&item))
	}
	return out
}

func ToPembayaranStatusCardDto(src *model.PembayaranStatusCardModel) *dto.PembayaranStatusCardDto {
	if src == nil {
		return nil
	}
	return &dto.PembayaranStatusCardDto{Title: src.Title, Message: src.Message, IsSuccess: src.IsSuccess}
}

func ToPembayaranQrisSectionDto(src *model.PembayaranQrisSectionModel) *dto.PembayaranQrisSectionDto {
	if src == nil {
		return nil
	}
	return &dto.PembayaranQrisSectionDto{Title: src.Title, QrContent: ToIsiQrDto(src.QrContent), Information: src.Information, Countdown: src.Countdown, AlternativeLabel: src.AlternativeLabel}
}

func ToIsiQrDto(src *model.IsiQrModel) *dto.IsiQrDto {
	if src == nil {
		return nil
	}
	return &dto.IsiQrDto{SessionID: src.SessionID, PlatNomor: src.PlatNomor, Lokasi: src.Lokasi, WaktuMasuk: src.WaktuMasuk, Durasi: src.Durasi, Nominal: src.Nominal, IsPaid: src.IsPaid, PaymentStatus: src.PaymentStatus, IsExpired: src.IsExpired, StatusMessage: src.StatusMessage}
}

func ToPembayaranOptionDto(src *model.PembayaranOptionModel) *dto.PembayaranOptionDto {
	if src == nil {
		return nil
	}
	return &dto.PembayaranOptionDto{Type: src.Type, Title: src.Title, Subtitle: src.Subtitle}
}
