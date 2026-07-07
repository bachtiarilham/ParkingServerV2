package payment

// import (
// 	"context"
// 	"fmt"
// 	"math"
// 	"modulegue/internal/domain/payment"
// 	"time"
// )

// type GetScanDetailInput struct {
// 	SessionID  int64
// 	CustomerID int64 // Didapat dari context JWT
// }

// type GetScanDetailOutput struct {
// 	CustomerID   string
// 	CustomerName string // Diambil dari user_repo
// 	LocationName string // Diambil dari parking_location
// 	Duration     string // Format HH:mm
// 	IsMember     bool   // Diambil dari user_repo
// 	Total        int64
// 	Breakdown    []PriceItemDto
// }

// type PriceItemDto struct {
// 	Label  string `json:"label"`
// 	Amount int64  `json:"amount"`
// }

// type GetScanDetailUseCase struct {
// 	repo payment.Repository
// 	// userRepo user.Repository // Untuk nama, member status
// 	// locationRepo location.Repository // Untuk nama lokasi
// }

// func NewGetScanDetailUseCase(repo payment.Repository /*, userRepo user.Repository, locationRepo location.Repository*/) *GetScanDetailUseCase {
// 	return &GetScanDetailUseCase{
// 		repo: repo,
// 		// userRepo: userRepo,
// 		// locationRepo: locationRepo,
// 	}
// }

// func (uc *GetScanDetailUseCase) Execute(ctx context.Context, input GetScanDetailInput) (GetScanDetailOutput, error) {
// 	// 1. Ambil detail sesi parkir
// 	session, err := uc.repo.GetActiveSessionByCode(ctx, fmt.Sprintf("%d", input.SessionID)) // Asumsi session.Code adalah string dari ID
// 	if err != nil {
// 		return GetScanDetailOutput{}, fmt.Errorf("session not found: %w", err)
// 	}

// 	// 2. Validasi kepemilikan (opsional)
// 	if session.CustomerID != input.CustomerID {
// 		// return GetScanDetailOutput{}, errors.New("session does not belong to this customer")
// 	}

// 	// 3. Hitung durasi (dari started_at ke sekarang)
// 	duration := time.Since(session.StartedAt)
// 	durationHours := math.Ceil(duration.Hours()) // Tarif biasanya per jam, dibulatkan ke atas

// 	// 4. Ambil tariff
// 	tariff, err := uc.repo.GetTariffForLocationAndVehicle(ctx, session.LocationID, session.VehicleTypeID)
// 	if err != nil {
// 		return GetScanDetailOutput{}, fmt.Errorf("failed to get tariff: %w", err)
// 	}

// 	// 5. Hitung total (sederhana: tariff * hours)
// 	total := int64(durationHours) * tariff

// 	// 6. Buat breakdown (contoh sederhana)
// 	breakdown := []PriceItemDto{
// 		{Label: fmt.Sprintf("Parking Fee (%.0fh x Rp %d)", durationHours, tariff), Amount: total},
// 		// Tambahkan diskon/penalty jika ada
// 	}

// 	// 7. Ambil info tambahan (nama customer, nama lokasi, status member)
// 	// customer, _ := uc.userRepo.GetByID(ctx, session.CustomerID) // Implementasi di user_repo
// 	// location, _ := uc.locationRepo.GetByID(ctx, session.LocationID) // Implementasi di location_repo

// 	// 8. Format durasi ke HH:mm
// 	hours := int(duration.Hours())
// 	minutes := int(duration.Minutes()) % 60
// 	durationStr := fmt.Sprintf("%02dh %02dm", hours, minutes)

// 	return GetScanDetailOutput{
// 		CustomerID:   fmt.Sprintf("%d", session.CustomerID),
// 		CustomerName: "Customer Name Placeholder", // Ambil dari customer entity
// 		LocationName: "Location Name Placeholder", // Ambil dari location entity
// 		Duration:     durationStr,
// 		IsMember:     false, // Ambil dari customer entity
// 		Total:        total,
// 		Breakdown:    breakdown,
// 	}, nil
// }
