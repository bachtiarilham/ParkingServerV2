package location

// import (
// 	"context"
// 	"fmt"
// 	"modulegue/internal/domain/web/location"
// )

// type GetLocationsOutput struct {
// 	Locations []*location.ParkingLocation `json:"locations"`
// 	// Tambahkan field lain jika perlu, misalnya metadata pagination
// }

// type GetLocationsUseCase struct {
// 	repo location.Repository
// }

// func NewGetLocationsUseCase(repo location.Repository) *GetLocationsUseCase {
// 	return &GetLocationsUseCase{repo: repo}
// }

// func (uc *GetLocationsUseCase) Execute(ctx context.Context) (GetLocationsOutput, error) {
// 	locations, err := uc.repo.GetLocations(ctx)
// 	if err != nil {
// 		return GetLocationsOutput{}, fmt.Errorf("get all locations from repo: %w", err)
// 	}

// 	// --- (Opsional) Post-processing di Use Case ---
// 	// Misalnya, hitung occupancy label di sini jika memerlukan perbandingan antar lokasi
// 	// atau jika logika occupancy kompleks.
// 	// Kita gunakan helper untuk menghitung maxTraffic global dulu
// 	maxTransactions := int64(0)
// 	for _, loc := range locations {
// 		totalVehicles := loc.Cars + loc.Motorcycles
// 		if totalVehicles > maxTransactions {
// 			maxTransactions = totalVehicles
// 		}
// 	}
// 	if maxTransactions == 0 {
// 		maxTransactions = 1
// 	} // Hindari pembagian nol

// 	// Update occupancy label berdasarkan max global
// 	for _, loc := range locations {
// 		totalVehicles := loc.Cars + loc.Motorcycles
// 		occupancyPercent := (totalVehicles * 100) / maxTransactions

// 		occupancyLabel := "Lancar"
// 		if occupancyPercent >= 80 {
// 			occupancyLabel = "Zona Padat"
// 		} else if occupancyPercent >= 55 {
// 			occupancyLabel = "Ramai"
// 		} else if occupancyPercent >= 30 {
// 			occupancyLabel = "Normal"
// 		}
// 		loc.OccupancyLabel = occupancyLabel
// 	}

// 	return GetLocationsOutput{
// 		Locations: locations,
// 	}, nil
// }
