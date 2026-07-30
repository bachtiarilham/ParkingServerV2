package riwayat

import "time"

// Response utama API Riwayat Parkir
type RiwayatParkirModel struct {
	Summary ParkingSummaryModel      `json:"summary"`
	Items   []RiwayatParkirItemModel `json:"items"`
}

// ParkingSummary merepresentasikan data statistik di Top Metric Card
type ParkingSummaryModel struct {
	TotalDurationMinutes int    `json:"total_duration_minutes"` // e.g. 1125
	TotalDurationText    string `json:"total_duration_text"`    // e.g. "18 Jam 45 Mnt"
	TotalAmount          int64  `json:"total_amount"`           // e.g. 120000
	CompletedCount       int    `json:"completed_count"`        // e.g. 12
}

// ParkingHistoryItem merepresentasikan tiap item transaksi parkir di dalam list
type RiwayatParkirItemModel struct {
	ID           string     `json:"id"`
	TicketNo     string     `json:"ticket_no"`     // e.g. "TKT-2026-0091"
	LocationName string     `json:"location_name"` // e.g. "On-Street Area A (Jl. Sudirman)"
	ParkingType  string     `json:"parking_type"`  // "ON_STREET" | "OFF_STREET"
	LicensePlate string     `json:"license_plate"` // e.g. "B 1234 ABC"
	VehicleType  string     `json:"vehicle_type"`  // "Motor" | "Mobil"
	CheckInTime  time.Time  `json:"check_in_time"`
	CheckOutTime *time.Time `json:"check_out_time,omitempty"` // Nullable jika status BERLANGSUNG
	DurationText string     `json:"duration_text"`            // e.g. "1 Jam 20 Menit"
	Amount       int64      `json:"amount"`
	Status       string     `json:"status"` // "SELESAI" | "BERLANGSUNG" | "BATAL"
}
