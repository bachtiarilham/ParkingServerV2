package web

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"modulegue/internal/domain/dashboard"
)

type MySQLRepository struct {
	db              *sql.DB
	adminTablesOnce sync.Once
	adminTablesErr  error
}

func NewMySQL(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) FindAdminByIdentity(ctx context.Context, identity string) (*dashboard.AuthRecord, error) {
	identity = strings.TrimSpace(identity)
	if identity == "" {
		return nil, fmt.Errorf("identity wajib diisi")
	}

	row := r.db.QueryRowContext(ctx,
		//  `
		// SELECT
		// 	su.id,
		// 	COALESCE(su.full_name, ''),
		// 	COALESCE(su.phone_number, ''),
		// 	COALESCE(su.email, ''),
		// 	COALESCE(su.username, ''),
		// 	COALESCE(su.role_id, 0),
		// 	COALESCE(sr.role_code, ''),
		// 	COALESCE(su.password_hash, ''),
		// 	COALESCE(su.is_verified, 0)
		// FROM system_user su
		// LEFT JOIN system_role sr ON sr.id = su.role_id
		// WHERE su.username = ? OR su.email = ? OR su.phone_number = ?
		// LIMIT 1`,
		`SELECT
		ui.id AS userId,
		ui.role_id AS roleId,
		ua.password_hash AS passwordHash
	FROM user_identity ui
	JOIN user_auth ua
		ON ua.user_id = ui.id
	JOIN master_role mr
		ON mr.id = ui.role_id
	WHERE
		(
			ui.username = ?
			OR ui.email = ?
			OR ui.phone_number = ?
		)
		AND ui.status = 'ACTIVE'
		AND mr.is_active = 1
		AND (
			ua.locked_until IS NULL
			OR ua.locked_until <= NOW()
		)
	LIMIT 1`,
		identity, identity, identity,
	)

	var item dashboard.AuthRecord
	var verified bool
	if err := row.Scan(
		&item.ID,
		&item.FullName,
		&item.PhoneNumber,
		&item.Email,
		&item.Username,
		&item.RoleID,
		&item.RoleCode,
		&item.PasswordHash,
		&verified,
	); err != nil {
		return nil, err
	}
	item.IsVerified = verified
	return &item, nil
}

