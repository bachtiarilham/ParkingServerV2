package repository

// import (
// 	"context"
// 	"fmt"
// 	metrics "modulegue/internal/domain/web/metrics"
// 	metricsuc "modulegue/internal/usecase/web/metrics"
// 	"time"
// )

// func (uc *metricsuc.GetDashboardMetricsUseCase) buildComparisonMetrics(ctx context.Context, todayStart, todayEnd, yesterdayStart, yesterdayEnd time.Time) ([]location.ComparisonMetric, error) {
// 	// Query dari repository untuk mendapatkan data hari ini dan kemarin
// 	var revenueToday, revenueYesterday, txToday, txYesterday int64

// 	// Query untuk pendapatan
// 	err := uc.locationRepo.db.QueryRowContext(ctx, `
// 		SELECT COALESCE(SUM(final_amount), 0)
// 		FROM financial_parking_transaction
// 		WHERE paid_at >= ? AND paid_at < ?
// 	`, todayStart, todayEnd).Scan(&revenueToday)
// 	if err != nil {
// 		return nil, fmt.Errorf("query revenue today: %w", err)
// 	}

// 	err = uc.locationRepo.db.QueryRowContext(ctx, `
// 		SELECT COALESCE(SUM(final_amount), 0)
// 		FROM financial_parking_transaction
// 		WHERE paid_at >= ? AND paid_at < ?
// 	`, yesterdayStart, yesterdayEnd).Scan(&revenueYesterday)
// 	if err != nil {
// 		return nil, fmt.Errorf("query revenue yesterday: %w", err)
// 	}

// 	// Query untuk jumlah transaksi kendaraan masuk
// 	err = uc.locationRepo.db.QueryRowContext(ctx, `
// 		SELECT COUNT(*)
// 		FROM parking_session
// 		WHERE started_at >= ? AND started_at < ?
// 	`, todayStart, todayEnd).Scan(&txToday)
// 	if err != nil {
// 		return nil, fmt.Errorf("query transactions today: %w", err)
// 	}

// 	err = uc.locationRepo.db.QueryRowContext(ctx, `
// 		SELECT COUNT(*)
// 		FROM parking_session
// 		WHERE started_at >= ? AND started_at < ?
// 	`, yesterdayStart, yesterdayEnd).Scan(&txYesterday)
// 	if err != nil {
// 		return nil, fmt.Errorf("query transactions yesterday: %w", err)
// 	}

// 	// Query untuk occupancy puncak (contoh: rata-rata aktif per jam)
// 	var activeToday, activeYesterday int64
// 	err = uc.locationRepo.db.QueryRowContext(ctx, `
// 		SELECT AVG(active_count) * 24 FROM (
// 			SELECT HOUR(started_at) as hour, COUNT(*) as active_count
// 			FROM parking_session
// 			WHERE started_at >= ? AND started_at < ?
// 			GROUP BY HOUR(started_at)
// 		) as hourly_avg
// 	`, todayStart, todayEnd).Scan(&activeToday)
// 	if err != nil {
// 		// Jika query gagal, abaikan atau gunakan 0
// 		activeToday = 0
// 	}

// 	err = uc.locationRepo.db.QueryRowContext(ctx, `
// 		SELECT AVG(active_count) * 24 FROM (
// 			SELECT HOUR(started_at) as hour, COUNT(*) as active_count
// 			FROM parking_session
// 			WHERE started_at >= ? AND started_at < ?
// 			GROUP BY HOUR(started_at)
// 		) as hourly_avg
// 	`, yesterdayStart, yesterdayEnd).Scan(&activeYesterday)
// 	if err != nil {
// 		activeYesterday = 0
// 	}

// 	metrics := []metrics.ComparisonMetric{
// 		{Label: "Pendapatan", Today: revenueToday, Yesterday: revenueYesterday, Unit: "currency"},
// 		{Label: "Kendaraan Masuk", Today: txToday, Yesterday: txYesterday, Unit: "count"},
// 		{Label: "Occupancy Rata-rata", Today: activeToday, Yesterday: activeYesterday, Unit: "count"}, // Ubah label/unit jika perlu
// 	}

// 	return metrics, nil
// }
