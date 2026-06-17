package store

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	model "github.com/bachtiarilham/ParkServerFinal/internal/dashboardparkir/contracts"
)

type MySQLRepository struct {
	db              *sql.DB
	adminTablesOnce sync.Once
	adminTablesErr  error
}

func NewMySQL(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

type Service = MySQLRepository

func New(db *sql.DB) *Service {
	return NewMySQL(db)
}

type locationAggregate struct {
	ID                int64
	Slug              string
	Name              string
	Zone              string
	Address           string
	OperationType     string
	Lat               float64
	Lng               float64
	Motorcycles       int64
	Cars              int64
	Transactions      int64
	Revenue           int64
	Officers          int64
	OfficerName       string
	OfficerStatus     string
	OfficerShiftStart string
	OfficerShiftEnd   string
	TariffMotor       int64
	TariffMobil       int64
	OperationalNote   string
	ShiftTemplates    []model.ParkingShiftTemplate
}

func (r *MySQLRepository) ensureAdminTables(ctx context.Context) error {
	r.adminTablesOnce.Do(func() {
		statements := []string{
			`CREATE TABLE IF NOT EXISTS admin_location_setting (
				location_id BIGINT UNSIGNED NOT NULL,
				tariff_motor_amount BIGINT NOT NULL DEFAULT 0,
				tariff_car_amount BIGINT NOT NULL DEFAULT 0,
				operational_note TEXT NULL,
				updated_by_user_id BIGINT UNSIGNED NULL,
				updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				PRIMARY KEY (location_id),
				CONSTRAINT fk_admin_location_setting_location FOREIGN KEY (location_id) REFERENCES parking_location(id) ON DELETE RESTRICT ON UPDATE RESTRICT,
				CONSTRAINT fk_admin_location_setting_updated_by FOREIGN KEY (updated_by_user_id) REFERENCES system_user(id) ON DELETE RESTRICT ON UPDATE RESTRICT
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
			`CREATE TABLE IF NOT EXISTS admin_settings_alert_rule (
				id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
				rule_title VARCHAR(191) NOT NULL,
				threshold_text VARCHAR(255) NOT NULL,
				source_text VARCHAR(255) NOT NULL,
				pic_text VARCHAR(255) NOT NULL,
				sort_order INT NOT NULL DEFAULT 0,
				is_active TINYINT(1) NOT NULL DEFAULT 1,
				updated_by_user_id BIGINT UNSIGNED NULL,
				updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				PRIMARY KEY (id),
				UNIQUE KEY uk_admin_settings_alert_rule_title (rule_title),
				CONSTRAINT fk_admin_settings_alert_rule_updated_by FOREIGN KEY (updated_by_user_id) REFERENCES system_user(id) ON DELETE RESTRICT ON UPDATE RESTRICT
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
			`CREATE TABLE IF NOT EXISTS admin_settings_shift_template (
				id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
				shift_label VARCHAR(191) NOT NULL,
				start_time CHAR(5) NOT NULL,
				end_time CHAR(5) NOT NULL,
				use_case VARCHAR(255) NOT NULL,
				sort_order INT NOT NULL DEFAULT 0,
				is_active TINYINT(1) NOT NULL DEFAULT 1,
				updated_by_user_id BIGINT UNSIGNED NULL,
				updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				PRIMARY KEY (id),
				CONSTRAINT fk_admin_settings_shift_template_updated_by FOREIGN KEY (updated_by_user_id) REFERENCES system_user(id) ON DELETE RESTRICT ON UPDATE RESTRICT
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
			`CREATE TABLE IF NOT EXISTS admin_settings_notification (
				id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
				channel_name VARCHAR(191) NOT NULL,
				trigger_text VARCHAR(255) NOT NULL,
				response_text VARCHAR(255) NOT NULL,
				sort_order INT NOT NULL DEFAULT 0,
				is_active TINYINT(1) NOT NULL DEFAULT 1,
				updated_by_user_id BIGINT UNSIGNED NULL,
				updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				PRIMARY KEY (id),
				CONSTRAINT fk_admin_settings_notification_updated_by FOREIGN KEY (updated_by_user_id) REFERENCES system_user(id) ON DELETE RESTRICT ON UPDATE RESTRICT
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
			`CREATE TABLE IF NOT EXISTS admin_settings_payment_method (
				method_code VARCHAR(64) NOT NULL,
				method_label VARCHAR(191) NOT NULL,
				icon_code VARCHAR(32) NOT NULL,
				is_enabled TINYINT(1) NOT NULL DEFAULT 1,
				updated_by_user_id BIGINT UNSIGNED NULL,
				updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				PRIMARY KEY (method_code),
				CONSTRAINT fk_admin_settings_payment_method_updated_by FOREIGN KEY (updated_by_user_id) REFERENCES system_user(id) ON DELETE RESTRICT ON UPDATE RESTRICT
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
			`CREATE TABLE IF NOT EXISTS admin_settings_default_tariff (
				vehicle_type_code VARCHAR(64) NOT NULL,
				vehicle_type_label VARCHAR(191) NOT NULL,
				first_hour_amount BIGINT NOT NULL DEFAULT 0,
				next_hour_amount BIGINT NOT NULL DEFAULT 0,
				max_rate_amount BIGINT NOT NULL DEFAULT 0,
				updated_by_user_id BIGINT UNSIGNED NULL,
				updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				PRIMARY KEY (vehicle_type_code),
				CONSTRAINT fk_admin_settings_default_tariff_updated_by FOREIGN KEY (updated_by_user_id) REFERENCES system_user(id) ON DELETE RESTRICT ON UPDATE RESTRICT
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		}
		for _, statement := range statements {
			if _, err := r.db.ExecContext(ctx, statement); err != nil {
				r.adminTablesErr = err
				return
			}
		}
	})
	return r.adminTablesErr
}

func parseIDString(value string) (int64, error) {
	parsed, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil || parsed <= 0 {
		return 0, errors.New("id tidak valid")
	}
	return parsed, nil
}

func shiftCodeFromTemplate(shift model.ParkingShiftTemplate) string {
	label := strings.ToUpper(strings.TrimSpace(shift.Label))
	label = strings.ReplaceAll(label, " ", "_")
	label = strings.ReplaceAll(label, "-", "_")
	if label == "" {
		label = "SHIFT"
	}
	return fmt.Sprintf("%s_%s_%s", label, strings.ReplaceAll(shift.Start, ":", ""), strings.ReplaceAll(shift.End, ":", ""))
}

func toOperationalStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "aktif":
		return "on_duty"
	case "istirahat":
		return "off_duty"
	case "diberhentikan":
		return "inactive"
	default:
		return "off_duty"
	}
}

func formatIDR(value int64) string {
	sign := ""
	if value < 0 {
		sign = "-"
		value = -value
	}
	raw := fmt.Sprintf("%d", value)
	if len(raw) <= 3 {
		return sign + "Rp " + raw
	}
	parts := make([]string, 0, (len(raw)+2)/3)
	for len(raw) > 3 {
		parts = append([]string{raw[len(raw)-3:]}, parts...)
		raw = raw[:len(raw)-3]
	}
	parts = append([]string{raw}, parts...)
	return sign + "Rp " + strings.Join(parts, ".")
}

func splitHoursRange(raw string) (string, string) {
	parts := strings.Split(raw, "-")
	if len(parts) != 2 {
		return "", ""
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

func paymentMethodCode(label string) string {
	normalized := strings.ToLower(strings.TrimSpace(label))
	switch {
	case strings.Contains(normalized, "qris"):
		return "qris"
	case strings.Contains(normalized, "cash"):
		return "cash"
	default:
		return strings.ReplaceAll(normalized, " ", "_")
	}
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	replacer := strings.NewReplacer(
		" ", "-",
		".", "",
		",", "",
		"/", "-",
	)
	value = replacer.Replace(value)
	for strings.Contains(value, "--") {
		value = strings.ReplaceAll(value, "--", "-")
	}
	return strings.Trim(value, "-")
}

func titleRole(roleCode string) string {
	switch strings.ToLower(strings.TrimSpace(roleCode)) {
	case "jukir", "officer":
		return "Petugas Jukir"
	case "sco":
		return "Supervisor"
	case "admin_ops":
		return "Admin Operations"
	case "admin_finance":
		return "Admin Finance"
	case "admin_super":
		return "Super Admin"
	default:
		return strings.ReplaceAll(strings.Title(strings.ReplaceAll(roleCode, "_", " ")), "Ops", "Ops")
	}
}

func (r *MySQLRepository) FindAdminByIdentity(ctx context.Context, identity string) (model.AuthRecord, error) {
	identity = strings.TrimSpace(identity)
	row := r.db.QueryRowContext(ctx, `
		SELECT
			su.id,
			COALESCE(su.full_name, ''),
			COALESCE(su.phone_number, ''),
			COALESCE(su.email, ''),
			COALESCE(su.username, ''),
			COALESCE(sr.role_code, ''),
			COALESCE(su.password_hash, ''),
			COALESCE(su.is_verified, 0)
		FROM system_user su
		JOIN system_role sr ON sr.id = su.role_id
		WHERE su.username = ? OR su.email = ? OR su.phone_number = ?
		LIMIT 1`, identity, identity, identity)

	var item model.AuthRecord
	var verified bool
	if err := row.Scan(&item.ID, &item.FullName, &item.PhoneNumber, &item.Email, &item.Username, &item.RoleCode, &item.PasswordHash, &verified); err != nil {
		return model.AuthRecord{}, err
	}
	item.IsVerified = verified
	return item, nil
}

func (r *MySQLRepository) FindAdminByID(ctx context.Context, userID int64) (model.AuthRecord, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT
			su.id,
			COALESCE(su.full_name, ''),
			COALESCE(su.phone_number, ''),
			COALESCE(su.email, ''),
			COALESCE(su.username, ''),
			COALESCE(sr.role_code, ''),
			COALESCE(su.password_hash, ''),
			COALESCE(su.is_verified, 0)
		FROM system_user su
		JOIN system_role sr ON sr.id = su.role_id
		WHERE su.id = ?
		LIMIT 1`, userID)

	var item model.AuthRecord
	var verified bool
	if err := row.Scan(&item.ID, &item.FullName, &item.PhoneNumber, &item.Email, &item.Username, &item.RoleCode, &item.PasswordHash, &verified); err != nil {
		return model.AuthRecord{}, err
	}
	item.IsVerified = verified
	return item, nil
}

func officerStatusLabel(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "on_duty":
		return "Aktif"
	case "off_duty":
		return "Istirahat"
	case "inactive":
		return "Diberhentikan"
	default:
		return "Aktif"
	}
}

func toneFromRank(index int) string {
	switch index {
	case 0:
		return "orange"
	case 1:
		return "blue"
	case 2:
		return "green"
	default:
		return "gold"
	}
}

func (r *MySQLRepository) listLocationAggregates(ctx context.Context) ([]locationAggregate, error) {
	if err := r.ensureAdminTables(ctx); err != nil {
		return nil, err
	}
	query := `
		SELECT
			pl.id,
			pl.location_name,
			COALESCE(pz.zone_name, ''),
			COALESCE(pl.street_address, ''),
			COALESCE(pl.operation_type, 'onstreet'),
			COALESCE(pl.center_latitude, 0),
			COALESCE(pl.center_longitude, 0),
			COALESCE(als.tariff_motor_amount, 0),
			COALESCE(als.tariff_car_amount, 0),
			COALESCE(als.operational_note, ''),
			SUM(CASE WHEN LOWER(COALESCE(vt.vehicle_type_name, '')) LIKE '%motor%' THEN 1 ELSE 0 END) AS motorcycles,
			SUM(CASE WHEN LOWER(COALESCE(vt.vehicle_type_name, '')) LIKE '%mobil%' OR LOWER(COALESCE(vt.vehicle_type_name, '')) LIKE '%car%' THEN 1 ELSE 0 END) AS cars,
			COUNT(fpt.id) AS transactions,
			COALESCE(SUM(CASE WHEN fpt.transaction_status IN ('paid', 'refunded_full', 'refunded_partial') THEN fpt.final_amount ELSE 0 END), 0) AS revenue
		FROM parking_location pl
		LEFT JOIN parking_zone pz ON pz.id = pl.zone_id
		LEFT JOIN admin_location_setting als ON als.location_id = pl.id
		LEFT JOIN financial_parking_transaction fpt
			ON fpt.location_id = pl.id
		LEFT JOIN vehicle_type vt ON vt.id = fpt.vehicle_type_id
		GROUP BY pl.id, pl.location_name, pz.zone_name, pl.street_address, pl.operation_type, pl.center_latitude, pl.center_longitude, als.tariff_motor_amount, als.tariff_car_amount, als.operational_note
		ORDER BY pl.id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	aggregates := make([]locationAggregate, 0)
	for rows.Next() {
		var item locationAggregate
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Zone,
			&item.Address,
			&item.OperationType,
			&item.Lat,
			&item.Lng,
			&item.TariffMotor,
			&item.TariffMobil,
			&item.OperationalNote,
			&item.Motorcycles,
			&item.Cars,
			&item.Transactions,
			&item.Revenue,
		); err != nil {
			return nil, err
		}
		item.Slug = fmt.Sprintf("%d", item.ID)
		aggregates = append(aggregates, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	assignmentQuery := `
		SELECT
			oac.location_id,
			oac.operational_status,
			COALESCE(su.full_name, ''),
			COALESCE(pst.shift_name, ''),
			COALESCE(DATE_FORMAT(pst.start_time, '%H:%i'), ''),
			COALESCE(DATE_FORMAT(pst.end_time, '%H:%i'), '')
		FROM officer_assignment_current oac
		LEFT JOIN system_user su ON su.id = oac.officer_user_id
		LEFT JOIN parking_shift_template pst ON pst.id = oac.shift_template_id
		WHERE oac.effective_to IS NULL
		ORDER BY oac.location_id, oac.id`
	assignmentRows, err := r.db.QueryContext(ctx, assignmentQuery)
	if err != nil {
		return nil, err
	}
	defer assignmentRows.Close()

	assignmentMap := map[int64][]struct {
		Status string
		Name   string
		Label  string
		Start  string
		End    string
	}{}
	for assignmentRows.Next() {
		var locationID int64
		var status, name, label, start, end string
		if err := assignmentRows.Scan(&locationID, &status, &name, &label, &start, &end); err != nil {
			return nil, err
		}
		assignmentMap[locationID] = append(assignmentMap[locationID], struct {
			Status string
			Name   string
			Label  string
			Start  string
			End    string
		}{Status: status, Name: name, Label: label, Start: start, End: end})
	}
	if err := assignmentRows.Err(); err != nil {
		return nil, err
	}

	shiftQuery := `
		SELECT
			pl.id,
			pst.id,
			COALESCE(pst.shift_name, ''),
			DATE_FORMAT(pst.start_time, '%H:%i'),
			DATE_FORMAT(pst.end_time, '%H:%i')
		FROM parking_location pl
		LEFT JOIN parking_shift_template pst ON pst.location_id = pl.id AND pst.is_active = 1
		ORDER BY pl.id, pst.start_time, pst.id`
	shiftRows, err := r.db.QueryContext(ctx, shiftQuery)
	if err != nil {
		return nil, err
	}
	defer shiftRows.Close()

	shiftMap := map[int64][]model.ParkingShiftTemplate{}
	for shiftRows.Next() {
		var locationID int64
		var shiftID sql.NullInt64
		var label, start, end sql.NullString
		if err := shiftRows.Scan(&locationID, &shiftID, &label, &start, &end); err != nil {
			return nil, err
		}
		if !shiftID.Valid {
			continue
		}
		shiftMap[locationID] = append(shiftMap[locationID], model.ParkingShiftTemplate{
			ID:    fmt.Sprintf("%d", shiftID.Int64),
			Label: label.String,
			Start: start.String,
			End:   end.String,
		})
	}
	if err := shiftRows.Err(); err != nil {
		return nil, err
	}

	for i := range aggregates {
		assignments := assignmentMap[aggregates[i].ID]
		aggregates[i].Officers = int64(len(assignments))
		if len(assignments) > 0 {
			aggregates[i].OfficerName = assignments[0].Name
			aggregates[i].OfficerStatus = officerStatusLabel(assignments[0].Status)
			aggregates[i].OfficerShiftStart = assignments[0].Start
			aggregates[i].OfficerShiftEnd = assignments[0].End
		} else {
			aggregates[i].OfficerName = "Belum ada jukir aktif"
			aggregates[i].OfficerStatus = "Istirahat"
		}
		aggregates[i].ShiftTemplates = shiftMap[aggregates[i].ID]
		if aggregates[i].TariffMotor <= 0 {
			aggregates[i].TariffMotor = maxInt64(2000, aggregates[i].RevenuePerVehicle("motor"))
		}
		if aggregates[i].TariffMobil <= 0 {
			aggregates[i].TariffMobil = maxInt64(5000, aggregates[i].RevenuePerVehicle("mobil"))
		}
	}

	return aggregates, nil
}

func (l locationAggregate) RevenuePerVehicle(kind string) int64 {
	switch kind {
	case "motor":
		if l.Motorcycles <= 0 {
			return 0
		}
		return l.Revenue / maxInt64(1, l.Motorcycles+l.Cars)
	case "mobil":
		if l.Cars <= 0 {
			return 0
		}
		return l.Revenue / maxInt64(1, l.Motorcycles+l.Cars)
	default:
		return 0
	}
}

func maxInt64(values ...int64) int64 {
	var max int64
	for i, value := range values {
		if i == 0 || value > max {
			max = value
		}
	}
	return max
}

func (r *MySQLRepository) buildParkingLocations(ctx context.Context) ([]locationAggregate, error) {
	return r.listLocationAggregates(ctx)
}

func toParkingLocation(item locationAggregate, maxTransactions int64) model.ParkingLocation {
	occupancyPercent := int64(0)
	if maxTransactions > 0 {
		occupancyPercent = (item.Transactions * 100) / maxTransactions
	}
	occupancyLabel := "Lancar"
	if occupancyPercent >= 80 {
		occupancyLabel = "Zona Padat"
	} else if occupancyPercent >= 55 {
		occupancyLabel = "Ramai"
	} else if occupancyPercent >= 30 {
		occupancyLabel = "Normal"
	}
	return model.ParkingLocation{
		ID:                item.Slug,
		Name:              item.Name,
		Zone:              item.Zone,
		Address:           item.Address,
		Lat:               item.Lat,
		Lng:               item.Lng,
		OfficerName:       item.OfficerName,
		OfficerShiftStart: item.OfficerShiftStart,
		OfficerShiftEnd:   item.OfficerShiftEnd,
		OfficerStatus:     item.OfficerStatus,
		DismissalReason:   item.OperationalNote,
		TariffMotor:       item.TariffMotor,
		TariffMobil:       item.TariffMobil,
		Motorcycles:       item.Motorcycles,
		Cars:              item.Cars,
		Officers:          item.Officers,
		OccupancyLabel:    occupancyLabel,
		ShiftTemplates:    item.ShiftTemplates,
	}
}

func (r *MySQLRepository) listOfficerOptions(ctx context.Context, locations []locationAggregate) ([]model.ParkingOfficerOption, error) {
	locationByID := make(map[int64]locationAggregate, len(locations))
	for _, location := range locations {
		locationByID[location.ID] = location
	}
	query := `
		SELECT
			su.id,
			COALESCE(su.full_name, ''),
			COALESCE(sr.role_code, ''),
			COALESCE(home_zone.zone_name, current_zone.zone_name, ''),
			COALESCE(current_location.location_name, ''),
			COALESCE(current_location.id, 0),
			COALESCE(pst.id, 0),
			COALESCE(DATE_FORMAT(pst.start_time, '%H:%i'), ''),
			COALESCE(DATE_FORMAT(pst.end_time, '%H:%i'), ''),
			COALESCE(oac.operational_status, osh.new_operational_status, 'off_duty')
		FROM system_user su
		LEFT JOIN system_role sr ON sr.id = su.role_id
		LEFT JOIN officer_assignment_current oac
			ON oac.officer_user_id = su.id
		   AND oac.effective_to IS NULL
		LEFT JOIN (
			SELECT h1.officer_user_id, h1.new_operational_status
			FROM officer_status_history h1
			INNER JOIN (
				SELECT officer_user_id, MAX(id) AS max_id
				FROM officer_status_history
				GROUP BY officer_user_id
			) h2
				ON h2.officer_user_id = h1.officer_user_id
			   AND h2.max_id = h1.id
		) osh ON osh.officer_user_id = su.id
		LEFT JOIN parking_location current_location ON current_location.id = oac.location_id
		LEFT JOIN parking_zone current_zone ON current_zone.id = current_location.zone_id
		LEFT JOIN user_region_scope urs ON urs.user_id = su.id
		LEFT JOIN parking_zone home_zone ON home_zone.id = urs.zone_id
		LEFT JOIN parking_shift_template pst ON pst.id = oac.shift_template_id
		WHERE COALESCE(sr.role_code, '') IN ('officer', 'sco')
		ORDER BY sr.role_code, su.full_name`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.ParkingOfficerOption, 0)
	for rows.Next() {
		var (
			id                                                                             int64
			name, roleCode, homeZone, currentLocationName, shiftStart, shiftEnd, statusRaw string
			currentLocationID, shiftID                                                     int64
		)
		if err := rows.Scan(&id, &name, &roleCode, &homeZone, &currentLocationName, &currentLocationID, &shiftID, &shiftStart, &shiftEnd, &statusRaw); err != nil {
			return nil, err
		}
		availability := "Tersedia"
		availabilityNote := "Petugas siap dirotasi mengikuti kebutuhan lokasi aktif."
		if currentLocationID > 0 {
			availability = "Bertugas di Lokasi Lain"
			availabilityNote = "Petugas sedang terikat di lokasi aktif dan perlu validasi saat dipindahkan."
		}
		if roleCode == "sco" {
			availability = "Cadangan"
			availabilityNote = "Supervisor dipakai sebagai fallback operasional bila diperlukan."
		}
		currentAssignment := "Belum ditempatkan"
		if currentLocationName != "" {
			currentAssignment = currentLocationName
		}
		if homeZone == "" && currentLocationID > 0 {
			homeZone = locationByID[currentLocationID].Zone
		}
		items = append(items, model.ParkingOfficerOption{
			ID:                fmt.Sprintf("%d", id),
			Name:              name,
			Role:              titleRole(roleCode),
			HomeZone:          homeZone,
			Availability:      availability,
			AvailabilityNote:  availabilityNote,
			CurrentAssignment: currentAssignment,
			CurrentLocationID: fmt.Sprintf("%d", currentLocationID),
			CurrentShiftID:    fmt.Sprintf("%d", shiftID),
			Status:            officerStatusLabel(statusRaw),
			DefaultShiftStart: shiftStart,
			DefaultShiftEnd:   shiftEnd,
			DefaultStatus:     officerStatusLabel(statusRaw),
		})
	}
	return items, rows.Err()
}

func todayWindow(now time.Time) (time.Time, time.Time) {
	location := now.Location()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	return start, start.Add(24 * time.Hour)
}

func yesterdayWindow(now time.Time) (time.Time, time.Time) {
	todayStart, _ := todayWindow(now)
	return todayStart.Add(-24 * time.Hour), todayStart
}

func dayNameID(t time.Time) string {
	switch t.Weekday() {
	case time.Monday:
		return "Sen"
	case time.Tuesday:
		return "Sel"
	case time.Wednesday:
		return "Rab"
	case time.Thursday:
		return "Kam"
	case time.Friday:
		return "Jum"
	case time.Saturday:
		return "Sab"
	default:
		return "Min"
	}
}

func (r *MySQLRepository) GetDashboardOverview(ctx context.Context) (model.DashboardOverview, error) {
	now := time.Now()
	locations, err := r.buildParkingLocations(ctx)
	if err != nil {
		return model.DashboardOverview{}, err
	}
	maxTransactions := int64(1)
	for _, location := range locations {
		if location.Transactions > maxTransactions {
			maxTransactions = location.Transactions
		}
	}
	parkingLocations := make([]model.ParkingLocation, 0, len(locations))
	for _, location := range locations {
		parkingLocations = append(parkingLocations, toParkingLocation(location, maxTransactions))
	}

	var walletBalance, totalTransactions, activeSessions, officerCount int64
	todayStart, todayEnd := todayWindow(now)
	_ = r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(current_balance_amount), 0) FROM user_wallet`).Scan(&walletBalance)
	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM financial_parking_transaction WHERE paid_at >= ? AND paid_at < ?`, todayStart, todayEnd).Scan(&totalTransactions)
	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM parking_session WHERE session_status = 'active'`).Scan(&activeSessions)
	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM system_user su JOIN system_role sr ON sr.id = su.role_id WHERE sr.role_code IN ('officer','sco') AND su.employment_status = 'active'`).Scan(&officerCount)

	stats := []model.StatCard{
		{Label: "Saldo", Value: formatIDR(walletBalance), Icon: "SL", Tone: "blue"},
		{Label: "Total Transaksi", Value: fmt.Sprintf("%d Transaksi", totalTransactions), Icon: "TX", Tone: "orange"},
		{Label: "Petugas", Value: fmt.Sprintf("%d Petugas", officerCount), Icon: "PT", Tone: "blue"},
		{Label: "Sedang Parkir", Value: fmt.Sprintf("%d", activeSessions), Icon: "PK", Tone: "cyan"},
	}

	txQuery := `
		SELECT
			COALESCE(fpt.plate_number, ''),
			COALESCE(vt.vehicle_type_name, ''),
			DATE_FORMAT(fpt.paid_at, '%H:%i'),
			COALESCE(fpt.final_amount, 0)
		FROM financial_parking_transaction fpt
		LEFT JOIN vehicle_type vt ON vt.id = fpt.vehicle_type_id
		WHERE fpt.paid_at IS NOT NULL
		ORDER BY fpt.paid_at DESC, fpt.id DESC
		LIMIT 5`
	rows, err := r.db.QueryContext(ctx, txQuery)
	if err != nil {
		return model.DashboardOverview{}, err
	}
	defer rows.Close()
	dashboardTransactions := make([]model.RowItem, 0, 5)
	for rows.Next() {
		var plate, vehicleType, paidTime string
		var amount int64
		if err := rows.Scan(&plate, &vehicleType, &paidTime, &amount); err != nil {
			return model.DashboardOverview{}, err
		}
		dashboardTransactions = append(dashboardTransactions, model.RowItem{
			Primary:   plate,
			Secondary: vehicleType,
			Time:      paidTime,
			Price:     formatIDR(amount),
		})
	}

	officers, err := r.listOfficerOptions(ctx, locations)
	if err != nil {
		return model.DashboardOverview{}, err
	}
	fieldOfficers := make([]model.RowItem, 0, 3)
	for i, officer := range officers {
		if i >= 3 {
			break
		}
		fieldOfficers = append(fieldOfficers, model.RowItem{
			Primary:    officer.Name,
			Secondary:  fmt.Sprintf("%s · %s", officer.Role, officer.CurrentAssignment),
			Status:     officer.Status,
			StatusTone: mapStatusTone(officer.Status),
			Time:       officer.DefaultShiftStart,
		})
	}

	alertsQuery := `
		SELECT severity, alert_title, COALESCE(alert_detail, '')
		FROM admin_alert_event
		ORDER BY triggered_at DESC, id DESC
		LIMIT 3`
	alertRows, err := r.db.QueryContext(ctx, alertsQuery)
	if err != nil {
		return model.DashboardOverview{}, err
	}
	defer alertRows.Close()
	dashboardAlerts := make([]model.AlertItem, 0, 3)
	for alertRows.Next() {
		var severity, title, detail string
		if err := alertRows.Scan(&severity, &title, &detail); err != nil {
			return model.DashboardOverview{}, err
		}
		dashboardAlerts = append(dashboardAlerts, model.AlertItem{
			Title:    title,
			Detail:   detail,
			Priority: normalizePriority(severity),
			Action:   actionForPriority(severity),
		})
	}
	if len(dashboardAlerts) == 0 {
		dashboardAlerts = []model.AlertItem{
			{Title: "Belum ada alert aktif", Detail: "Sistem admin belum menemukan anomali operasional untuk hari ini.", Priority: "low", Action: "Pantau monitoring lokasi"},
		}
	}

	hourlyTraffic, err := r.buildHourlyTraffic(ctx, todayStart, todayEnd)
	if err != nil {
		return model.DashboardOverview{}, err
	}
	revenueByLocation, occupancyByLocation, priorityActions := buildLocationMetrics(parkingLocations)
	comparisonMetrics, err := r.buildComparisonMetrics(ctx, now)
	if err != nil {
		return model.DashboardOverview{}, err
	}
	heatmap, err := r.buildHeatmap(ctx, now)
	if err != nil {
		return model.DashboardOverview{}, err
	}

	return model.DashboardOverview{
		DashboardStats:        stats,
		DashboardTransactions: dashboardTransactions,
		FieldOfficers:         fieldOfficers,
		DashboardAlerts:       dashboardAlerts,
		HourlyTraffic:         hourlyTraffic,
		RevenueByLocation:     revenueByLocation,
		OccupancyByLocation:   occupancyByLocation,
		ComparisonMetrics:     comparisonMetrics,
		ParkingHeatmap:        heatmap,
		PriorityActions:       priorityActions,
		ParkingLocations:      parkingLocations,
	}, nil
}

func mapStatusTone(status string) string {
	switch status {
	case "Aktif":
		return "green"
	case "Istirahat":
		return "blue"
	case "Diberhentikan":
		return "gray"
	default:
		return "blue"
	}
}

func normalizePriority(severity string) string {
	switch strings.ToLower(strings.TrimSpace(severity)) {
	case "high":
		return "high"
	case "medium":
		return "medium"
	default:
		return "low"
	}
}

func actionForPriority(severity string) string {
	switch normalizePriority(severity) {
	case "high":
		return "Naikkan prioritas monitoring"
	case "medium":
		return "Atur rotasi shift"
	default:
		return "Tinjau konfigurasi"
	}
}

func (r *MySQLRepository) buildHourlyTraffic(ctx context.Context, start, end time.Time) ([]model.HourlyTrafficPoint, error) {
	points := []model.HourlyTrafficPoint{
		{Label: "06:00"}, {Label: "08:00"}, {Label: "10:00"}, {Label: "12:00"}, {Label: "14:00"}, {Label: "16:00"},
	}
	query := `SELECT HOUR(COALESCE(started_at, occurred_at)) AS h, COUNT(*) FROM (
		SELECT started_at, NULL AS occurred_at FROM parking_session WHERE started_at >= ? AND started_at < ?
		UNION ALL
		SELECT NULL AS started_at, occurred_at FROM financial_parking_transaction WHERE occurred_at >= ? AND occurred_at < ?
	) q GROUP BY HOUR(COALESCE(started_at, occurred_at))`
	rows, err := r.db.QueryContext(ctx, query, start, end, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	incoming := map[int]int64{}
	for rows.Next() {
		var hour int
		var count int64
		if err := rows.Scan(&hour, &count); err != nil {
			return nil, err
		}
		incoming[hour] = count
	}
	exitQuery := `SELECT HOUR(COALESCE(ended_at, paid_at)) AS h, COUNT(*) FROM (
		SELECT ended_at, NULL AS paid_at FROM parking_session WHERE ended_at IS NOT NULL AND ended_at >= ? AND ended_at < ?
		UNION ALL
		SELECT NULL AS ended_at, paid_at FROM financial_parking_transaction WHERE paid_at IS NOT NULL AND paid_at >= ? AND paid_at < ?
	) q GROUP BY HOUR(COALESCE(ended_at, paid_at))`
	exitRows, err := r.db.QueryContext(ctx, exitQuery, start, end, start, end)
	if err != nil {
		return nil, err
	}
	defer exitRows.Close()
	outgoing := map[int]int64{}
	for exitRows.Next() {
		var hour int
		var count int64
		if err := exitRows.Scan(&hour, &count); err != nil {
			return nil, err
		}
		outgoing[hour] = count
	}
	hours := []int{6, 8, 10, 12, 14, 16}
	for i, hour := range hours {
		points[i].Masuk = incoming[hour]
		points[i].Keluar = outgoing[hour]
	}
	return points, nil
}

func buildLocationMetrics(locations []model.ParkingLocation) ([]model.LocationMetric, []model.LocationMetric, []model.ActionItem) {
	revenue := make([]model.LocationMetric, 0, len(locations))
	occupancy := make([]model.LocationMetric, 0, len(locations))
	maxTraffic := int64(1)
	for _, location := range locations {
		traffic := location.Motorcycles + location.Cars
		if traffic > maxTraffic {
			maxTraffic = traffic
		}
	}
	for _, location := range locations {
		revenue = append(revenue, model.LocationMetric{
			Name:      location.Name,
			Value:     (location.TariffMotor * location.Motorcycles) + (location.TariffMobil * location.Cars),
			Secondary: "Akumulasi estimasi per lokasi",
			Tone:      "blue",
		})
		occupancyPercent := ((location.Motorcycles + location.Cars) * 100) / maxTraffic
		tone := "green"
		if occupancyPercent >= 80 {
			tone = "orange"
		} else if occupancyPercent >= 55 {
			tone = "blue"
		} else if occupancyPercent >= 30 {
			tone = "gold"
		}
		occupancy = append(occupancy, model.LocationMetric{
			Name:      location.Name,
			Value:     occupancyPercent,
			Secondary: location.OccupancyLabel,
			Tone:      tone,
		})
	}
	sort.Slice(revenue, func(i, j int) bool { return revenue[i].Value > revenue[j].Value })
	sort.Slice(occupancy, func(i, j int) bool { return occupancy[i].Value > occupancy[j].Value })
	for i := range revenue {
		revenue[i].Tone = toneFromRank(i)
	}
	for i := range occupancy {
		occupancy[i].Tone = toneFromRank(i)
	}
	priority := make([]model.ActionItem, 0, minInt(3, len(occupancy)))
	for i := 0; i < len(occupancy) && i < 3; i++ {
		item := occupancy[i]
		locationID := ""
		for _, location := range locations {
			if location.Name == item.Name {
				locationID = location.ID
				break
			}
		}
		priority = append(priority, model.ActionItem{
			LocationID:     locationID,
			Location:       item.Name,
			Issue:          fmt.Sprintf("Occupancy indikator %d%%", item.Value),
			Recommendation: "Buka monitoring lokasi, cek shift global, dan siapkan rotasi petugas bila diperlukan.",
			Href:           fmt.Sprintf("/admin/monitoring?location=%s", locationID),
			Tone:           item.Tone,
		})
	}
	return revenue, occupancy, priority
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (r *MySQLRepository) buildComparisonMetrics(ctx context.Context, now time.Time) ([]model.ComparisonMetric, error) {
	todayStart, todayEnd := todayWindow(now)
	yesterdayStart, yesterdayEnd := yesterdayWindow(now)
	var revenueToday, revenueYesterday, txToday, txYesterday int64
	_ = r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(final_amount), 0) FROM financial_parking_transaction WHERE paid_at >= ? AND paid_at < ?`, todayStart, todayEnd).Scan(&revenueToday)
	_ = r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(final_amount), 0) FROM financial_parking_transaction WHERE paid_at >= ? AND paid_at < ?`, yesterdayStart, yesterdayEnd).Scan(&revenueYesterday)
	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM financial_parking_transaction WHERE paid_at >= ? AND paid_at < ?`, todayStart, todayEnd).Scan(&txToday)
	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM financial_parking_transaction WHERE paid_at >= ? AND paid_at < ?`, yesterdayStart, yesterdayEnd).Scan(&txYesterday)
	var activeToday, activeYesterday int64
	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM parking_session WHERE started_at >= ? AND started_at < ?`, todayStart, todayEnd).Scan(&activeToday)
	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM parking_session WHERE started_at >= ? AND started_at < ?`, yesterdayStart, yesterdayEnd).Scan(&activeYesterday)
	return []model.ComparisonMetric{
		{Label: "Pendapatan", Today: revenueToday, Yesterday: revenueYesterday, Unit: "currency"},
		{Label: "Kendaraan Masuk", Today: txToday, Yesterday: txYesterday, Unit: "count"},
		{Label: "Occupancy Puncak", Today: activeToday * 10, Yesterday: activeYesterday * 10, Unit: "percent"},
	}, nil
}

