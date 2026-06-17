package metrics

// import (
// 	"context"
// 	"fmt"
// 	"modulegue/internal/domain/web/location"
// 	"modulegue/internal/domain/web/metrics" // Asumsikan ada domain metrics entity.go
// 	"sort"
// 	"time"
// )

// type GetDashboardMetricsOutput struct {
// 	LocationMetrics     []metrics.LocationMetric     `json:"locationMetrics"`
// 	ComparisonMetrics   []metrics.ComparisonMetric   `json:"comparisonMetrics"`
// 	HourlyTrafficPoints []metrics.HourlyTrafficPoint `json:"hourlyTrafficPoints"` // Jika di sini
// 	ParkingHeatmap      []metrics.HeatmapPoint       `json:"parkingHeatmap"`      // Jika di sini
// 	// Tambahkan field lain sesuai kebutuhan dashboard
// }

// type GetDashboardMetricsUseCase struct {
// 	locationRepo location.Repository
// 	// metricsRepo metrics.Repository // Jika ada repo khusus metrics
// }

// func NewGetDashboardMetricsUseCase(locationRepo location.Repository /*, metricsRepo metrics.Repository*/) *GetDashboardMetricsUseCase {
// 	return &GetDashboardMetricsUseCase{
// 		locationRepo: locationRepo,
// 		// metricsRepo: metricsRepo,
// 	}
// }

// func (uc *GetDashboardMetricsUseCase) Execute(ctx context.Context, now time.Time) (GetDashboardMetricsOutput, error) {
// 	// 1. Ambil data dasar (lokasi, traffic, heatmap)
// 	locations, err := uc.locationRepo.GetLocations(ctx)
// 	if err != nil {
// 		return GetDashboardMetricsOutput{}, fmt.Errorf("get locations: %w", err)
// 	}

// 	todayStart, todayEnd := todayWindow(now)             // Helper untuk mendapatkan rentang hari ini
// 	yesterdayStart, yesterdayEnd := yesterdayWindow(now) // Helper untuk kemarin

// 	hourlyTraffic, err := uc.locationRepo.GetHourlyTraffic(ctx, todayStart, todayEnd)
// 	if err != nil {
// 		// Log error, gunakan slice kosong
// 		hourlyTraffic = []*metrics.HourlyTrafficPoint{}
// 	}

// 	heatmap, err := uc.locationRepo.GetHeatmapData(ctx, now.AddDate(0, 0, -7), now) // Ambil 7 hari terakhir
// 	if err != nil {
// 		// Log error, gunakan slice kosong
// 		heatmap = []*metrics.HeatmapPoint{}
// 	}

// 	// 2. Hitung Metrik Lokasi (Revenue, Occupancy)
// 	revenueMetrics, occupancyMetrics := buildLocationMetrics(locations)

// 	// 3. Hitung Metrik Perbandingan (Hari Ini vs Kemarin)
// 	comparisonMetrics, err := uc.buildComparisonMetrics(ctx, todayStart, todayEnd, yesterdayStart, yesterdayEnd)
// 	if err != nil {
// 		// Log error, gunakan slice kosong
// 		comparisonMetrics = []metrics.ComparisonMetric{}
// 	}

// 	return GetDashboardMetricsOutput{
// 		LocationMetrics:     revenueMetrics, // Atau gabungkan revenue dan occupancy jika frontend butuh
// 		ComparisonMetrics:   comparisonMetrics,
// 		HourlyTrafficPoints: hourlyTraffic,
// 		ParkingHeatmap:      heatmap,
// 		// Tambahkan field lain
// 	}, nil
// }

// // Helper untuk menghitung revenue dan occupancy metrics
// func buildLocationMetrics(locations []*location.LocationAggregate) ([]metrics.LocationMetric, []metrics.LocationMetric) {
// 	revenue := make([]metrics.LocationMetric, 0, len(locations))
// 	occupancy := make([]metrics.LocationMetric, 0, len(locations))

// 	// Hitung max traffic untuk persentase occupancy
// 	maxTraffic := int64(1)
// 	for _, loc := range locations {
// 		totalVehicles := loc.Cars + loc.Motorcycles
// 		if totalVehicles > maxTraffic {
// 			maxTraffic = totalVehicles
// 		}
// 	}

// 	for _, loc := range locations {
// 		// Revenue Metric
// 		totalRevenue := (loc.TariffMotor * loc.Motorcycles) + (loc.TariffMobil * loc.Cars)
// 		revenue = append(revenue, metrics.LocationMetric{
// 			Name:      loc.Name,
// 			Value:     totalRevenue,
// 			Secondary: "Akumulasi estimasi per lokasi",
// 			Tone:      "blue", // Default tone, bisa diupdate berdasarkan ranking
// 		})

// 		// Occupancy Metric
// 		totalVehicles := loc.Cars + loc.Motorcycles
// 		occupancyPercent := (totalVehicles * 100) / maxTraffic
// 		occupancyTone := "green"
// 		if occupancyPercent >= 80 {
// 			occupancyTone = "orange"
// 		} else if occupancyPercent >= 55 {
// 			occupancyTone = "blue"
// 		} else if occupancyPercent >= 30 {
// 			occupancyTone = "gold"
// 		}
// 		occupancy = append(occupancy, metrics.LocationMetric{
// 			Name:      loc.Name,
// 			Value:     occupancyPercent,
// 			Secondary: loc.OccupancyLabel,
// 			Tone:      occupancyTone,
// 		})
// 	}

// 	// Urutkan berdasarkan nilai (descending)
// 	sort.Slice(revenue, func(i, j int) bool { return revenue[i].Value > revenue[j].Value })
// 	sort.Slice(occupancy, func(i, j int) bool { return occupancy[i].Value > occupancy[j].Value })

// 	// Beri tone berdasarkan ranking (opsional, sesuai kebijakan UI)
// 	for i := range revenue {
// 		revenue[i].Tone = toneFromRank(i, len(revenue)) // Helper untuk warna berdasarkan rank
// 	}
// 	for i := range occupancy {
// 		occupancy[i].Tone = toneFromRank(i, len(occupancy)) // Helper untuk warna berdasarkan rank
// 	}

// 	return revenue, occupancy
// }

// // Helper untuk rentang waktu
// func todayWindow(now time.Time) (time.Time, time.Time) {
// 	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
// 	end := start.Add(24 * time.Hour)
// 	return start, end
// }

// func yesterdayWindow(now time.Time) (time.Time, time.Time) {
// 	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -1)
// 	end := start.Add(24 * time.Hour)
// 	return start, end
// }

// // Helper untuk warna berdasarkan rank
// func toneFromRank(rank, total int) string {
// 	if total <= 0 {
// 		return "gray"
// 	}
// 	percentage := float64(rank+1) / float64(total)

// 	if percentage <= 0.33 {
// 		return "green" // Top 33%
// 	} else if percentage <= 0.66 {
// 		return "blue" // Middle 33%
// 	} else {
// 		return "orange" // Bottom 33%
// 	}
// }
