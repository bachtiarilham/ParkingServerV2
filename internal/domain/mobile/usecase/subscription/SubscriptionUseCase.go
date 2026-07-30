package subscription

import (
	"context"
	"fmt"
	"strconv"

	model "modulegue/internal/domain/mobile/model/subscription"
	"modulegue/internal/domain/mobile/repository"
)

type SubscriptionUseCase struct {
	subsRepo repository.SubscriptionRepository
}

func NewSubscriptionUseCase(
	subsRepo repository.SubscriptionRepository,
) *SubscriptionUseCase {
	return &SubscriptionUseCase{
		subsRepo: subsRepo,
	}
}

func (uc *SubscriptionUseCase) Execute(ctx context.Context, userId int64) (*model.SubscribeResponseModel, error) {
	result, err := uc.subsRepo.GetSubscribe(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("get subscription: %w", err)
	}

	if result == nil {
		return emptySubscriptionResponse(), nil
	}

	// 1. Format Statistik (Query B logic)
	if result.Statistik != nil {
		rawMinutes := result.Statistik.TotalJamParkirBulanLalu
		totalJam := rawMinutes / 60

		// Parse the temporarily stored raw values
		rawAmountPaid, _ := strconv.ParseInt(result.Statistik.TotalBiayaParkirBulanLaluText, 10, 64)
		rawAmountSaved, _ := strconv.ParseInt(result.Statistik.TotalPersentaseHematText, 10, 64)

		totalPaidAndSaved := rawAmountPaid + rawAmountSaved
		var pctHemat string
		if totalPaidAndSaved > 0 {
			percent := (float64(rawAmountSaved) / float64(totalPaidAndSaved)) * 100
			pctHemat = fmt.Sprintf("%.0f%%", percent)
		} else {
			pctHemat = "0%"
		}

		result.Statistik.TotalJamParkirBulanLalu = totalJam
		result.Statistik.TotalBiayaParkirBulanLaluText = formatRupiah(rawAmountPaid)
		result.Statistik.TotalPersentaseHematText = pctHemat
	}

	// 2. Format List Paket (Query C logic)
	for i := range result.ListPaket {
		pkg := &result.ListPaket[i]

		// Format Price to PriceLabel
		pkg.PriceLabel = formatRupiah(pkg.Price)

		// Map PeriodLabel
		periodCode := pkg.PeriodLabel
		switch periodCode {
		case "MONTHLY":
			pkg.PeriodLabel = "Bulan"
		case "YEARLY":
			pkg.PeriodLabel = "Tahun"
		}

		// Map BadgeLabel
		if pkg.BadgeLabel != nil {
			discPcnt, err := strconv.Atoi(*pkg.BadgeLabel)
			if err == nil && discPcnt > 0 {
				lbl := fmt.Sprintf("Diskon %d%%", discPcnt)
				pkg.BadgeLabel = &lbl
			} else {
				pkg.BadgeLabel = nil
			}
		}
	}

	if result.Benefits == nil {
		result.Benefits = []model.BenefitsModel{}
	}
	if result.ListPaket == nil {
		result.ListPaket = []model.DetailPaketModel{}
	}
	if result.Faq == nil {
		result.Faq = []model.FaqModel{}
	}

	return result, nil
}

func formatRupiah(amount int64) string {
	str := fmt.Sprintf("%d", amount)
	var result []rune
	count := 0
	for i := len(str) - 1; i >= 0; i-- {
		if count > 0 && count%3 == 0 {
			result = append([]rune{'.'}, result...)
		}
		result = append([]rune{rune(str[i])}, result...)
		count++
	}
	return "Rp " + string(result)
}

func emptySubscriptionResponse() *model.SubscribeResponseModel {
	return &model.SubscribeResponseModel{
		Benefits:  []model.BenefitsModel{},
		ListPaket: []model.DetailPaketModel{},
		Faq:       []model.FaqModel{},
	}
}