func (r *MySQLRepository) buildHeatmap(ctx context.Context, now time.Time) ([]model.HeatmapPoint, error) {
	start := now.AddDate(0, 0, -6)
	query := `
		SELECT DATE(paid_at) AS d, HOUR(paid_at) AS h, COUNT(*)
		FROM financial_parking_transaction
		WHERE paid_at >= ?
		GROUP BY DATE(paid_at), HOUR(paid_at)`
	rows, err := r.db.QueryContext(ctx, query, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	counts := map[string]int64{}
	for rows.Next() {
		var d time.Time
		var hour int
		var count int64
		if err := rows.Scan(&d, &hour, &count); err != nil {
			return nil, err
		}
		counts[fmt.Sprintf("%s-%02d", d.Format("2006-01-02"), hour)] = count
	}
	points := make([]model.HeatmapPoint, 0, 20)
	hours := []string{"08", "10", "12", "14"}
	for i := 4; i >= 0; i-- {
		day := now.AddDate(0, 0, -i)
		for _, hour := range hours {
			key := fmt.Sprintf("%s-%s", day.Format("2006-01-02"), hour)
			points = append(points, model.HeatmapPoint{
				Day:   dayNameID(day),
				Hour:  hour,
				Value: counts[key],
			})
		}
	}
	return points, nil
}

func (r *MySQLRepository) GetMonitoringOverview(ctx context.Context) (model.MonitoringOverview, error) {
	locations, err := r.buildParkingLocations(ctx)
	if err != nil {
		return model.MonitoringOverview{}, err
	}
	maxTransactions := int64(1)
	for _, location := range locations {
		if location.Transactions > maxTransactions {
			maxTransactions = location.Transactions
		}
	}
	parkingLocations := make([]model.ParkingLocation, 0, len(locations))
	monitoringZones := make([]model.RowItem, 0, len(locations))
	zoneSet := map[string]struct{}{}
	for _, location := range locations {
		item := toParkingLocation(location, maxTransactions)
		parkingLocations = append(parkingLocations, item)
		zoneSet[item.Zone] = struct{}{}
		monitoringZones = append(monitoringZones, model.RowItem{
			Primary:    item.Name,
			Secondary:  item.Zone,
			Status:     item.OccupancyLabel,
			StatusTone: mapOccupancyTone(item.OccupancyLabel),
			ValueA:     fmt.Sprintf("%d", item.Motorcycles),
			ValueB:     fmt.Sprintf("%d", item.Cars),
			Location:   fmt.Sprintf("%d", item.Motorcycles+item.Cars),
		})
	}
	zones := []string{"Semua Zona"}
	for zone := range zoneSet {
		zones = append(zones, zone)
	}
	sort.Strings(zones[1:])
	officers, err := r.listOfficerOptions(ctx, locations)
	if err != nil {
		return model.MonitoringOverview{}, err
	}
	officerFilters := []string{"Semua Petugas"}
	officerFilterSet := map[string]struct{}{}
	for _, officer := range officers {
		if strings.TrimSpace(officer.DefaultShiftStart) == "" || strings.TrimSpace(officer.DefaultShiftEnd) == "" {
			continue
		}
		label := fmt.Sprintf("%s - %s", officer.DefaultShiftStart, officer.DefaultShiftEnd)
		if _, exists := officerFilterSet[label]; exists {
			continue
		}
		officerFilterSet[label] = struct{}{}
		officerFilters = append(officerFilters, label)
	}
	return model.MonitoringOverview{
		TopFilters: model.TopFilters{
			Zones:    zones,
			Dates:    time.Now().Format("02 Jan 2006"),
			Officers: officerFilters,
		},
		MonitoringZones:       monitoringZones,
		ParkingLocations:      parkingLocations,
		ParkingOfficerOptions: officers,
	}, nil
}

func mapOccupancyTone(label string) string {
	switch label {
	case "Zona Padat":
		return "gold"
	case "Ramai":
		return "blue"
	case "Normal":
		return "green"
	default:
		return "gray"
	}
}

func (r *MySQLRepository) GetOfficerOverview(ctx context.Context) (model.OfficerOverview, error) {
	monitoring, err := r.GetMonitoringOverview(ctx)
	if err != nil {
		return model.OfficerOverview{}, err
	}
	active, onShift, reserve := int64(0), int64(0), int64(0)
	for _, officer := range monitoring.ParkingOfficerOptions {
		if officer.Status == "Aktif" {
			active++
		}
		if officer.Availability == "Bertugas di Lokasi Lain" {
			onShift++
		}
		if officer.Availability == "Cadangan" {
			reserve++
		}
	}
	stats := []model.StatCard{
		{Label: "Petugas Aktif", Value: fmt.Sprintf("%d", active), Icon: "PT", Tone: "blue"},
		{Label: "Sedang Bertugas", Value: fmt.Sprintf("%d Petugas", onShift), Icon: "KD", Tone: "orange"},
		{Label: "Cadangan", Value: fmt.Sprintf("%d", reserve), Icon: "CK", Tone: "cyan"},
	}
	return model.OfficerOverview{
		OfficerStats:          stats,
		ParkingOfficerOptions: monitoring.ParkingOfficerOptions,
		ParkingLocations:      monitoring.ParkingLocations,
	}, nil
}

func (r *MySQLRepository) GetTransactionsOverview(ctx context.Context) (model.TransactionsOverview, error) {
	now := time.Now()
	todayStart, todayEnd := todayWindow(now)
	var totalTransactions, totalRevenue, motorTransactions, carTransactions int64
	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*), COALESCE(SUM(final_amount),0) FROM financial_parking_transaction WHERE paid_at >= ? AND paid_at < ?`, todayStart, todayEnd).Scan(&totalTransactions, &totalRevenue)
	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM financial_parking_transaction fpt JOIN vehicle_type vt ON vt.id = fpt.vehicle_type_id WHERE fpt.paid_at >= ? AND fpt.paid_at < ? AND LOWER(vt.vehicle_type_name) LIKE '%motor%'`, todayStart, todayEnd).Scan(&motorTransactions)
	_ = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM financial_parking_transaction fpt JOIN vehicle_type vt ON vt.id = fpt.vehicle_type_id WHERE fpt.paid_at >= ? AND fpt.paid_at < ? AND (LOWER(vt.vehicle_type_name) LIKE '%mobil%' OR LOWER(vt.vehicle_type_name) LIKE '%car%')`, todayStart, todayEnd).Scan(&carTransactions)
	stats := []model.StatCard{
		{Label: "Total Transaksi", Value: fmt.Sprintf("%d", totalTransactions), Icon: "TX", Tone: "orange"},
		{Label: "Total Pendapatan", Value: formatIDR(totalRevenue), Icon: "RP", Tone: "blue"},
		{Label: "Transaksi Motor", Value: fmt.Sprintf("%d", motorTransactions), Icon: "MT", Tone: "cyan"},
		{Label: "Transaksi Mobil", Value: fmt.Sprintf("%d", carTransactions), Icon: "MB", Tone: "blue"},
	}

	query := `
		SELECT
			fpt.id,
			COALESCE(fpt.successful_payment_event_id, 0),
			fpt.location_id,
			COALESCE(fpt.plate_number, ''),
			COALESCE(pz.zone_name, ''),
			COALESCE(su.full_name, ''),
			COALESCE(vt.vehicle_type_name, ''),
			DATE_FORMAT(fpt.paid_at, '%d %M %Y'),
			COALESCE(fpt.final_amount, 0)
		FROM financial_parking_transaction fpt
		LEFT JOIN parking_location pl ON pl.id = fpt.location_id
		LEFT JOIN parking_zone pz ON pz.id = pl.zone_id
		LEFT JOIN system_user su ON su.id = fpt.officer_user_id
		LEFT JOIN vehicle_type vt ON vt.id = fpt.vehicle_type_id
		WHERE fpt.paid_at IS NOT NULL
		ORDER BY fpt.paid_at DESC, fpt.id DESC
		LIMIT 10`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return model.TransactionsOverview{}, err
	}
	defer rows.Close()
	transactionRows := make([]model.RowItem, 0, 10)
	for rows.Next() {
		var transactionID, paymentEventID, locationID int64
		var plate, zone, officer, vehicleType, paidDate string
		var amount int64
		if err := rows.Scan(&transactionID, &paymentEventID, &locationID, &plate, &zone, &officer, &vehicleType, &paidDate, &amount); err != nil {
			return model.TransactionsOverview{}, err
		}
		transactionRows = append(transactionRows, model.RowItem{
			TransactionID:  fmt.Sprintf("%d", transactionID),
			PaymentEventID: formatOptionalID(paymentEventID),
			LocationID:     fmt.Sprintf("%d", locationID),
			Primary:        plate,
			Secondary:      zone,
			Status:         officer,
			StatusTone:     "blue",
			ValueA:         vehicleType,
			ValueB:         paidDate,
			Price:          formatIDR(amount),
		})
	}

	var cashAmount, qrisAmount int64
	_ = r.db.QueryRowContext(ctx, `
		SELECT
			COALESCE(SUM(CASE WHEN LOWER(pm.payment_method_code) = 'cash' THEN fpt.final_amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN LOWER(pm.payment_method_code) = 'qris' THEN fpt.final_amount ELSE 0 END), 0)
		FROM financial_parking_transaction fpt
		LEFT JOIN payment_method pm ON pm.id = fpt.payment_method_id
		WHERE fpt.paid_at >= ? AND fpt.paid_at < ?`, todayStart, todayEnd).Scan(&cashAmount, &qrisAmount)
	paymentBreakdownItems := []model.PaymentBreakdownItem{
		{Label: "QRIS", Amount: formatIDR(qrisAmount), Share: percentageLabel(qrisAmount, totalRevenue), Tone: "blue"},
		{Label: "Cash", Amount: formatIDR(cashAmount), Share: percentageLabel(cashAmount, totalRevenue), Tone: "gold"},
	}

	transactionIssueItems := []model.TransactionIssueItem{
		{Title: "Sinkronisasi QRIS", Detail: "Pantau payment event settled yang belum masuk reconciliation batch.", Action: "Audit callback dan batch rekonsiliasi", Tone: "orange"},
		{Title: "Tarif lokasi turunan", Detail: "Tarif lokasi saat ini masih diturunkan dari estimasi transaksi karena tabel tarif v2 belum dipisah.", Action: "Tambahkan domain tarif versi produksi", Tone: "blue"},
		{Title: "Validasi closing harian", Detail: "Pastikan variance closing tetap nol sebelum hari ditutup.", Action: "Review closing batch dan settlement", Tone: "gold"},
	}

	exportQueueItems := []model.ExportQueueItem{
		{Label: "Laporan Harian", Status: "Siap diunduh", Note: now.Format("02 Jan 2006")},
		{Label: "Rekap Per Petugas", Status: "Butuh filter", Note: "Gunakan filter lokasi dan periode"},
		{Label: "Settlement QRIS", Status: "Menunggu sinkron", Note: "Cek batch rekonsiliasi terbaru"},
	}

	return model.TransactionsOverview{
		TransactionStats:      stats,
		TransactionRows:       transactionRows,
		PaymentBreakdownItems: paymentBreakdownItems,
		TransactionIssueItems: transactionIssueItems,
		ExportQueueItems:      exportQueueItems,
	}, nil
}

func percentageLabel(value, total int64) string {
	if total <= 0 {
		return "0%"
	}
	return fmt.Sprintf("%d%%", (value*100)/total)
}

func (r *MySQLRepository) GetSettingsOverview(ctx context.Context) (model.SettingsOverview, error) {
	if err := r.ensureAdminTables(ctx); err != nil {
		return model.SettingsOverview{}, err
	}
	alertRuleItems, err := r.listSettingsAlertRules(ctx)
	if err != nil {
		return model.SettingsOverview{}, err
	}
	shiftTemplates, err := r.listSettingsShiftTemplates(ctx)
	if err != nil {
		return model.SettingsOverview{}, err
	}

	roleRows, err := r.db.QueryContext(ctx, `SELECT role_code, role_name FROM system_role ORDER BY id`)
	if err != nil {
		return model.SettingsOverview{}, err
	}
	defer roleRows.Close()
	adminRoleItems := make([]model.AdminRoleItem, 0)
	for roleRows.Next() {
		var code, name string
		if err := roleRows.Scan(&code, &name); err != nil {
			return model.SettingsOverview{}, err
		}
		if !strings.HasPrefix(code, "admin") && code != "sco" {
			continue
		}
		adminRoleItems = append(adminRoleItems, model.AdminRoleItem{
			Role:   name,
			Access: describeRoleAccess(code),
			Owner:  describeRoleOwner(code),
		})
	}

	notificationItems, err := r.listSettingsNotifications(ctx)
	if err != nil {
		return model.SettingsOverview{}, err
	}
	paymentMethodItems, err := r.listSettingsPaymentMethods(ctx)
	if err != nil {
		return model.SettingsOverview{}, err
	}
	defaultTariffItems, err := r.listSettingsDefaultTariffs(ctx)
	if err != nil {
		return model.SettingsOverview{}, err
	}

	return model.SettingsOverview{
		AlertRuleItems:        alertRuleItems,
		DefaultShiftTemplates: shiftTemplates,
		DefaultTariffItems:    defaultTariffItems,
		AdminRoleItems:        adminRoleItems,
		NotificationItems:     notificationItems,
		PaymentMethodItems:    paymentMethodItems,
	}, nil
}

func (r *MySQLRepository) UpdateSettingsOverview(ctx context.Context, adminUserID int64, req model.UpdateSettingsOverviewRequest) (model.SettingsOverview, error) {
	if err := r.ensureAdminTables(ctx); err != nil {
		return model.SettingsOverview{}, err
	}
	if len(req.AlertRuleItems) == 0 {
		return model.SettingsOverview{}, errors.New("alert rule wajib diisi")
	}
	if len(req.DefaultShiftTemplates) == 0 {
		return model.SettingsOverview{}, errors.New("template shift wajib diisi")
	}
	if len(req.DefaultTariffItems) == 0 {
		return model.SettingsOverview{}, errors.New("tarif default wajib diisi")
	}
	if len(req.NotificationItems) == 0 {
		return model.SettingsOverview{}, errors.New("notifikasi wajib diisi")
	}
	if len(req.PaymentMethodItems) == 0 {
		return model.SettingsOverview{}, errors.New("metode pembayaran wajib diisi")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.SettingsOverview{}, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM admin_settings_alert_rule`); err != nil {
		return model.SettingsOverview{}, err
	}
	for index, item := range req.AlertRuleItems {
		if strings.TrimSpace(item.Title) == "" || strings.TrimSpace(item.Threshold) == "" || strings.TrimSpace(item.Source) == "" || strings.TrimSpace(item.PIC) == "" {
			return model.SettingsOverview{}, errors.New("setiap rule alert wajib memiliki title, threshold, source, dan PIC")
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO admin_settings_alert_rule (
				rule_title, threshold_text, source_text, pic_text, sort_order, is_active, updated_by_user_id
			) VALUES (?, ?, ?, ?, ?, 1, ?)`,
			strings.TrimSpace(item.Title),
			strings.TrimSpace(item.Threshold),
			strings.TrimSpace(item.Source),
			strings.TrimSpace(item.PIC),
			index,
			adminUserID,
		); err != nil {
			return model.SettingsOverview{}, err
		}
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM admin_settings_shift_template`); err != nil {
		return model.SettingsOverview{}, err
	}
	for index, item := range req.DefaultShiftTemplates {
		start, end := splitHoursRange(item.Hours)
		if strings.TrimSpace(item.Label) == "" || start == "" || end == "" || strings.TrimSpace(item.UseCase) == "" {
			return model.SettingsOverview{}, errors.New("setiap template shift wajib memiliki label, jam, dan use case")
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO admin_settings_shift_template (
				shift_label, start_time, end_time, use_case, sort_order, is_active, updated_by_user_id
			) VALUES (?, ?, ?, ?, ?, 1, ?)`,
			strings.TrimSpace(item.Label),
			start,
			end,
			strings.TrimSpace(item.UseCase),
			index,
			adminUserID,
		); err != nil {
			return model.SettingsOverview{}, err
		}
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM admin_settings_default_tariff`); err != nil {
		return model.SettingsOverview{}, err
	}
	for _, item := range req.DefaultTariffItems {
		vehicleType := strings.TrimSpace(item.VehicleType)
		if vehicleType == "" {
			return model.SettingsOverview{}, errors.New("jenis kendaraan tarif default wajib diisi")
		}
		if item.FirstHour < 0 || item.NextHour < 0 || item.MaxRate < 0 {
			return model.SettingsOverview{}, errors.New("nominal tarif default tidak boleh negatif")
		}
		vehicleCode := strings.ToLower(strings.ReplaceAll(vehicleType, " ", "_"))
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO admin_settings_default_tariff (
				vehicle_type_code, vehicle_type_label, first_hour_amount, next_hour_amount, max_rate_amount, updated_by_user_id
			) VALUES (?, ?, ?, ?, ?, ?)`,
			vehicleCode,
			vehicleType,
			item.FirstHour,
			item.NextHour,
			item.MaxRate,
			adminUserID,
		); err != nil {
			return model.SettingsOverview{}, err
		}
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM admin_settings_notification`); err != nil {
		return model.SettingsOverview{}, err
	}
	for index, item := range req.NotificationItems {
		if strings.TrimSpace(item.Channel) == "" || strings.TrimSpace(item.Trigger) == "" || strings.TrimSpace(item.Response) == "" {
			return model.SettingsOverview{}, errors.New("setiap notifikasi wajib memiliki channel, trigger, dan response")
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO admin_settings_notification (
				channel_name, trigger_text, response_text, sort_order, is_active, updated_by_user_id
			) VALUES (?, ?, ?, ?, 1, ?)`,
			strings.TrimSpace(item.Channel),
			strings.TrimSpace(item.Trigger),
			strings.TrimSpace(item.Response),
			index,
			adminUserID,
		); err != nil {
			return model.SettingsOverview{}, err
		}
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM admin_settings_payment_method`); err != nil {
		return model.SettingsOverview{}, err
	}
	for _, item := range req.PaymentMethodItems {
		if strings.TrimSpace(item.Label) == "" {
			return model.SettingsOverview{}, errors.New("label metode pembayaran wajib diisi")
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO admin_settings_payment_method (
				method_code, method_label, icon_code, is_enabled, updated_by_user_id
			) VALUES (?, ?, ?, ?, ?)`,
			paymentMethodCode(item.Label),
			strings.TrimSpace(item.Label),
			strings.TrimSpace(item.Icon),
			item.Enabled,
			adminUserID,
		); err != nil {
			return model.SettingsOverview{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return model.SettingsOverview{}, err
	}
	return r.GetSettingsOverview(ctx)
}

func (r *MySQLRepository) listSettingsAlertRules(ctx context.Context) ([]model.AlertRuleItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT rule_title, threshold_text, source_text, pic_text
		FROM admin_settings_alert_rule
		WHERE is_active = 1
		ORDER BY sort_order, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.AlertRuleItem, 0)
	for rows.Next() {
		var item model.AlertRuleItem
		if err := rows.Scan(&item.Title, &item.Threshold, &item.Source, &item.PIC); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(items) > 0 {
		return items, nil
	}

	alertRows, err := r.db.QueryContext(ctx, `
		SELECT
			COALESCE(ar.rule_title, ''),
			COALESCE(ar.threshold_config_json, '{}'),
			COALESCE(ar.source_type, ''),
			COALESCE(su.full_name, '')
		FROM admin_alert_rule ar
		LEFT JOIN system_user su ON su.id = ar.created_by_user_id
		WHERE ar.is_active = 1
		ORDER BY ar.id`)
	if err != nil {
		return nil, err
	}
	defer alertRows.Close()
	for alertRows.Next() {
		var title, thresholdJSON, sourceType, creator string
		if err := alertRows.Scan(&title, &thresholdJSON, &sourceType, &creator); err != nil {
			return nil, err
		}
		items = append(items, model.AlertRuleItem{
			Title:     title,
			Threshold: summarizeThreshold(thresholdJSON),
			Source:    sourceType,
			PIC:       creator,
		})
	}
	return items, alertRows.Err()
}

func (r *MySQLRepository) listSettingsShiftTemplates(ctx context.Context) ([]model.ShiftTemplateItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT shift_label, start_time, end_time, use_case
		FROM admin_settings_shift_template
		WHERE is_active = 1
		ORDER BY sort_order, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.ShiftTemplateItem, 0)
	for rows.Next() {
		var label, start, end, useCase string
		if err := rows.Scan(&label, &start, &end, &useCase); err != nil {
			return nil, err
		}
		items = append(items, model.ShiftTemplateItem{
			Label:   label,
			Hours:   fmt.Sprintf("%s - %s", start, end),
			UseCase: useCase,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(items) > 0 {
		return items, nil
	}

	shiftRows, err := r.db.QueryContext(ctx, `
		SELECT
			COALESCE(shift_name, ''),
			DATE_FORMAT(start_time, '%H:%i'),
			DATE_FORMAT(end_time, '%H:%i')
		FROM parking_shift_template
		WHERE is_active = 1
		GROUP BY shift_name, start_time, end_time
		ORDER BY start_time, shift_name`)
	if err != nil {
		return nil, err
	}
	defer shiftRows.Close()
	for shiftRows.Next() {
		var label, start, end string
		if err := shiftRows.Scan(&label, &start, &end); err != nil {
			return nil, err
		}
		items = append(items, model.ShiftTemplateItem{
			Label:   label,
			Hours:   fmt.Sprintf("%s - %s", start, end),
			UseCase: "Template operasional lokasi aktif",
		})
	}
	return items, shiftRows.Err()
}

func (r *MySQLRepository) listSettingsNotifications(ctx context.Context) ([]model.NotificationItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT channel_name, trigger_text, response_text
		FROM admin_settings_notification
		WHERE is_active = 1
		ORDER BY sort_order, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.NotificationItem, 0)
	for rows.Next() {
		var item model.NotificationItem
		if err := rows.Scan(&item.Channel, &item.Trigger, &item.Response); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *MySQLRepository) listSettingsPaymentMethods(ctx context.Context) ([]model.PaymentMethodItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT method_label, is_enabled, icon_code
		FROM admin_settings_payment_method
		ORDER BY method_code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.PaymentMethodItem, 0)
	for rows.Next() {
		var item model.PaymentMethodItem
		if err := rows.Scan(&item.Label, &item.Enabled, &item.Icon); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(items) > 0 {
		return items, nil
	}
	methodRows, err := r.db.QueryContext(ctx, `
		SELECT payment_method_name, payment_method_code
		FROM payment_method
		ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer methodRows.Close()

	for methodRows.Next() {
		var label, code string
		if err := methodRows.Scan(&label, &code); err != nil {
			return nil, err
		}
		icon := strings.ToUpper(strings.TrimSpace(code))
		if len(icon) > 2 {
			icon = icon[:2]
		}
		items = append(items, model.PaymentMethodItem{
			Label:   label,
			Enabled: true,
			Icon:    icon,
		})
	}
	if err := methodRows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *MySQLRepository) listSettingsDefaultTariffs(ctx context.Context) ([]model.DefaultTariffItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT vehicle_type_label, first_hour_amount, next_hour_amount, max_rate_amount
		FROM admin_settings_default_tariff
		ORDER BY vehicle_type_code`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.DefaultTariffItem, 0)
	for rows.Next() {
		var item model.DefaultTariffItem
		if err := rows.Scan(&item.VehicleType, &item.FirstHour, &item.NextHour, &item.MaxRate); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(items) > 0 {
		return items, nil
	}
	vehicleRows, err := r.db.QueryContext(ctx, `
		SELECT vehicle_type_name
		FROM vehicle_type
		ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer vehicleRows.Close()

	for vehicleRows.Next() {
		var name string
		if err := vehicleRows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, model.DefaultTariffItem{
			VehicleType: name,
			FirstHour:   0,
			NextHour:    0,
			MaxRate:     0,
		})
	}
	if err := vehicleRows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func summarizeThreshold(raw string) string {
	payload := map[string]any{}
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return raw
	}
	parts := make([]string, 0, len(payload))
	keys := make([]string, 0, len(payload))
	for key := range payload {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", key, payload[key]))
	}
	return strings.Join(parts, " · ")
}

func describeRoleAccess(code string) string {
	switch code {
	case "admin_super":
		return "Semua modul + rule sistem"
	case "admin_finance":
		return "Transaksi, refund, settlement, closing"
	case "admin_ops":
		return "Monitoring, mutasi petugas, alert"
	case "sco":
		return "Pengawasan lapangan dan tindak lanjut alert"
	default:
		return "Akses operasional"
	}
}

func describeRoleOwner(code string) string {
	switch code {
	case "admin_super":
		return "Admin utama / pengelola sistem"
	case "admin_finance":
		return "Tim keuangan"
	case "admin_ops":
		return "Tim operasional"
	case "sco":
		return "Supervisor lapangan"
	default:
		return "Admin"
	}
}

func (r *MySQLRepository) findLocationByStringID(ctx context.Context, locationID string) (locationAggregate, error) {
	locations, err := r.buildParkingLocations(ctx)
	if err != nil {
		return locationAggregate{}, err
	}
	for _, location := range locations {
		if location.Slug == strings.TrimSpace(locationID) {
			return location, nil
		}
	}
	return locationAggregate{}, sql.ErrNoRows
}

func (r *MySQLRepository) findOfficerByStringID(ctx context.Context, officerID string) (model.ParkingOfficerOption, error) {
	locations, err := r.buildParkingLocations(ctx)
	if err != nil {
		return model.ParkingOfficerOption{}, err
	}
	officers, err := r.listOfficerOptions(ctx, locations)
	if err != nil {
		return model.ParkingOfficerOption{}, err
	}
	for _, officer := range officers {
		if officer.ID == strings.TrimSpace(officerID) {
			return officer, nil
		}
	}
	return model.ParkingOfficerOption{}, sql.ErrNoRows
}

func (r *MySQLRepository) UpdateLocationSettings(ctx context.Context, adminUserID int64, locationID string, req model.UpdateLocationSettingsRequest) (model.ParkingLocation, error) {
	if err := r.ensureAdminTables(ctx); err != nil {
		return model.ParkingLocation{}, err
	}
	locationPK, err := parseIDString(locationID)
	if err != nil {
		return model.ParkingLocation{}, err
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO admin_location_setting (location_id, tariff_motor_amount, tariff_car_amount, operational_note, updated_by_user_id)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			tariff_motor_amount = VALUES(tariff_motor_amount),
			tariff_car_amount = VALUES(tariff_car_amount),
			operational_note = VALUES(operational_note),
			updated_by_user_id = VALUES(updated_by_user_id),
			updated_at = CURRENT_TIMESTAMP`,
		locationPK,
		req.TariffMotor,
		req.TariffMobil,
		strings.TrimSpace(req.DismissalReason),
		adminUserID,
	)
	if err != nil {
		return model.ParkingLocation{}, err
	}
	location, err := r.findLocationByStringID(ctx, locationID)
	if err != nil {
		return model.ParkingLocation{}, err
	}
	maxTraffic := maxInt64(1, location.Motorcycles+location.Cars)
	return toParkingLocation(location, maxTraffic), nil
}

func (r *MySQLRepository) SaveLocationShiftTemplates(ctx context.Context, adminUserID int64, locationID string, shiftTemplates []model.ParkingShiftTemplate) (model.ParkingLocation, error) {
	locationPK, err := parseIDString(locationID)
	if err != nil {
		return model.ParkingLocation{}, err
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.ParkingLocation{}, err
	}
	defer tx.Rollback()

	existingRows, err := tx.QueryContext(ctx, `SELECT id FROM parking_shift_template WHERE location_id = ?`, locationPK)
	if err != nil {
		return model.ParkingLocation{}, err
	}
	defer existingRows.Close()
	existing := map[int64]struct{}{}
	for existingRows.Next() {
		var id int64
		if err := existingRows.Scan(&id); err != nil {
			return model.ParkingLocation{}, err
		}
		existing[id] = struct{}{}
	}
	if err := existingRows.Err(); err != nil {
		return model.ParkingLocation{}, err
	}

	used := map[int64]struct{}{}
	for _, shift := range shiftTemplates {
		label := strings.TrimSpace(shift.Label)
		start := strings.TrimSpace(shift.Start)
		end := strings.TrimSpace(shift.End)
		if label == "" || start == "" || end == "" {
			return model.ParkingLocation{}, errors.New("shift template wajib memiliki nama, waktu mulai, dan waktu selesai")
		}
		if id, parseErr := strconv.ParseInt(strings.TrimSpace(shift.ID), 10, 64); parseErr == nil && id > 0 {
			used[id] = struct{}{}
			if _, ok := existing[id]; ok {
				if _, err := tx.ExecContext(ctx, `
					UPDATE parking_shift_template
					SET shift_name = ?, shift_code = ?, start_time = ?, end_time = ?, is_active = 1
					WHERE id = ? AND location_id = ?`,
					label,
					shiftCodeFromTemplate(shift),
					start,
					end,
					id,
					locationPK,
				); err != nil {
					return model.ParkingLocation{}, err
				}
				continue
			}
		}
		result, err := tx.ExecContext(ctx, `
			INSERT INTO parking_shift_template (location_id, shift_code, shift_name, start_time, end_time, is_active, created_at)
			VALUES (?, ?, ?, ?, ?, 1, CURRENT_TIMESTAMP)`,
			locationPK,
			shiftCodeFromTemplate(shift),
			label,
			start,
			end,
		)
		if err != nil {
			return model.ParkingLocation{}, err
		}
		insertID, _ := result.LastInsertId()
		if insertID > 0 {
			used[insertID] = struct{}{}
		}
	}

	for existingID := range existing {
		if _, ok := used[existingID]; ok {
			continue
		}
		if _, err := tx.ExecContext(ctx, `UPDATE parking_shift_template SET is_active = 0 WHERE id = ? AND location_id = ?`, existingID, locationPK); err != nil {
			return model.ParkingLocation{}, err
		}
	}

	if _, err := tx.ExecContext(ctx, `UPDATE officer_assignment_current SET assigned_by_user_id = ?, updated_at = CURRENT_TIMESTAMP WHERE location_id = ? AND effective_to IS NULL`, adminUserID, locationPK); err != nil {
		return model.ParkingLocation{}, err
	}

	if err := tx.Commit(); err != nil {
		return model.ParkingLocation{}, err
	}
	location, err := r.findLocationByStringID(ctx, locationID)
	if err != nil {
		return model.ParkingLocation{}, err
	}
	maxTraffic := maxInt64(1, location.Motorcycles+location.Cars)
	return toParkingLocation(location, maxTraffic), nil
}

func (r *MySQLRepository) UpdateOfficerStatus(ctx context.Context, adminUserID int64, officerID string, status string) (model.ParkingOfficerOption, error) {
	officerPK, err := parseIDString(officerID)
	if err != nil {
		return model.ParkingOfficerOption{}, err
	}
	nextStatus := toOperationalStatus(status)
	var previous sql.NullString
	_ = r.db.QueryRowContext(ctx, `SELECT operational_status FROM officer_assignment_current WHERE officer_user_id = ? AND effective_to IS NULL LIMIT 1`, officerPK).Scan(&previous)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.ParkingOfficerOption{}, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `
		UPDATE officer_assignment_current
		SET operational_status = ?, assigned_by_user_id = ?, updated_at = CURRENT_TIMESTAMP
		WHERE officer_user_id = ? AND effective_to IS NULL`,
		nextStatus, adminUserID, officerPK,
	); err != nil {
		return model.ParkingOfficerOption{}, err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO officer_status_history (officer_user_id, old_operational_status, new_operational_status, change_reason, changed_by_user_id, changed_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`,
		officerPK, nullableStringValue(previous), nextStatus, "Status diperbarui dari dashboard admin", adminUserID,
	); err != nil {
		return model.ParkingOfficerOption{}, err
	}
	if err := tx.Commit(); err != nil {
		return model.ParkingOfficerOption{}, err
	}
	return r.findOfficerByStringID(ctx, officerID)
}

func (r *MySQLRepository) ApplyOfficerMutation(ctx context.Context, adminUserID int64, req model.ApplyOfficerMutationRequest) (model.ParkingOfficerOption, error) {
	officerPK, err := parseIDString(req.OfficerID)
	if err != nil {
		return model.ParkingOfficerOption{}, err
	}
	locationPK, err := parseIDString(req.TargetLocationID)
	if err != nil {
		return model.ParkingOfficerOption{}, err
	}
	shiftPK, err := parseIDString(req.TargetShiftID)
	if err != nil {
		return model.ParkingOfficerOption{}, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.ParkingOfficerOption{}, err
	}
	defer tx.Rollback()

	var currentAssignmentID int64
	var previousLocationID, previousShiftID, zoneID, areaID sql.NullInt64
	var previousStatus string
	err = tx.QueryRowContext(ctx, `
		SELECT id, location_id, shift_template_id, operational_status
		FROM officer_assignment_current
		WHERE officer_user_id = ? AND effective_to IS NULL
		LIMIT 1`, officerPK).Scan(&currentAssignmentID, &previousLocationID, &previousShiftID, &previousStatus)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return model.ParkingOfficerOption{}, err
	}
	if strings.TrimSpace(previousStatus) == "" {
		previousStatus = "off_duty"
	}

	if err := tx.QueryRowContext(ctx, `SELECT zone_id, area_id FROM parking_location WHERE id = ? LIMIT 1`, locationPK).Scan(&zoneID, &areaID); err != nil {
		return model.ParkingOfficerOption{}, err
	}

	if currentAssignmentID > 0 {
		if _, err := tx.ExecContext(ctx, `
			UPDATE officer_assignment_current
			SET location_id = ?, zone_id = ?, area_id = ?, shift_template_id = ?, assigned_by_user_id = ?, updated_at = CURRENT_TIMESTAMP
			WHERE id = ?`,
			locationPK, nullableInt64Value(zoneID), nullableInt64Value(areaID), shiftPK, adminUserID, currentAssignmentID,
		); err != nil {
			return model.ParkingOfficerOption{}, err
		}
	} else {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO officer_assignment_current (
				officer_user_id, location_id, zone_id, area_id, shift_template_id, operational_status, effective_from, effective_to, assigned_by_user_id
			) VALUES (?, ?, ?, ?, ?, 'off_duty', CURRENT_TIMESTAMP, NULL, ?)`,
			officerPK, locationPK, nullableInt64Value(zoneID), nullableInt64Value(areaID), shiftPK, adminUserID,
		); err != nil {
			return model.ParkingOfficerOption{}, err
		}
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO officer_assignment_history (
			officer_user_id,
			from_location_id,
			to_location_id,
			from_shift_template_id,
			to_shift_template_id,
			from_operational_status,
			to_operational_status,
			change_reason,
			changed_by_user_id,
			changed_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`,
		officerPK,
		nullableInt64Value(previousLocationID),
		locationPK,
		nullableInt64Value(previousShiftID),
		shiftPK,
		nullableStringOrNil(previousStatus),
		previousStatus,
		strings.TrimSpace(req.Note),
		adminUserID,
	); err != nil {
		return model.ParkingOfficerOption{}, err
	}

	if err := tx.Commit(); err != nil {
		return model.ParkingOfficerOption{}, err
	}
	return r.findOfficerByStringID(ctx, req.OfficerID)
}

func (r *MySQLRepository) CreateDisputeCase(ctx context.Context, adminUserID int64, req model.CreateDisputeCaseRequest) (model.DisputeCaseSummary, error) {
	if strings.TrimSpace(req.ReferenceEntityType) == "" || req.ReferenceEntityID <= 0 {
		return model.DisputeCaseSummary{}, errors.New("reference entity dispute wajib diisi")
	}
	if strings.TrimSpace(req.CaseType) == "" {
		return model.DisputeCaseSummary{}, errors.New("case type wajib diisi")
	}

	code := nextBusinessCode("DIS")
	now := time.Now()
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO financial_dispute_case (
			dispute_case_code,
			reference_entity_type,
			reference_entity_id,
			case_type,
			case_status,
			opened_by_user_id,
			assigned_to_user_id,
			opened_at,
			resolution_note
		) VALUES (?, ?, ?, ?, 'open', ?, ?, ?, NULL)`,
		code,
		strings.TrimSpace(req.ReferenceEntityType),
		req.ReferenceEntityID,
		strings.TrimSpace(req.CaseType),
		adminUserID,
		nullableZeroInt64(req.AssignedToUserID),
		now,
	)
	if err != nil {
		return model.DisputeCaseSummary{}, err
	}
	disputeID, err := result.LastInsertId()
	if err != nil {
		return model.DisputeCaseSummary{}, err
	}
	if _, err := r.db.ExecContext(ctx, `
		INSERT INTO financial_dispute_history (
			dispute_case_id, old_case_status, new_case_status, change_note, changed_by_user_id, changed_at
		) VALUES (?, NULL, 'open', ?, ?, ?)`,
		disputeID,
		strings.TrimSpace(req.ChangeNote),
		adminUserID,
		now,
	); err != nil {
		return model.DisputeCaseSummary{}, err
	}
	return r.fetchDisputeCaseSummary(ctx, disputeID)
}

func (r *MySQLRepository) UpdateDisputeCaseStatus(ctx context.Context, adminUserID int64, disputeID string, req model.UpdateDisputeCaseStatusRequest) (model.DisputeCaseSummary, error) {
	disputePK, err := parseIDString(disputeID)
	if err != nil {
		return model.DisputeCaseSummary{}, err
	}
	nextStatus := normalizeDisputeStatus(req.Status)
	if nextStatus == "" {
		return model.DisputeCaseSummary{}, errors.New("status dispute tidak valid")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.DisputeCaseSummary{}, err
	}
	defer tx.Rollback()

	var oldStatus string
	if err := tx.QueryRowContext(ctx, `SELECT case_status FROM financial_dispute_case WHERE id = ?`, disputePK).Scan(&oldStatus); err != nil {
		return model.DisputeCaseSummary{}, err
	}
	now := time.Now()
	resolvedAt := any(nil)
	resolutionNote := any(nil)
	if nextStatus == "resolved" || nextStatus == "closed" {
		resolvedAt = now
		if note := strings.TrimSpace(req.ChangeNote); note != "" {
			resolutionNote = note
		}
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE financial_dispute_case
		SET case_status = ?,
			assigned_to_user_id = COALESCE(?, assigned_to_user_id),
			resolved_at = COALESCE(?, resolved_at),
			resolution_note = COALESCE(?, resolution_note)
		WHERE id = ?`,
		nextStatus,
		nullableZeroInt64(req.AssignedToUserID),
		resolvedAt,
		resolutionNote,
		disputePK,
	); err != nil {
		return model.DisputeCaseSummary{}, err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO financial_dispute_history (
			dispute_case_id, old_case_status, new_case_status, change_note, changed_by_user_id, changed_at
		) VALUES (?, ?, ?, ?, ?, ?)`,
		disputePK,
		oldStatus,
		nextStatus,
		strings.TrimSpace(req.ChangeNote),
		adminUserID,
		now,
	); err != nil {
		return model.DisputeCaseSummary{}, err
	}
	if err := tx.Commit(); err != nil {
		return model.DisputeCaseSummary{}, err
	}
	return r.fetchDisputeCaseSummary(ctx, disputePK)
}

func (r *MySQLRepository) CreateRefundTransaction(ctx context.Context, adminUserID int64, req model.CreateRefundTransactionRequest) (model.RefundTransactionSummary, error) {
	if strings.TrimSpace(req.ReferenceEntityType) == "" || req.ReferenceEntityID <= 0 {
		return model.RefundTransactionSummary{}, errors.New("reference entity refund wajib diisi")
	}
	if req.RefundAmount <= 0 {
		return model.RefundTransactionSummary{}, errors.New("refund amount harus lebih besar dari nol")
	}
	if strings.TrimSpace(req.RefundReason) == "" {
		return model.RefundTransactionSummary{}, errors.New("refund reason wajib diisi")
	}

	code := nextBusinessCode("REF")
	now := time.Now()
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO financial_refund_transaction (
			refund_transaction_code,
			reference_entity_type,
			reference_entity_id,
			payment_event_id,
			wallet_id,
			refund_amount,
			currency_code,
			refund_reason,
			refund_status,
			requested_by_user_id,
			requested_at
		) VALUES (?, ?, ?, ?, ?, ?, 'IDR', ?, 'requested', ?, ?)`,
		code,
		strings.TrimSpace(req.ReferenceEntityType),
		req.ReferenceEntityID,
		nullableZeroInt64(req.PaymentEventID),
		nullableZeroInt64(req.WalletID),
		req.RefundAmount,
		strings.TrimSpace(req.RefundReason),
		adminUserID,
		now,
	)
	if err != nil {
		return model.RefundTransactionSummary{}, err
	}
	refundID, err := result.LastInsertId()
	if err != nil {
		return model.RefundTransactionSummary{}, err
	}
	if _, err := r.db.ExecContext(ctx, `
		INSERT INTO financial_refund_history (
			refund_transaction_id, old_refund_status, new_refund_status, status_reason, changed_by_user_id, changed_at
		) VALUES (?, NULL, 'requested', ?, ?, ?)`,
		refundID,
		strings.TrimSpace(req.RefundReason),
		adminUserID,
		now,
	); err != nil {
		return model.RefundTransactionSummary{}, err
	}
	return r.fetchRefundTransactionSummary(ctx, refundID)
}

func (r *MySQLRepository) UpdateRefundTransactionStatus(ctx context.Context, adminUserID int64, refundID string, req model.UpdateRefundStatusRequest) (model.RefundTransactionSummary, error) {
	refundPK, err := parseIDString(refundID)
	if err != nil {
		return model.RefundTransactionSummary{}, err
	}
	nextStatus := normalizeRefundStatus(req.Status)
	if nextStatus == "" {
		return model.RefundTransactionSummary{}, errors.New("status refund tidak valid")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.RefundTransactionSummary{}, err
	}
	defer tx.Rollback()

	var oldStatus, refType string
	var refID, refundAmount int64
	if err := tx.QueryRowContext(ctx, `
		SELECT refund_status, reference_entity_type, reference_entity_id, refund_amount
		FROM financial_refund_transaction
		WHERE id = ?`,
		refundPK,
	).Scan(&oldStatus, &refType, &refID, &refundAmount); err != nil {
		return model.RefundTransactionSummary{}, err
	}

	now := time.Now()
	approvedBy := any(nil)
	approvedAt := any(nil)
	processedAt := any(nil)
	if nextStatus == "approved" || nextStatus == "processed" || nextStatus == "settled" {
		approvedBy = adminUserID
		approvedAt = now
	}
	if nextStatus == "processed" || nextStatus == "settled" {
		processedAt = now
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE financial_refund_transaction
		SET refund_status = ?,
			approved_by_user_id = COALESCE(?, approved_by_user_id),
			approved_at = COALESCE(?, approved_at),
			processed_at = COALESCE(?, processed_at)
		WHERE id = ?`,
		nextStatus,
		approvedBy,
		approvedAt,
		processedAt,
		refundPK,
	); err != nil {
		return model.RefundTransactionSummary{}, err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO financial_refund_history (
			refund_transaction_id, old_refund_status, new_refund_status, status_reason, changed_by_user_id, changed_at
		) VALUES (?, ?, ?, ?, ?, ?)`,
		refundPK,
		oldStatus,
		nextStatus,
		strings.TrimSpace(req.StatusReason),
		adminUserID,
		now,
	); err != nil {
		return model.RefundTransactionSummary{}, err
	}

	if (nextStatus == "processed" || nextStatus == "settled") && refType == "financial_parking_transaction" {
		var txStatus string
		var finalAmount int64
		if err := tx.QueryRowContext(ctx, `SELECT transaction_status, final_amount FROM financial_parking_transaction WHERE id = ?`, refID).Scan(&txStatus, &finalAmount); err != nil {
			return model.RefundTransactionSummary{}, err
		}
		nextTxStatus := "refunded_partial"
		if refundAmount >= finalAmount {
			nextTxStatus = "refunded_full"
		}
		if _, err := tx.ExecContext(ctx, `
			UPDATE financial_parking_transaction
			SET transaction_status = ?
			WHERE id = ?`,
			nextTxStatus,
			refID,
		); err != nil {
			return model.RefundTransactionSummary{}, err
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO financial_parking_transaction_history (
				parking_transaction_id, old_transaction_status, new_transaction_status, status_reason, changed_by_user_id, changed_at
			) VALUES (?, ?, ?, ?, ?, ?)`,
			refID,
			txStatus,
			nextTxStatus,
			strings.TrimSpace(req.StatusReason),
			adminUserID,
			now,
		); err != nil {
			return model.RefundTransactionSummary{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return model.RefundTransactionSummary{}, err
	}
	return r.fetchRefundTransactionSummary(ctx, refundPK)
}

func (r *MySQLRepository) CreateClosingBatch(ctx context.Context, adminUserID int64, req model.CreateClosingBatchRequest) (model.ClosingBatchSummary, error) {
	if req.LocationID <= 0 {
		return model.ClosingBatchSummary{}, errors.New("location id closing wajib diisi")
	}
	closingDate, err := time.Parse("2006-01-02", strings.TrimSpace(req.ClosingDate))
	if err != nil {
		return model.ClosingBatchSummary{}, errors.New("closing date harus berformat YYYY-MM-DD")
	}

	var cashSales, cashlessSales int64
	if err := r.db.QueryRowContext(ctx, `
		SELECT
			COALESCE(SUM(CASE WHEN LOWER(pm.payment_method_code) = 'cash' THEN fpt.final_amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN LOWER(pm.payment_method_code) <> 'cash' THEN fpt.final_amount ELSE 0 END), 0)
		FROM financial_parking_transaction fpt
		JOIN payment_method pm ON pm.id = fpt.payment_method_id
		WHERE fpt.location_id = ?
		  AND DATE(COALESCE(fpt.paid_at, fpt.occurred_at)) = ?
		  AND fpt.transaction_status IN ('paid', 'refunded_partial', 'refunded_full')`,
		req.LocationID,
		closingDate.Format("2006-01-02"),
	).Scan(&cashSales, &cashlessSales); err != nil {
		return model.ClosingBatchSummary{}, err
	}

	var refundAmount int64
	if err := r.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(fr.refund_amount), 0)
		FROM financial_refund_transaction fr
		JOIN financial_parking_transaction fpt
		  ON fr.reference_entity_type = 'financial_parking_transaction'
		 AND fr.reference_entity_id = fpt.id
		WHERE fpt.location_id = ?
		  AND DATE(fr.requested_at) = ?
		  AND fr.refund_status IN ('approved', 'processed', 'settled')`,
		req.LocationID,
		closingDate.Format("2006-01-02"),
	).Scan(&refundAmount); err != nil {
		return model.ClosingBatchSummary{}, err
	}

	const openingBalance, topupAmount, adjustmentAmount int64 = 0, 0, 0
	expected := openingBalance + cashSales + cashlessSales + topupAmount + adjustmentAmount - refundAmount
	actual := req.ActualClosingBalanceAmount
	variance := actual - expected

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.ClosingBatchSummary{}, err
	}
	defer tx.Rollback()

	code := nextBusinessCode("CLOSE")
	now := time.Now()
	result, err := tx.ExecContext(ctx, `
		INSERT INTO location_daily_closing_batch (
			closing_batch_code,
			location_id,
			closing_date,
			opening_balance_amount,
			cash_sales_amount,
			cashless_sales_amount,
			topup_amount,
			refund_amount,
			adjustment_amount,
			expected_closing_balance_amount,
			actual_closing_balance_amount,
			variance_amount,
			closing_status,
			submitted_by_user_id,
			submitted_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'submitted', ?, ?)`,
		code,
		req.LocationID,
		closingDate.Format("2006-01-02"),
		openingBalance,
		cashSales,
		cashlessSales,
		topupAmount,
		refundAmount,
		adjustmentAmount,
		expected,
		actual,
		variance,
		adminUserID,
		now,
	)
	if err != nil {
		return model.ClosingBatchSummary{}, err
	}
	closingID, err := result.LastInsertId()
	if err != nil {
		return model.ClosingBatchSummary{}, err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO location_daily_closing_history (
			closing_batch_id, old_closing_status, new_closing_status, change_note, changed_by_user_id, changed_at
		) VALUES (?, NULL, 'submitted', ?, ?, ?)`,
		closingID,
		strings.TrimSpace(req.ChangeNote),
		adminUserID,
		now,
	); err != nil {
		return model.ClosingBatchSummary{}, err
	}

	type closingItemSeed struct {
		referenceType string
		referenceID   int64
		amount        int64
		direction     string
		note          string
	}
	transactionRows, err := tx.QueryContext(ctx, `
		SELECT id, final_amount
		FROM financial_parking_transaction
		WHERE location_id = ?
		  AND DATE(COALESCE(paid_at, occurred_at)) = ?
		  AND transaction_status IN ('paid', 'refunded_partial', 'refunded_full')`,
		req.LocationID,
		closingDate.Format("2006-01-02"),
	)
	if err != nil {
		return model.ClosingBatchSummary{}, err
	}
	items := make([]closingItemSeed, 0)
	for transactionRows.Next() {
		var id, amount int64
		if err := transactionRows.Scan(&id, &amount); err != nil {
			transactionRows.Close()
			return model.ClosingBatchSummary{}, err
		}
		items = append(items, closingItemSeed{
			referenceType: "financial_parking_transaction",
			referenceID:   id,
			amount:        amount,
			direction:     "in",
			note:          "Auto-linked paid transaction",
		})
	}
	if err := transactionRows.Close(); err != nil {
		return model.ClosingBatchSummary{}, err
	}

	refundRows, err := tx.QueryContext(ctx, `
		SELECT fr.id, fr.refund_amount
		FROM financial_refund_transaction fr
		JOIN financial_parking_transaction fpt
		  ON fr.reference_entity_type = 'financial_parking_transaction'
		 AND fr.reference_entity_id = fpt.id
		WHERE fpt.location_id = ?
		  AND DATE(fr.requested_at) = ?
		  AND fr.refund_status IN ('approved', 'processed', 'settled')`,
		req.LocationID,
		closingDate.Format("2006-01-02"),
	)
	if err != nil {
		return model.ClosingBatchSummary{}, err
	}
	for refundRows.Next() {
		var id, amount int64
		if err := refundRows.Scan(&id, &amount); err != nil {
			refundRows.Close()
			return model.ClosingBatchSummary{}, err
		}
		items = append(items, closingItemSeed{
			referenceType: "financial_refund_transaction",
			referenceID:   id,
			amount:        amount,
			direction:     "out",
			note:          "Auto-linked refund transaction",
		})
	}
	if err := refundRows.Close(); err != nil {
		return model.ClosingBatchSummary{}, err
	}
	for _, item := range items {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO location_daily_closing_item (
				closing_batch_id, reference_entity_type, reference_entity_id, amount, entry_direction, entry_note
			) VALUES (?, ?, ?, ?, ?, ?)`,
			closingID,
			item.referenceType,
			item.referenceID,
			item.amount,
			item.direction,
			item.note,
		); err != nil {
			return model.ClosingBatchSummary{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return model.ClosingBatchSummary{}, err
	}
	return r.fetchClosingBatchSummary(ctx, closingID)
}

func (r *MySQLRepository) UpdateClosingBatchStatus(ctx context.Context, adminUserID int64, closingID string, req model.UpdateClosingStatusRequest) (model.ClosingBatchSummary, error) {
	closingPK, err := parseIDString(closingID)
	if err != nil {
		return model.ClosingBatchSummary{}, err
	}
	nextStatus := normalizeClosingStatus(req.Status)
	if nextStatus == "" {
		return model.ClosingBatchSummary{}, errors.New("status closing tidak valid")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.ClosingBatchSummary{}, err
	}
	defer tx.Rollback()

	var oldStatus string
	if err := tx.QueryRowContext(ctx, `SELECT closing_status FROM location_daily_closing_batch WHERE id = ?`, closingPK).Scan(&oldStatus); err != nil {
		return model.ClosingBatchSummary{}, err
	}

	now := time.Now()
	if _, err := tx.ExecContext(ctx, `
		UPDATE location_daily_closing_batch
		SET closing_status = ?,
			reviewed_by_user_id = CASE WHEN ? = 'reviewed' THEN ? ELSE reviewed_by_user_id END,
			reviewed_at = CASE WHEN ? = 'reviewed' THEN ? ELSE reviewed_at END,
			approved_by_user_id = CASE WHEN ? = 'approved' THEN ? ELSE approved_by_user_id END,
			approved_at = CASE WHEN ? = 'approved' THEN ? ELSE approved_at END
		WHERE id = ?`,
		nextStatus,
		nextStatus, adminUserID,
		nextStatus, now,
		nextStatus, adminUserID,
		nextStatus, now,
		closingPK,
	); err != nil {
		return model.ClosingBatchSummary{}, err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO location_daily_closing_history (
			closing_batch_id, old_closing_status, new_closing_status, change_note, changed_by_user_id, changed_at
		) VALUES (?, ?, ?, ?, ?, ?)`,
		closingPK,
		oldStatus,
		nextStatus,
		strings.TrimSpace(req.ChangeNote),
		adminUserID,
		now,
	); err != nil {
		return model.ClosingBatchSummary{}, err
	}
	if err := tx.Commit(); err != nil {
		return model.ClosingBatchSummary{}, err
	}
	return r.fetchClosingBatchSummary(ctx, closingPK)
}

func (r *MySQLRepository) fetchDisputeCaseSummary(ctx context.Context, disputeID int64) (model.DisputeCaseSummary, error) {
	var item model.DisputeCaseSummary
	var assignedBy sql.NullInt64
	var resolvedAt sql.NullTime
	var resolutionNote sql.NullString
	var openedAt time.Time
	if err := r.db.QueryRowContext(ctx, `
		SELECT id, dispute_case_code, reference_entity_type, reference_entity_id, case_type, case_status,
		       opened_by_user_id, assigned_to_user_id, opened_at, resolved_at, resolution_note
		FROM financial_dispute_case
		WHERE id = ?`,
		disputeID,
	).Scan(
		&item.ID,
		&item.DisputeCaseCode,
		&item.ReferenceEntityType,
		&item.ReferenceEntityID,
		&item.CaseType,
		&item.Status,
		&item.OpenedByUserID,
		&assignedBy,
		&openedAt,
		&resolvedAt,
		&resolutionNote,
	); err != nil {
		return model.DisputeCaseSummary{}, err
	}
	item.AssignedToUserID = assignedBy.Int64
	item.OpenedAt = openedAt.Format(time.RFC3339)
	if resolvedAt.Valid {
		item.ResolvedAt = resolvedAt.Time.Format(time.RFC3339)
	}
	if resolutionNote.Valid {
		item.ResolutionNote = resolutionNote.String
	}
	return item, nil
}

func (r *MySQLRepository) fetchRefundTransactionSummary(ctx context.Context, refundID int64) (model.RefundTransactionSummary, error) {
	var item model.RefundTransactionSummary
	var paymentEventID, walletID, approvedBy sql.NullInt64
	var approvedAt, processedAt sql.NullTime
	var requestedAt time.Time
	if err := r.db.QueryRowContext(ctx, `
		SELECT id, refund_transaction_code, reference_entity_type, reference_entity_id, payment_event_id, wallet_id,
		       refund_amount, currency_code, refund_reason, refund_status, requested_by_user_id, approved_by_user_id,
		       requested_at, approved_at, processed_at
		FROM financial_refund_transaction
		WHERE id = ?`,
		refundID,
	).Scan(
		&item.ID,
		&item.RefundTransactionCode,
		&item.ReferenceEntityType,
		&item.ReferenceEntityID,
		&paymentEventID,
		&walletID,
		&item.RefundAmount,
		&item.CurrencyCode,
		&item.RefundReason,
		&item.Status,
		&item.RequestedByUserID,
		&approvedBy,
		&requestedAt,
		&approvedAt,
		&processedAt,
	); err != nil {
		return model.RefundTransactionSummary{}, err
	}
	item.PaymentEventID = paymentEventID.Int64
	item.WalletID = walletID.Int64
	item.ApprovedByUserID = approvedBy.Int64
	item.RequestedAt = requestedAt.Format(time.RFC3339)
	if approvedAt.Valid {
		item.ApprovedAt = approvedAt.Time.Format(time.RFC3339)
	}
	if processedAt.Valid {
		item.ProcessedAt = processedAt.Time.Format(time.RFC3339)
	}
	return item, nil
}

func (r *MySQLRepository) fetchClosingBatchSummary(ctx context.Context, closingID int64) (model.ClosingBatchSummary, error) {
	var item model.ClosingBatchSummary
	var reviewedBy, approvedBy sql.NullInt64
	var actual sql.NullInt64
	var submittedAt, reviewedAt, approvedAt sql.NullTime
	var closingDate time.Time
	if err := r.db.QueryRowContext(ctx, `
		SELECT id, closing_batch_code, location_id, closing_date, opening_balance_amount, cash_sales_amount,
		       cashless_sales_amount, topup_amount, refund_amount, adjustment_amount,
		       expected_closing_balance_amount, actual_closing_balance_amount, variance_amount, closing_status,
		       submitted_by_user_id, reviewed_by_user_id, approved_by_user_id, submitted_at, reviewed_at, approved_at
		FROM location_daily_closing_batch
		WHERE id = ?`,
		closingID,
	).Scan(
		&item.ID,
		&item.ClosingBatchCode,
		&item.LocationID,
		&closingDate,
		&item.OpeningBalanceAmount,
		&item.CashSalesAmount,
		&item.CashlessSalesAmount,
		&item.TopupAmount,
		&item.RefundAmount,
		&item.AdjustmentAmount,
		&item.ExpectedClosingBalanceAmount,
		&actual,
		&item.VarianceAmount,
		&item.Status,
		&item.SubmittedByUserID,
		&reviewedBy,
		&approvedBy,
		&submittedAt,
		&reviewedAt,
		&approvedAt,
	); err != nil {
		return model.ClosingBatchSummary{}, err
	}
	item.ClosingDate = closingDate.Format("2006-01-02")
	item.ActualClosingBalanceAmount = actual.Int64
	item.ReviewedByUserID = reviewedBy.Int64
	item.ApprovedByUserID = approvedBy.Int64
	if submittedAt.Valid {
		item.SubmittedAt = submittedAt.Time.Format(time.RFC3339)
	}
	if reviewedAt.Valid {
		item.ReviewedAt = reviewedAt.Time.Format(time.RFC3339)
	}
	if approvedAt.Valid {
		item.ApprovedAt = approvedAt.Time.Format(time.RFC3339)
	}
	return item, nil
}

func normalizeDisputeStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "open":
		return "open"
	case "investigating":
		return "investigating"
	case "resolved":
		return "resolved"
	case "closed":
		return "closed"
	default:
		return ""
	}
}

func normalizeRefundStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "requested":
		return "requested"
	case "approved":
		return "approved"
	case "processed":
		return "processed"
	case "settled":
		return "settled"
	case "cancelled":
		return "cancelled"
	default:
		return ""
	}
}

func normalizeClosingStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "draft":
		return "draft"
	case "submitted":
		return "submitted"
	case "reviewed":
		return "reviewed"
	case "approved":
		return "approved"
	case "rejected":
		return "rejected"
	default:
		return ""
	}
}

func nextBusinessCode(prefix string) string {
	buffer := make([]byte, 3)
	if _, err := rand.Read(buffer); err != nil {
		return fmt.Sprintf("%s-%s-%d", prefix, time.Now().Format("20060102150405"), time.Now().UnixNano())
	}
	return fmt.Sprintf("%s-%s-%s", prefix, time.Now().Format("20060102150405"), strings.ToUpper(hex.EncodeToString(buffer)))
}

func nullableZeroInt64(value int64) any {
	if value <= 0 {
		return nil
	}
	return value
}

func formatOptionalID(value int64) string {
	if value <= 0 {
		return ""
	}
	return fmt.Sprintf("%d", value)
}

func nullableStringValue(value sql.NullString) any {
	if value.Valid {
		return value.String
	}
	return nil
}

func nullableStringOrNil(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}

func nullableInt64Value(value sql.NullInt64) any {
	if value.Valid {
		return value.Int64
	}
	return nil
}
