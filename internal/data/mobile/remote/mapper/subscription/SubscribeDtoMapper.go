package subscription

import (
	dto "modulegue/internal/data/mobile/remote/dto/subscription"
	model "modulegue/internal/domain/mobile/model/subscription"
)

func ToSubscribeDto(src *model.SubscriptionResponseModel) *dto.SubscriptionResponseDto {
	if src == nil {
		return nil
	}
	out := &dto.SubscriptionResponseDto{
		ActivePackageName:    src.ActivePackageName,
		ActivePackageExpired: src.ActivePackageExpired,
		ActivePackageBenefit: append([]string(nil), src.ActivePackageBenefit...),
		ListPaket:            ToListPaketDto(src.ListPaket),
		PromoTersedia:        ToPromoTersediaDto(src.PromoTersedia),
	}
	return out
}

func ToListPaketDto(src model.ListPaket) dto.ListPaket {
	return dto.ListPaket{
		Bulanan:   ToDetailPaketDtos(src.Bulanan),
		EnamBulan: ToDetailPaketDtos(src.EnamBulan),
		Tahunan:   ToDetailPaketDtos(src.Tahunan),
	}
}

func ToDetailPaketDtos(src []model.DetailPaket) []dto.DetailPaket {
	if len(src) == 0 {
		return []dto.DetailPaket{}
	}
	out := make([]dto.DetailPaket, 0, len(src))
	for _, item := range src {
		out = append(out, dto.DetailPaket{
			NamaPaket:      item.NamaPaket,
			Harga:          item.Harga,
			CoverageLokasi: append([]string(nil), item.CoverageLokasi...),
			BenefitPackage: append([]string(nil), item.BenefitPackage...),
		})
	}
	return out
}

func ToPromoTersediaDto(src model.PromoTersedia) dto.PromoTersedia {
	return dto.PromoTersedia{
		SyaratDanKetentuan: append([]string(nil), src.SyaratDanKetentuan...),
		EachPromo:          ToDetailPromoDtos(src.EachPromo),
	}
}

func ToDetailPromoDtos(src []model.DetailPromo) []dto.DetailPromo {
	if len(src) == 0 {
		return []dto.DetailPromo{}
	}
	out := make([]dto.DetailPromo, 0, len(src))
	for _, item := range src {
		out = append(out, dto.DetailPromo{
			NamaPromo:   item.NamaPromo,
			BesarDiskon: item.BesarDiskon,
		})
	}
	return out
}
