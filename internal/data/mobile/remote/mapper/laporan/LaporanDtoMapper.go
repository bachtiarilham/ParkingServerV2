package laporan

import (
	dto "modulegue/internal/data/mobile/remote/dto/laporan"
	model "modulegue/internal/domain/mobile/model/laporan"
)

func ToLaporanDto(src *model.LaporanModel) *dto.LaporanResponseDto {
	if src == nil {
		return nil
	}

	out := &dto.LaporanResponseDto{
		TanggalTerpilih:    src.TanggalTerpilih,
		Periode:            ToLaporanDateRangeDto(src.Periode),
		Summary:            ToLaporanSummaryDto(src.Summary),
		ChartBars:          make([]dto.LaporanChartBarDto, 0, len(src.ChartBars)),
		PaymentSummaries:   make([]dto.LaporanPaymentSummaryDto, 0, len(src.PaymentSummaries)),
		RecentTransactions: make([]dto.LaporanRecentTransactionDto, 0, len(src.RecentTransactions)),
	}

	for _, item := range src.ChartBars {
		out.ChartBars = append(out.ChartBars, *ToLaporanChartBarDto(&item))
	}
	for _, item := range src.PaymentSummaries {
		out.PaymentSummaries = append(out.PaymentSummaries, *ToLaporanPaymentSummaryDto(&item))
	}
	for _, item := range src.RecentTransactions {
		out.RecentTransactions = append(out.RecentTransactions, *ToLaporanRecentTransactionDto(&item))
	}

	return out
}

func ToLaporanDateRangeDto(src *model.LaporanDateRangeModel) *dto.LaporanDateRangeDto {
	if src == nil {
		return nil
	}
	return &dto.LaporanDateRangeDto{StartDate: src.StartDate, EndDate: src.EndDate, Label: src.Label}
}

func ToLaporanSummaryDto(src *model.LaporanSummaryModel) *dto.LaporanSummaryDto {
	if src == nil {
		return nil
	}
	return &dto.LaporanSummaryDto{TotalTransaksi: src.TotalTransaksi, TotalPendapatan: src.TotalPendapatan, RataRataTransaksi: src.RataRataTransaksi}
}

func ToLaporanChartBarDto(src *model.LaporanChartBarModel) *dto.LaporanChartBarDto {
	if src == nil {
		return nil
	}
	return &dto.LaporanChartBarDto{Tanggal: src.Tanggal, Amount: src.Amount, Value: src.Value, PeriodLabel: src.PeriodLabel, PeriodStartDate: src.PeriodStartDate, PeriodEndDate: src.PeriodEndDate}
}

func ToLaporanPaymentSummaryDto(src *model.LaporanPaymentSummaryModel) *dto.LaporanPaymentSummaryDto {
	if src == nil {
		return nil
	}
	return &dto.LaporanPaymentSummaryDto{Label: src.Label, Amount: src.Amount, Percentage: src.Percentage}
}

func ToLaporanRecentTransactionDto(src *model.LaporanRecentTransactionModel) *dto.LaporanRecentTransactionDto {
	if src == nil {
		return nil
	}
	return &dto.LaporanRecentTransactionDto{Code: src.Code, Time: src.Time, Total: src.Total, PaymentTag: src.PaymentTag}
}
