package subscription

import (
	dto "modulegue/internal/data/mobile/remote/dto/subscription"
	model "modulegue/internal/domain/mobile/model/subscription"
)

func ToSubscribeDto(src *model.SubscribeResponseModel) *dto.SubscribeResponseDto {
	if src == nil {
		return nil
	}

	out := &dto.SubscribeResponseDto{
		Benefits:  ToBenefitsDtos(src.Benefits),
		ListPaket: ToDetailPaketDtos(src.ListPaket),
		Faq:       ToFaqDtos(src.Faq),
	}

	if src.ActivePaket != nil {
		out.ActivePaket = &dto.ActivePaketDto{
			ActivePackageName:    src.ActivePaket.ActivePackageName,
			ActivePackageExpired: src.ActivePaket.ActivePackageExpired,
		}
	}

	if src.Statistik != nil {
		out.Statistik = &dto.StatistikDto{
			TotalJamParkirBulanLalu:       src.Statistik.TotalJamParkirBulanLalu,
			TotalBiayaParkirBulanLaluText: src.Statistik.TotalBiayaParkirBulanLaluText,
			TotalPersentaseHematText:      src.Statistik.TotalPersentaseHematText,
		}
	}

	return out
}

func ToBenefitsDtos(src []model.BenefitsModel) []dto.BenefitsDto {
	if len(src) == 0 {
		return []dto.BenefitsDto{}
	}
	out := make([]dto.BenefitsDto, 0, len(src))
	for _, item := range src {
		out = append(out, dto.BenefitsDto{
			Name:        item.Name,
			Description: item.Description,
		})
	}
	return out
}

func ToDetailPaketDtos(src []model.DetailPaketModel) []dto.DetailPaketDto {
	if len(src) == 0 {
		return []dto.DetailPaketDto{}
	}
	out := make([]dto.DetailPaketDto, 0, len(src))
	for _, item := range src {
		out = append(out, dto.DetailPaketDto{
			Name:        item.Name,
			Price:       item.Price,
			PriceLabel:  item.PriceLabel,
			PeriodLabel: item.PeriodLabel,
			InfoLabel:   item.InfoLabel,
			BadgeLabel:  item.BadgeLabel,
			Benefits:    append([]string(nil), item.Benefits...),
		})
	}
	return out
}

func ToFaqDtos(src []model.FaqModel) []dto.FaqDto {
	if len(src) == 0 {
		return []dto.FaqDto{}
	}
	out := make([]dto.FaqDto, 0, len(src))
	for _, item := range src {
		out = append(out, dto.FaqDto{
			Question: item.Question,
			Answer:   item.Answer,
		})
	}
	return out
}
