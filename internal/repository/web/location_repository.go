package repository

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"modulegue/internal/domain/web/location"
// 	// Import driver database kamu
// )

// type LocationRepository struct {
// 	db *sql.DB
// }

// func NewLocationRepository(db *sql.DB) location.Repository {
// 	return &LocationRepository{db: db}
// }

// // GetLocations mengambil semua lokasi dengan data agregat
// // Ini adalah adaptasi dari listLocationAggregates/buildParkingLocations di store.go
// func (r *LocationRepository) GetLocations(ctx context.Context) ([]*location.ParkingLocation, error) {
// 	// --- 1. Query Utama Lokasi ---
// 	// Ambil data dasar dari store.go: listLocationAggregates
// 	query := `
// 		SELECT
// 		    pl.id, pl.location_name, COALESCE(pz.zone_name, ''), COALESCE(pl.street_address, ''),
// 		    COALESCE(pl.operation_type, 'onstreet'),
// 		    COALESCE(pl.center_latitude, 0), COALESCE(pl.center_longitude, 0),
// 		    COALESCE(als.tariff_motor_amount, 0), COALESCE(als.tariff_car_amount, 0),
// 		    COALESCE(als.operational_note, '')
// 		FROM parking_location pl
// 		LEFT JOIN parking_zone pz ON pz.id = pl.zone_id
// 		LEFT JOIN (
// 		    SELECT location_id, tariff_motor_amount, tariff_car_amount, operational_note
// 		    FROM location_admin_settings -- Ganti nama tabel jika berbeda
// 		) als ON als.location_id = pl.id
// 		ORDER BY pl.id
// 	`

// 	rows, err := r.db.QueryContext(ctx, query)
// 	if err != nil {
// 		return nil, fmt.Errorf("query basic location data: %w", err)
// 	}
// 	defer rows.Close()

// 	// Ambil hasil query utama
// 	var rawLocations []locationAggregate // Gunakan struct sementara atau sesuaikan dengan domain
// 	for rows.Next() {
// 		var item locationAggregate // Asumsikan struct ini mirip dengan struct dari store.go
// 		err := rows.Scan(
// 			&item.ID, &item.Name, &item.Zone, &item.Address, &item.OperationType,
// 			&item.Latitude, &item.Longitude, &item.TariffMotor, &item.TariffMobil, &item.OperationalNote,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("scan basic location row: %w", err)
// 		}
// 		rawLocations = append(rawLocations, item)
// 	}
// 	if err = rows.Err(); err != nil {
// 		return nil, fmt.Errorf("iterate basic location rows: %w", err)
// 	}

// 	// Konversi ke []*domain.Location (inisialisasi dulu)
// 	locations := make([]*location.ParkingLocation, len(rawLocations))
// 	for i, rawLoc := range rawLocations {
// 		locations[i] = &location.ParkingLocation{
// 			ID:              fmt.Sprintf("%d", rawLoc.ID),
// 			Name:            rawLoc.Name,
// 			Zone:            rawLoc.Zone,
// 			Address:         rawLoc.Address,
// 			Latitude:        rawLoc.Latitude,
// 			Longitude:       rawLoc.Longitude,
// 			DismissalReason: rawLoc.OperationalNote,
// 			TariffMotor:     rawLoc.TariffMotor,
// 			TariffMobil:     rawLoc.TariffMobil,
// 			// Field lainnya diisi nanti
// 		}
// 	}

// 	// --- 2. Ambil dan Gabungkan Assignment Petugas ---
// 	// Ambil data assignment dari store.go: listLocationAggregates
// 	assignmentQuery := `
// 		SELECT
// 		    oac.location_id,
// 		    oac.operational_status,
// 		    COALESCE(su.full_name, ''),
// 		    COALESCE(pst.shift_name, ''),
// 		    COALESCE(DATE_FORMAT(pst.start_time, '%H:%i'), ''),
// 		    COALESCE(DATE_FORMAT(pst.end_time, '%H:%i'), '')
// 		FROM officer_assignment_current oac
// 		LEFT JOIN system_user su ON su.id = oac.officer_user_id
// 		LEFT JOIN parking_shift_template pst ON pst.id = oac.shift_template_id
// 		WHERE oac.effective_to IS NULL
// 		ORDER BY oac.location_id, oac.id
// 	`

// 	assignmentRows, err := r.db.QueryContext(ctx, assignmentQuery)
// 	if err != nil {
// 		return nil, fmt.Errorf("query officer assignment data: %w", err)
// 	}
// 	defer assignmentRows.Close()

// 	assignmentMap := make(map[int64][]struct {
// 		Status string
// 		Name   string
// 		Label  string
// 		Start  string
// 		End    string
// 	})