func (r *MySQLRepository) FindAdminByID(ctx context.Context, userID int64) (*dashboard.AuthRecord, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT
			su.id,
			COALESCE(su.full_name, ''),
			COALESCE(su.phone_number, ''),
			COALESCE(su.email, ''),
			COALESCE(su.username, ''),
			COALESCE(su.role_id, 0),
			COALESCE(sr.role_code, ''),
			COALESCE(su.password_hash, ''),
			COALESCE(su.is_verified, 0)
		FROM system_user su
		LEFT JOIN system_role sr ON sr.id = su.role_id
		WHERE su.id = ?
		LIMIT 1`,
		userID,
	)

	var item dashboard.AuthRecord
	var verified bool
	if err := row.Scan(
		&item.ID,
		&item.FullName,
		&item.PhoneNumber,
		&item.Email,
		&item.Username,
		&item.RoleID,
		&item.RoleCode,
		&item.PasswordHash,
		&verified,
	); err != nil {
		return nil, err
	}
	item.IsVerified = verified
	return &item, nil
}

func (r *MySQLRepository) GetDashboardOverview(ctx context.Context) (dashboard.DashboardOverview, error) {
	stats := []dashboard.StatCard{
		{Label: "Transaksi", Value: "0", Note: "Total transaksi", Icon: "receipt", Tone: "blue"},
		{Label: "Pendapatan", Value: "Rp 0", Note: "Pendapatan masuk", Icon: "cash", Tone: "green"},
		{Label: "Lokasi", Value: "0", Note: "Zona aktif", Icon: "map", Tone: "orange"},
		{Label: "Petugas", Value: "0", Note: "Petugas terpantau", Icon: "users", Tone: "purple"},
	}

	totalTransactions, _ := r.countQuery(ctx, `SELECT COUNT(*) FROM financial_parking_transaction`)
	paidTransactions, _ := r.countQuery(ctx, `SELECT COUNT(*) FROM financial_parking_transaction WHERE transaction_status IN ('paid', 'refunded_partial', 'refunded_full')`)
	totalLocations, _ := r.countQuery(ctx, `SELECT COUNT(*) FROM parking_location`)
	totalOfficers, _ := r.countQuery(ctx, `
		SELECT COUNT(*)
		FROM system_user su
		LEFT JOIN system_role sr ON sr.id = su.role_id
		WHERE LOWER(COALESCE(sr.role_code, '')) IN ('officer', 'sco')`)
	revenue, _ := r.sumQuery(ctx, `SELECT COALESCE(SUM(final_amount), 0) FROM financial_parking_transaction WHERE transaction_status IN ('paid', 'refunded_partial', 'refunded_full')`)

	stats[0].Value = fmt.Sprintf("%d", totalTransactions)
	stats[1].Value = formatIDR(revenue)
	stats[2].Value = fmt.Sprintf("%d", totalLocations)
	stats[3].Value = fmt.Sprintf("%d", totalOfficers)
	stats[0].Note = fmt.Sprintf("%d sudah lunas", paidTransactions)

	transactions, _ := r.listRecentTransactions(ctx, 10)
	locations, _ := r.listParkingLocations(ctx, 10)
	alerts, _ := r.listAlerts(ctx, 5)

	return dashboard.DashboardOverview{
		DashboardStats:        stats,
		DashboardTransactions: transactions,
		DashboardAlerts:       alerts,
		ParkingLocations:      locations,
		HourlyTraffic:         []dashboard.HourlyTrafficPoint{},
		RevenueByLocation:     []dashboard.LocationMetric{},
		OccupancyByLocation:   []dashboard.LocationMetric{},
		ComparisonMetrics:     []dashboard.ComparisonMetric{},
		ParkingHeatmap:        []dashboard.HeatmapPoint{},
		PriorityActions:       []dashboard.ActionItem{},
		FieldOfficers:         []dashboard.RowItem{},
	}, nil
}

func (r *MySQLRepository) GetMonitoringOverview(ctx context.Context) (dashboard.MonitoringOverview, error) {
	zones, _ := r.listZones(ctx)
	locations, _ := r.listParkingLocations(ctx, 50)
	officers, _ := r.listOfficerOptions(ctx)

	topFilters := dashboard.TopFilters{
		Zones:    zones,
		Dates:    time.Now().Format("2006-01-02"),
		Officers: officerNames(officers),
	}

	monitoringZones := make([]dashboard.RowItem, 0, len(locations))
	for _, loc := range locations {
		monitoringZones = append(monitoringZones, dashboard.RowItem{
			LocationID: loc.ID,
			Primary:    loc.Name,
			Secondary:  loc.Zone,
			Location:   loc.Address,
			Status:     loc.OfficerStatus,
			StatusTone: occupancyTone(loc.OccupancyLabel),
			ValueA:     fmt.Sprintf("%d motor", loc.Motorcycles),
			ValueB:     fmt.Sprintf("%d mobil", loc.Cars),
			Note:       loc.OfficerName,
		})
	}

	return dashboard.MonitoringOverview{
		TopFilters:            topFilters,
		MonitoringZones:       monitoringZones,
		ParkingLocations:      locations,
		ParkingOfficerOptions: officers,
	}, nil
}

func (r *MySQLRepository) GetOfficerOverview(ctx context.Context) (dashboard.OfficerOverview, error) {
	locations, _ := r.listParkingLocations(ctx, 50)
	officers, _ := r.listOfficerOptions(ctx)

	activeCount, _ := r.countQuery(ctx, `
		SELECT COUNT(*)
		FROM system_user su
		LEFT JOIN system_role sr ON sr.id = su.role_id
		WHERE LOWER(COALESCE(sr.role_code, '')) IN ('officer', 'sco')
		  AND COALESCE(su.employment_status, '') <> 'inactive'`)
	inactiveCount, _ := r.countQuery(ctx, `
		SELECT COUNT(*)
		FROM system_user su
		LEFT JOIN system_role sr ON sr.id = su.role_id
		WHERE LOWER(COALESCE(sr.role_code, '')) IN ('officer', 'sco')
		  AND LOWER(COALESCE(su.employment_status, '')) = 'inactive'`)
	totalCount, _ := r.countQuery(ctx, `
		SELECT COUNT(*)
		FROM system_user su
		LEFT JOIN system_role sr ON sr.id = su.role_id
		WHERE LOWER(COALESCE(sr.role_code, '')) IN ('officer', 'sco')`)

	stats := []dashboard.StatCard{
		{Label: "Petugas", Value: fmt.Sprintf("%d", totalCount), Note: "Seluruh petugas", Icon: "users", Tone: "blue"},
		{Label: "Aktif", Value: fmt.Sprintf("%d", activeCount), Note: "Sedang bertugas", Icon: "check", Tone: "green"},
		{Label: "Nonaktif", Value: fmt.Sprintf("%d", inactiveCount), Note: "Perlu perhatian", Icon: "pause", Tone: "orange"},
	}

	return dashboard.OfficerOverview{
		OfficerStats:          stats,
		ParkingOfficerOptions: officers,
		ParkingLocations:      locations,
	}, nil
}

func (r *MySQLRepository) GetTransactionsOverview(ctx context.Context) (dashboard.TransactionsOverview, error) {
	totalTransactions, _ := r.countQuery(ctx, `SELECT COUNT(*) FROM financial_parking_transaction`)
	paidTransactions, _ := r.countQuery(ctx, `SELECT COUNT(*) FROM financial_parking_transaction WHERE transaction_status = 'paid'`)
	unpaidTransactions, _ := r.countQuery(ctx, `SELECT COUNT(*) FROM financial_parking_transaction WHERE transaction_status = 'unpaid'`)
	refundedTransactions, _ := r.countQuery(ctx, `SELECT COUNT(*) FROM financial_parking_transaction WHERE transaction_status IN ('refunded_partial', 'refunded_full')`)
	revenue, _ := r.sumQuery(ctx, `SELECT COALESCE(SUM(final_amount), 0) FROM financial_parking_transaction WHERE transaction_status IN ('paid', 'refunded_partial', 'refunded_full')`)

	stats := []dashboard.StatCard{
		{Label: "Transaksi", Value: fmt.Sprintf("%d", totalTransactions), Note: "Semua transaksi", Icon: "receipt", Tone: "blue"},
		{Label: "Lunas", Value: fmt.Sprintf("%d", paidTransactions), Note: "Sudah dibayar", Icon: "check", Tone: "green"},
		{Label: "Belum Lunas", Value: fmt.Sprintf("%d", unpaidTransactions), Note: "Butuh tindak lanjut", Icon: "clock", Tone: "orange"},
		{Label: "Refund", Value: fmt.Sprintf("%d", refundedTransactions), Note: formatIDR(revenue), Icon: "refresh", Tone: "purple"},
	}

	rows, _ := r.listTransactionRows(ctx, 20)
	breakdown, _ := r.listPaymentBreakdown(ctx)
	issues, _ := r.listTransactionIssues(ctx)

	return dashboard.TransactionsOverview{
		TransactionStats:      stats,
		TransactionRows:       rows,
		PaymentBreakdownItems: breakdown,
		TransactionIssueItems: issues,
		ExportQueueItems:      []dashboard.ExportQueueItem{},
	}, nil
}

func (r *MySQLRepository) GetSettingsOverview(ctx context.Context) (dashboard.SettingsOverview, error) {
	alertRules, _ := r.listAlertRules(ctx)
	shiftTemplates, _ := r.listShiftTemplates(ctx)
	defaultTariffs, _ := r.listDefaultTariffs(ctx)
	adminRoles, _ := r.listAdminRoles(ctx)
	notifications, _ := r.listNotifications(ctx)
	paymentMethods, _ := r.listPaymentMethods(ctx)

	return dashboard.SettingsOverview{
		AlertRuleItems:        alertRules,
		DefaultShiftTemplates: shiftTemplates,
		DefaultTariffItems:    defaultTariffs,
		AdminRoleItems:        adminRoles,
		NotificationItems:     notifications,
		PaymentMethodItems:    paymentMethods,
	}, nil
}

func (r *MySQLRepository) countQuery(ctx context.Context, query string) (int64, error) {
	var value int64
	if err := r.db.QueryRowContext(ctx, query).Scan(&value); err != nil {
		return 0, err
	}
	return value, nil
}

func (r *MySQLRepository) sumQuery(ctx context.Context, query string) (int64, error) {
	var value sql.NullInt64
	if err := r.db.QueryRowContext(ctx, query).Scan(&value); err != nil {
		return 0, err
	}
	if value.Valid {
		return value.Int64, nil
	}
	return 0, nil
}

func (r *MySQLRepository) listRecentTransactions(ctx context.Context, limit int) ([]dashboard.RowItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			fpt.id,
			COALESCE(fpt.transaction_code, ''),
			COALESCE(fpt.transaction_status, ''),
			COALESCE(fpt.operation_type, ''),
			COALESCE(fpt.transaction_source, ''),
			COALESCE(pl.location_name, ''),
			COALESCE(fpt.final_amount, 0),
			COALESCE(DATE_FORMAT(fpt.occurred_at, '%Y-%m-%d %H:%i'), '')
		FROM financial_parking_transaction fpt
		LEFT JOIN parking_location pl ON pl.id = fpt.location_id
		ORDER BY fpt.created_at DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]dashboard.RowItem, 0)
	for rows.Next() {
		var id int64
		var code, status, opType, source, location string
		var amount int64
		var occurredAt string
		if err := rows.Scan(&id, &code, &status, &opType, &source, &location, &amount, &occurredAt); err != nil {
			return nil, err
		}
		items = append(items, dashboard.RowItem{
			TransactionID: fmt.Sprintf("%d", id),
			Primary:       code,
			Secondary:     fmt.Sprintf("%s • %s", opType, source),
			Status:        status,
			StatusTone:    toneFromStatus(status),
			Location:      location,
			Price:         formatIDR(amount),
			Time:          occurredAt,
		})
	}
	return items, nil
}

func (r *MySQLRepository) listParkingLocations(ctx context.Context, limit int) ([]dashboard.ParkingLocation, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			pl.id,
			COALESCE(pl.location_name, ''),
			COALESCE(pz.zone_name, ''),
			COALESCE(pl.street_address, ''),
			COALESCE(pl.center_latitude, 0),
			COALESCE(pl.center_longitude, 0),
			COALESCE(als.tariff_motor_amount, 0),
			COALESCE(als.tariff_car_amount, 0)
		FROM parking_location pl
		LEFT JOIN parking_zone pz ON pz.id = pl.zone_id
		LEFT JOIN admin_location_setting als ON als.location_id = pl.id
		ORDER BY pl.id DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]dashboard.ParkingLocation, 0)
	for rows.Next() {
		var id int64
		var name, zone, address string
		var lat, lng, tariffMotor, tariffCar float64
		if err := rows.Scan(&id, &name, &zone, &address, &lat, &lng, &tariffMotor, &tariffCar); err != nil {
			return nil, err
		}
		items = append(items, dashboard.ParkingLocation{
			ID:             fmt.Sprintf("%d", id),
			Name:           name,
			Zone:           zone,
			Address:        address,
			Lat:            lat,
			Lng:            lng,
			TariffMotor:    int64(tariffMotor),
			TariffMobil:    int64(tariffCar),
			OccupancyLabel: "Normal",
		})
	}
	return items, nil
}

func (r *MySQLRepository) listAlerts(ctx context.Context, limit int) ([]dashboard.AlertItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			COALESCE(alert_title, ''),
			COALESCE(alert_detail, ''),
			COALESCE(alert_priority, ''),
			COALESCE(alert_action, '')
		FROM admin_alert_event
		WHERE COALESCE(alert_status, '') = 'open'
		ORDER BY id DESC
		LIMIT ?`, limit)
	if err != nil {
		return []dashboard.AlertItem{}, nil
	}
	defer rows.Close()

	items := make([]dashboard.AlertItem, 0)
	for rows.Next() {
		var item dashboard.AlertItem
		if err := rows.Scan(&item.Title, &item.Detail, &item.Priority, &item.Action); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *MySQLRepository) listZones(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT DISTINCT COALESCE(zone_name, '') FROM parking_zone ORDER BY zone_name`)
	if err != nil {
		return []string{}, nil
	}
	defer rows.Close()

	items := make([]string, 0)
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		if strings.TrimSpace(value) != "" {
			items = append(items, value)
		}
	}
	return items, nil
}

func (r *MySQLRepository) listOfficerOptions(ctx context.Context) ([]dashboard.ParkingOfficerOption, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			su.id,
			COALESCE(su.full_name, ''),
			COALESCE(sr.role_code, ''),
			COALESCE(pz.zone_name, ''),
			COALESCE(pl.location_name, ''),
			COALESCE(pl.id, 0),
			COALESCE(pst.id, 0),
			COALESCE(DATE_FORMAT(pst.start_time, '%H:%i'), ''),
			COALESCE(DATE_FORMAT(pst.end_time, '%H:%i'), ''),
			COALESCE(oac.operational_status, 'off_duty')
		FROM system_user su
		LEFT JOIN system_role sr ON sr.id = su.role_id
		LEFT JOIN officer_assignment_current oac
			ON oac.officer_user_id = su.id
		   AND oac.effective_to IS NULL
		LEFT JOIN parking_location pl ON pl.id = oac.location_id
		LEFT JOIN parking_zone pz ON pz.id = pl.zone_id
		LEFT JOIN parking_shift_template pst ON pst.id = oac.shift_template_id
		WHERE LOWER(COALESCE(sr.role_code, '')) IN ('officer', 'sco')
		ORDER BY su.full_name`)
	if err != nil {
		return []dashboard.ParkingOfficerOption{}, nil
	}
	defer rows.Close()

	items := make([]dashboard.ParkingOfficerOption, 0)
	for rows.Next() {
		var (
			id                   int64
			name, roleCode       string
			homeZone, location   string
			locationID, shiftID  int64
			shiftStart, shiftEnd string
			status               string
		)
		if err := rows.Scan(&id, &name, &roleCode, &homeZone, &location, &locationID, &shiftID, &shiftStart, &shiftEnd, &status); err != nil {
			return nil, err
		}
		items = append(items, dashboard.ParkingOfficerOption{
			ID:                fmt.Sprintf("%d", id),
			Name:              name,
			Role:              roleCode,
			HomeZone:          homeZone,
			Availability:      status,
			AvailabilityNote:  "",
			CurrentAssignment: location,
			CurrentLocationID: fmt.Sprintf("%d", locationID),
			CurrentShiftID:    fmt.Sprintf("%d", shiftID),
			Status:            status,
			DefaultShiftStart: shiftStart,
			DefaultShiftEnd:   shiftEnd,
			DefaultStatus:     status,
		})
	}
	return items, nil
}

