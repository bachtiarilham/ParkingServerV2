package subscription

import (
	dto "modulegue/internal/data/mobile/remote/dto/subscription"
	model "modulegue/internal/domain/mobile/model/subscription"
)

func ToSubscribeDto(src *model.SubscribeModel) *dto.SubscriptionResponseDto {
	if src == nil {
		return nil
	}
	out := &dto.SubscriptionResponseDto{
		StatusCard:  ToStatusCardDto(src.StatusCard),
		PackageCard: make([]dto.PackageCardDto, 0, len(src.PackageCard)),
		Promo:       make([]dto.PromoDto, 0, len(src.Promo)),
	}
	for _, item := range src.PackageCard {
		out.PackageCard = append(out.PackageCard, *ToPackageCardDto(&item))
	}
	for _, item := range src.Promo {
		out.Promo = append(out.Promo, *ToPromoDto(&item))
	}
	return out
}

func ToStatusCardDto(src *model.StatusCardModel) *dto.StatusCardDto {
	if src == nil {
		return nil
	}
	return &dto.StatusCardDto{PaketAktif: src.PaketAktif, Kadaluarsa: src.Kadaluarsa, Benefit: src.Benefit}
}

func ToPackageCardDto(src *model.PackageCardModel) *dto.PackageCardDto {
	if src == nil {
		return nil
	}
	return &dto.PackageCardDto{NamaPaket: src.NamaPaket, Harga: src.Harga, MasaBerlaku: src.MasaBerlaku, JumlahDiskon: src.JumlahDiskon, Deskripsi: src.Deskripsi, Benefit: append([]string(nil), src.Benefit...)}
}

func ToPromoDto(src *model.PromoModel) *dto.PromoDto {
	if src == nil {
		return nil
	}
	out := &dto.PromoDto{SNk: append([]string(nil), src.SNk...), Promo: make([]dto.PromoTerpilihDto, 0, len(src.Promo))}
	for _, item := range src.Promo {
		out.Promo = append(out.Promo, *ToPromoTerpilihDto(&item))
	}
	return out
}

func ToPromoTerpilihDto(src *model.PromoTerpilihModel) *dto.PromoTerpilihDto {
	if src == nil {
		return nil
	}
	return &dto.PromoTerpilihDto{NamaPromo: src.NamaPromo, Deskripsi: src.Deskripsi, JumlahDiskon: src.JumlahDiskon}
}
