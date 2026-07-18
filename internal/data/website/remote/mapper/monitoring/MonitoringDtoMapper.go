package monitoring

import (
	dto "modulegue/internal/data/website/remote/dto/monitoring"
	model "modulegue/internal/domain/web/model/monitoring"
)

func ToMonitoringRequestDto(src *model.MonitoringRequestModel) *dto.MonitoringRequestDto {
	if src == nil {
		return nil
	}
	return &dto.MonitoringRequestDto{
		TglAwal:     src.TglAwal,
		TglAkhir:    src.TglAkhir,
		IDLokasi:    src.IDLokasi,
		NamaPetugas: src.NamaPetugas,
	}
}

func ToMonitoringRequestModel(src *dto.MonitoringRequestDto) *model.MonitoringRequestModel {
	if src == nil {
		return nil
	}
	return &model.MonitoringRequestModel{
		TglAwal:     src.TglAwal,
		TglAkhir:    src.TglAkhir,
		IDLokasi:    src.IDLokasi,
		NamaPetugas: src.NamaPetugas,
	}
}

func ToMonitoringResponseDto(src *model.MonitoringResponseModel) *dto.MonitoringResponseDto {
	if src == nil {
		return nil
	}

	out := &dto.MonitoringResponseDto{}
	if src.Parlok != nil {
		items := make([]dto.ParlokItemDto, 0, len(*src.Parlok))
		for _, item := range *src.Parlok {
			items = append(items, dto.ParlokItemDto{
				NamaParlok:      item.NamaParlok,
				IDZona:          item.IDZona,
				NamaZona:        item.NamaZona,
				LatMin:          item.LatMin,
				LatMax:          item.LatMax,
				LngMin:          item.LngMin,
				LngMax:          item.LngMax,
				CenterX:         item.CenterX,
				CenterY:         item.CenterY,
				Altitude:        item.Altitude,
				PendapatanMotor: item.PendapatanMotor,
				PendapatanMobil: item.PendapatanMobil,
				TotalPendapatan: item.TotalPendapatan,
			})
		}
		out.Parlok = &items
	}
	if src.Transaksi != nil {
		items := make([]dto.TransaksiItemDto, 0, len(*src.Transaksi))
		for _, item := range *src.Transaksi {
			items = append(items, dto.TransaksiItemDto{
				NamaJukir:  item.NamaJukir,
				Parlok:     item.Parlok,
				Zona:       item.Zona,
				Plat:       item.Plat,
				Waktu:      item.Waktu,
				Kendaraan:  item.Kendaraan,
				Pembayaran: item.Pembayaran,
				Tarif:      item.Tarif,
			})
		}
		out.Transaksi = &items
	}
	return out
}

func ToMonitoringResponseModel(src *dto.MonitoringResponseDto) *model.MonitoringResponseModel {
	if src == nil {
		return nil
	}

	out := &model.MonitoringResponseModel{}
	if src.Parlok != nil {
		items := make([]model.ParlokItemModel, 0, len(*src.Parlok))
		for _, item := range *src.Parlok {
			items = append(items, model.ParlokItemModel{
				NamaParlok:      item.NamaParlok,
				IDZona:          item.IDZona,
				NamaZona:        item.NamaZona,
				LatMin:          item.LatMin,
				LatMax:          item.LatMax,
				LngMin:          item.LngMin,
				LngMax:          item.LngMax,
				CenterX:         item.CenterX,
				CenterY:         item.CenterY,
				Altitude:        item.Altitude,
				PendapatanMotor: item.PendapatanMotor,
				PendapatanMobil: item.PendapatanMobil,
				TotalPendapatan: item.TotalPendapatan,
			})
		}
		out.Parlok = &items
	}
	if src.Transaksi != nil {
		items := make([]model.TransaksiItemModel, 0, len(*src.Transaksi))
		for _, item := range *src.Transaksi {
			items = append(items, model.TransaksiItemModel{
				NamaJukir:  item.NamaJukir,
				Parlok:     item.Parlok,
				Zona:       item.Zona,
				Plat:       item.Plat,
				Waktu:      item.Waktu,
				Kendaraan:  item.Kendaraan,
				Pembayaran: item.Pembayaran,
				Tarif:      item.Tarif,
			})
		}
		out.Transaksi = &items
	}
	return out
}