// 	for assignmentRows.Next() {
// 		var locationID int64
// 		var status, name, label, start, end string
// 		err := assignmentRows.Scan(&locationID, &status, &name, &label, &start, &end)
// 		if err != nil {
// 			return nil, fmt.Errorf("scan assignment row: %w", err)
// 		}
// 		assignmentMap[locationID] = append(assignmentMap[locationID], struct {
// 			Status string
// 			Name   string
// 			Label  string
// 			Start  string
// 			End    string
// 		}{
// 			Status: status,
// 			Name:   name,
// 			Label:  label,
// 			Start:  start,
// 			End:    end,
// 		})
// 	}
// 	if err = assignmentRows.Err(); err != nil {
// 		return nil, fmt.Errorf("iterate assignment rows: %w", err)
// 	}

// 	// Gabungkan assignment ke struct domain
// 	for i := range locations {
// 		assignments := assignmentMap[getIDInt64(locations[i].ID)] // Helper untuk parse ID string ke int64
// 		locations[i].Officers = int64(len(assignments))
// 		if len(assignments) > 0 {
// 			locations[i].OfficerName = assignments[0].Name
// 			locations[i].OfficerStatus = officerStatusLabel(assignments[0].Status) // Gunakan fungsi helper dari store.go
// 			locations[i].OfficerShiftStart = assignments[0].Start
// 			locations[i].OfficerShiftEnd = assignments[0].End
// 		} else {
// 			locations[i].OfficerName = "Belum ada jukir aktif"
// 			locations[i].OfficerStatus = "Istirahat"
// 		}
// 	}

// 	// --- 3. Ambil dan Gabungkan Shift Templates ---
// 	// Ambil data shift dari store.go: listLocationAggregates
// 	shiftQuery := `
// 		SELECT
// 		    pl.id, -- location_id
// 		    pst.id, -- shift_template_id
// 		    COALESCE(pst.shift_name, ''),
// 		    DATE_FORMAT(pst.start_time, '%H:%i'),
// 		    DATE_FORMAT(pst.end_time, '%H:%i')
// 		FROM parking_location pl
// 		LEFT JOIN parking_shift_template pst ON pst.location_id = pl.id AND pst.is_active = 1
// 		ORDER BY pl.id, pst.start_time, pst.id
// 	`

// 	shiftRows, err := r.db.QueryContext(ctx, shiftQuery)
// 	if err != nil {
// 		return nil, fmt.Errorf("query shift template data: %w", err)
// 	}
// 	defer shiftRows.Close()

// 	shiftMap := make(map[int64][]location.ParkingShiftTemplate)
// 	for shiftRows.Next() {
// 		var locationID int64
// 		var shiftID sql.NullInt64 // ID shift bisa null jika tidak ada template aktif
// 		var label, start, end sql.NullString
// 		err := shiftRows.Scan(&locationID, &shiftID, &label, &start, &end)
// 		if err != nil {
// 			return nil, fmt.Errorf("scan shift template row: %w", err)
// 		}
// 		if !shiftID.Valid {
// 			continue // Lewati jika tidak ada shift template
// 		}
// 		shiftMap[locationID] = append(shiftMap[locationID], location.ParkingShiftTemplate{
// 			ID:    fmt.Sprintf("%d", shiftID.Int64),
// 			Label: label.String,
// 			Start: start.String,
// 			End:   end.String,
// 		})
// 	}
// 	if err = shiftRows.Err(); err != nil {
// 		return nil, fmt.Errorf("iterate shift template rows: %w", err)
// 	}

// 	// Gabungkan shift template ke struct domain
// 	for i := range locations {
// 		locations[i].ShiftTemplates = shiftMap[getIDInt64(locations[i].ID)]
// 	}

// 	// --- 4. Ambil dan Gabungkan Transaksi/Statistik ---
// 	// Ambil data transaksi dari store.go: listLocationAggregates
// 	statsQuery := `
// 		SELECT
// 		    fpt.location_id,
// 		    SUM(CASE WHEN LOWER(COALESCE(vt.vehicle_type_name, '')) LIKE '%motor%' THEN 1 ELSE 0 END) AS motorcycles,
// 		    SUM(CASE WHEN LOWER(COALESCE(vt.vehicle_type_name, '')) LIKE '%car%' OR LOWER(COALESCE(vt.vehicle_type_name, '')) LIKE '%mobil%' THEN 1 ELSE 0 END) AS cars,
// 		    COUNT(fpt.id) AS transactions,
// 		    SUM(fpt.final_amount) AS revenue
// 		FROM financial_parking_transaction fpt
// 		LEFT JOIN vehicle_type vt ON vt.id = fpt.vehicle_type_id
// 		-- Tambahkan filter status dan waktu jika perlu, misalnya hanya 'paid'
// 		-- WHERE fpt.transaction_status = 'paid' AND ...
// 		GROUP BY fpt.location_id
// 	`

// 	statsRows, err := r.db.QueryContext(ctx, statsQuery)
// 	if err != nil {
// 		// Log error, lanjutkan dengan data sebelumnya
// 		log.Printf("Warning: Could not fetch location stats: %v", err)
// 		// Kita tetap lanjutkan, field motorcycle, cars, etc. akan kosong (0)
// 	} else {
// 		defer statsRows.Close()

