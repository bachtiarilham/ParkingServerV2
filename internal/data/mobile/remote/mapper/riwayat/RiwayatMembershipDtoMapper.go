package riwayat

import (
	dto "modulegue/internal/data/mobile/remote/dto/riwayat"
	model "modulegue/internal/domain/mobile/model/riwayat"
)

func ToRiwayatMembershipDto(src *model.RiwayatMembershipModel) *dto.RiwayatMembershipDto {
	if src == nil {
		return nil
	}

	items := make([]dto.RiwayatMembershipItemDto, len(src.Items))
	for i, item := range src.Items {
		items[i] = dto.RiwayatMembershipItemDto{
			ID:          item.ID,
			InvoiceNo:   item.InvoiceNo,
			PackageName: item.PackageName,
			PeriodStart: item.PeriodStart,
			PeriodEnd:   item.PeriodEnd,
			Amount:      item.Amount,
			PaidAt:      item.PaidAt,
			Status:      item.Status,
		}
	}

	return &dto.RiwayatMembershipDto{
		Summary: dto.RiwayatSummaryDto{
			PackageName:       src.Summary.PackageName,
			IsActive:          src.Summary.IsActive,
			ActiveUntil:       src.Summary.ActiveUntil,
			NextBillingAmount: src.Summary.NextBillingAmount,
			IsAutoRenew:       src.Summary.IsAutoRenew,
		},
		Items: items,
	}
}
