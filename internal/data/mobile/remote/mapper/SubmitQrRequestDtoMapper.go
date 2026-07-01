package mapper

import (
	"modulegue/internal/data/mobile/remote/dto"
	"modulegue/internal/domain/mobile/model"
)

func ToSubmitQrRequestModel(src *dto.SubmitQrRequestDto) *model.SubmitQrRequestModel {
	if src == nil {
		return nil
	}
	return &model.SubmitQrRequestModel{
		SessionID:     src.SessionID,
		PlatNomor:     src.PlatNomor,
		Lokasi:        src.Lokasi,
		WaktuMasuk:    src.WaktuMasuk,
		Durasi:        src.Durasi,
		Nominal:       src.Nominal,
		IsPaid:        src.IsPaid,
		PaymentStatus: src.PaymentStatus,
		IsExpired:     src.IsExpired,
		StatusMessage: src.StatusMessage,
	}
}