func (r *MySQLRepository) listTransactionRows(ctx context.Context, limit int) ([]dashboard.RowItem, error) {
	return r.listRecentTransactions(ctx, limit)
}

func (r *MySQLRepository) listPaymentBreakdown(ctx context.Context) ([]dashboard.PaymentBreakdownItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			COALESCE(NULLIF(payment_method, ''), 'unknown') AS method,
			COALESCE(SUM(final_amount), 0) AS amount
		FROM financial_parking_transaction
		GROUP BY COALESCE(NULLIF(payment_method, ''), 'unknown')
		ORDER BY amount DESC`)
	if err != nil {
		return []dashboard.PaymentBreakdownItem{}, nil
	}
	defer rows.Close()

	total, _ := r.sumQuery(ctx, `SELECT COALESCE(SUM(final_amount), 0) FROM financial_parking_transaction`)
	items := make([]dashboard.PaymentBreakdownItem, 0)
	for rows.Next() {
		var method string
		var amount int64
		if err := rows.Scan(&method, &amount); err != nil {
			return nil, err
		}
		share := "0%"
		if total > 0 {
			share = fmt.Sprintf("%.0f%%", (float64(amount)/float64(total))*100)
		}
		items = append(items, dashboard.PaymentBreakdownItem{
			Label:  strings.ToUpper(method),
			Amount: formatIDR(amount),
			Share:  share,
			Tone:   toneFromMethod(method),
		})
	}
	return items, nil
}

func (r *MySQLRepository) listTransactionIssues(ctx context.Context) ([]dashboard.TransactionIssueItem, error) {
	unpaid, _ := r.countQuery(ctx, `SELECT COUNT(*) FROM financial_parking_transaction WHERE transaction_status = 'unpaid'`)
	if unpaid == 0 {
		return []dashboard.TransactionIssueItem{}, nil
	}
	return []dashboard.TransactionIssueItem{
		{
			Title:  "Transaksi belum lunas",
			Detail: fmt.Sprintf("%d transaksi masih berstatus unpaid", unpaid),
			Action: "Tinjau transaksi pending",
			Tone:   "orange",
		},
	}, nil
}

func (r *MySQLRepository) listAlertRules(ctx context.Context) ([]dashboard.AlertRuleItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			COALESCE(rule_title, ''),
			COALESCE(threshold_text, ''),
			COALESCE(source_text, ''),
			COALESCE(pic_text, '')
		FROM admin_settings_alert_rule
		ORDER BY sort_order, id`)
	if err != nil {
		return []dashboard.AlertRuleItem{}, nil
	}
	defer rows.Close()

	items := make([]dashboard.AlertRuleItem, 0)
	for rows.Next() {
		var item dashboard.AlertRuleItem
		if err := rows.Scan(&item.Title, &item.Threshold, &item.Source, &item.PIC); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *MySQLRepository) listShiftTemplates(ctx context.Context) ([]dashboard.ShiftTemplateItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			COALESCE(shift_label, ''),
			COALESCE(CONCAT(start_time, '-', end_time), ''),
			COALESCE(use_case, '')
		FROM admin_settings_shift_template
		ORDER BY sort_order, id`)
	if err != nil {
		return []dashboard.ShiftTemplateItem{}, nil
	}
	defer rows.Close()

	items := make([]dashboard.ShiftTemplateItem, 0)
	for rows.Next() {
		var item dashboard.ShiftTemplateItem
		if err := rows.Scan(&item.Label, &item.Hours, &item.UseCase); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *MySQLRepository) listDefaultTariffs(ctx context.Context) ([]dashboard.DefaultTariffItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			COALESCE(vehicle_type_label, vehicle_type_code, ''),
			COALESCE(first_hour_amount, 0),
			COALESCE(next_hour_amount, 0),
			COALESCE(max_rate_amount, 0)
		FROM admin_settings_default_tariff
		ORDER BY vehicle_type_code`)
	if err != nil {
		return []dashboard.DefaultTariffItem{}, nil
	}
	defer rows.Close()

	items := make([]dashboard.DefaultTariffItem, 0)
	for rows.Next() {
		var item dashboard.DefaultTariffItem
		if err := rows.Scan(&item.VehicleType, &item.FirstHour, &item.NextHour, &item.MaxRate); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *MySQLRepository) listAdminRoles(ctx context.Context) ([]dashboard.AdminRoleItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			COALESCE(role_code, ''),
			COALESCE(role_name, role_code, ''),
			'system'
		FROM system_role
		ORDER BY role_code`)
	if err != nil {
		return []dashboard.AdminRoleItem{}, nil
	}
	defer rows.Close()

	items := make([]dashboard.AdminRoleItem, 0)
	for rows.Next() {
		var item dashboard.AdminRoleItem
		if err := rows.Scan(&item.Role, &item.Access, &item.Owner); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *MySQLRepository) listNotifications(ctx context.Context) ([]dashboard.NotificationItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			COALESCE(channel_name, ''),
			COALESCE(trigger_text, ''),
			COALESCE(response_text, '')
		FROM admin_settings_notification
		ORDER BY sort_order, id`)
	if err != nil {
		return []dashboard.NotificationItem{}, nil
	}
	defer rows.Close()

	items := make([]dashboard.NotificationItem, 0)
	for rows.Next() {
		var item dashboard.NotificationItem
		if err := rows.Scan(&item.Channel, &item.Trigger, &item.Response); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *MySQLRepository) listPaymentMethods(ctx context.Context) ([]dashboard.PaymentMethodItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			COALESCE(method_label, ''),
			COALESCE(is_enabled, 0),
			COALESCE(icon_code, '')
		FROM admin_settings_payment_method
		ORDER BY method_code`)
	if err != nil {
		return []dashboard.PaymentMethodItem{}, nil
	}
	defer rows.Close()

	items := make([]dashboard.PaymentMethodItem, 0)
	for rows.Next() {
		var item dashboard.PaymentMethodItem
		var enabled bool
		if err := rows.Scan(&item.Label, &enabled, &item.Icon); err != nil {
			return nil, err
		}
		item.Enabled = enabled
		items = append(items, item)
	}
	return items, nil
}

func officerNames(items []dashboard.ParkingOfficerOption) []string {
	names := make([]string, 0, len(items))
	for _, item := range items {
		if strings.TrimSpace(item.Name) != "" {
			names = append(names, item.Name)
		}
	}
	return names
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

func toneFromStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "paid":
		return "green"
	case "unpaid":
		return "orange"
	case "refunded_partial", "refunded_full":
		return "purple"
	default:
		return "blue"
	}
}

func toneFromMethod(method string) string {
	switch strings.ToLower(strings.TrimSpace(method)) {
	case "cash":
		return "green"
	case "qris":
		return "blue"
	default:
		return "orange"
	}
}

func occupancyTone(label string) string {
	switch strings.ToLower(strings.TrimSpace(label)) {
	case "zona padat":
		return "red"
	case "ramai":
		return "orange"
	case "normal":
		return "green"
	default:
		return "blue"
	}
}