// 		statsMap := make(map[int64]struct {
// 			Motorcycles  int64
// 			Cars         int64
// 			Transactions int64
// 			Revenue      int64
// 		})

// 		for statsRows.Next() {
// 			var locationID int64
// 			var motorcycles, cars, transactions, revenue int64
// 			err := statsRows.Scan(&locationID, &motorcycles, &cars, &transactions, &revenue)
// 			if err != nil {
// 				return nil, fmt.Errorf("scan location stats row: %w", err)
// 			}
// 			statsMap[locationID] = struct {
// 				Motorcycles  int64
// 				Cars         int64
// 				Transactions int64
// 				Revenue      int64
// 			}{
// 				Motorcycles:  motorcycles,
// 				Cars:         cars,
// 				Transactions: transactions,
// 				Revenue:      revenue,
// 			}
// 		}
// 		if err = statsRows.Err(); err != nil {
// 			return nil, fmt.Errorf("iterate location stats rows: %w", err)
// 		}

// 		// Gabungkan stats ke struct domain
// 		for i := range locations {
// 			stats := statsMap[getIDInt64(locations[i].ID)]
// 			locations[i].Motorcycles = stats.Motorcycles
// 			locations[i].Cars = stats.Cars
// 			// Revenue mungkin tidak langsung dimasukkan ke ParkingLocation jika tidak relevan
// 		}
// 	}

// 	// --- 5. Hitung Occupancy Label (menggunakan logika dari store.go) ---
// 	// Ambil max traffic global dulu
// 	maxTrafficGlobal := int64(1) // Default agar tidak dibagi nol
// 	for _, loc := range locations {
// 		totalVehicles := loc.Motorcycles + loc.Cars
// 		if totalVehicles > maxTrafficGlobal {
// 			maxTrafficGlobal = totalVehicles
// 		}
// 	}

// 	// Hitung label occupancy lokal berdasarkan max global
// 	for i := range locations {
// 		totalVehicles := locations[i].Cars + locations[i].Motorcycles
// 		occupancyPercent := (totalVehicles * 100) / maxTrafficGlobal
// 		occupancyLabel := "Lancar"
// 		if occupancyPercent >= 80 {
// 			occupancyLabel = "Zona Padat"
// 		} else if occupancyPercent >= 55 {
// 			occupancyLabel = "Ramai"
// 		} else if occupancyPercent >= 30 {
// 			occupancyLabel = "Normal"
// 		}
// 		locations[i].OccupancyLabel = occupancyLabel
// 	}

// 	// --- 6. Atur Tarif Default (menggunakan logika dari store.go) ---
// 	for i := range locations {
// 		if locations[i].TariffMotor <= 0 {
// 			// Ganti dengan logika RevenuePerVehicle jika diperlukan, atau nilai default
// 			// locations[i].TariffMotor = maxInt64(2000, locations[i].RevenuePerVehicle("motor")) // <-- Tidak bisa langsung, harus dihitung di sini atau query terpisah
// 			locations[i].TariffMotor = 2000 // Contoh default
// 		}
// 		if locations[i].TariffMobil <= 0 {
// 			// locations[i].TariffMobil = maxInt64(5000, locations[i].RevenuePerVehicle("mobil")) // <-- Sama seperti atas
// 			locations[i].TariffMobil = 5000 // Contoh default
// 		}
// 	}

// 	return locations, nil
// }

// // Helper untuk parse ID string ke int64
// func getIDInt64(idStr string) int64 {
// 	// Implementasi parsing, misalnya menggunakan strconv.ParseInt
// 	// Sederhanakan untuk saat ini
// 	var id int64
// 	fmt.Sscanf(idStr, "%d", &id)
// 	return id
// }

// // Helper untuk status (ambil dari store.go jika ada)
// func officerStatusLabel(status string) string {
// 	// Implementasi dari store.go: officerStatusLabel
// 	// Contoh sederhana:
// 	switch status {
// 	case "on_duty":
// 		return "Aktif Bertugas"
// 	case "off_duty":
// 		return "Istirahat"
// 	case "inactive":
// 		return "Tidak Aktif"
// 	default:
// 		return "Tidak Diketahui"
// 	}
// }

// // Helper untuk maxInt64 (ambil dari store.go jika ada)
// func maxInt64(a, b int64) int64 {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }

// // struct sementara untuk mapping hasil query awal
// // Ganti dengan struct dari domain jika semua field cocok
// type locationAggregate struct {
// 	ID              int64
// 	Name            string
// 	Zone            string
// 	Address         string
// 	OperationType   string
// 	Latitude        float64
// 	Longitude       float64
// 	TariffMotor     int64
// 	TariffMobil     int64
// 	OperationalNote string
// 	// Tambahkan field lain sesuai query awal
// 	Motorcycles int64
// 	Cars        int64
// 	// ... dan seterusnya
// }

// // ... (implementasi fungsi lain seperti GetLocationByID, GetHourlyTraffic, dll. mengikuti pola yang sama)
