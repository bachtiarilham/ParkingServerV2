package location

// import (
// 	"context"
// 	"modulegue/internal/domain/web/metrics"
// 	"time"
// )

// type Repository interface {
// 	// GetLocations mengambil semua lokasi dengan data agregat
// 	GetLocations(ctx context.Context) ([]*LocationAggregate, error)

// 	// GetLocationByID mengambil satu lokasi berdasarkan ID
// 	GetLocationByID(ctx context.Context, id int64) (*LocationAggregate, error)

// 	// GetHourlyTraffic mengambil data traffic per jam untuk rentang waktu tertentu
// 	GetHourlyTraffic(ctx context.Context, startDate, endDate time.Time) ([]*metrics.HourlyTrafficPoint, error)

// 	// GetHeatmapData mengambil data heatmap untuk rentang waktu tertentu
// 	GetHeatmapData(ctx context.Context, startDate, endDate time.Time) ([]*metrics.HeatmapPoint, error)

// 	// UpdateLocationSettings memperbarui setting lokasi (tarif, catatan)
// 	UpdateLocationSettings(ctx context.Context, id int64, tariffMotor, tariffMobil int64, operationalNote string) error

// 	// SaveShiftTemplates menyimpan template shift untuk lokasi
// 	SaveShiftTemplates(ctx context.Context, locationID int64, templates []ParkingShiftTemplate) error

// 	// Fungsi tambahan lainnya sesuai kebutuhan
// }
