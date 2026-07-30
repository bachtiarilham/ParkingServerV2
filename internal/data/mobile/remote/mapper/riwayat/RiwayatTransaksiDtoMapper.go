package riwayat

import (
	dto "modulegue/internal/data/mobile/remote/dto/riwayat"
	model "modulegue/internal/domain/mobile/model/riwayat"
)

func ToRiwayatTransaksiDto(src *model.RiwayatTransaksiModel) *dto.RiwayatTransaksiDto {
	if src == nil {
		return nil
	}

	items := make([]dto.RiwayatTransaksiItemDto, len(src.Items))
	for i, item := range src.Items {
		items[i] = dto.RiwayatTransaksiItemDto{
			ID:              item.ID,
			ReferenceNo:     item.ReferenceNo,
			Title:           item.Title,
			TransactionType: item.TransactionType,
			Flow:            item.Flow,
			Amount:          item.Amount,
			Status:          item.Status,
			CreatedAt:       item.CreatedAt,
		}
	}

	return &dto.RiwayatTransaksiDto{
		Summary: dto.WalletSummaryDto{
			TotalIncome:    src.Summary.TotalIncome,
			TotalExpense:   src.Summary.TotalExpense,
			CurrentBalance: src.Summary.CurrentBalance,
		},
		Items: items,
	}
}
