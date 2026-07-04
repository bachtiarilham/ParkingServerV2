package payment

import (
	dto "modulegue/internal/data/mobile/remote/dto/payment"
	model "modulegue/internal/domain/mobile/model/payment"
)

func ToPostPaymentParkingRequestModel(src *dto.PostPaymentParkingRequestDto) *model.PostPaymentParkingRequestModel {
	if src == nil {
		return nil
	}
	return &model.PostPaymentParkingRequestModel{
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
