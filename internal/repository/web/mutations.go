package web

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"modulegue/internal/domain/dashboard"
)

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
				PRIMARY KEY (location_id)
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
				PRIMARY KEY (id)
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
				PRIMARY KEY (id)
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
				PRIMARY KEY (id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
			`CREATE TABLE IF NOT EXISTS admin_settings_payment_method (
				method_code VARCHAR(64) NOT NULL,
				method_label VARCHAR(191) NOT NULL,
				icon_code VARCHAR(32) NOT NULL,
				is_enabled TINYINT(1) NOT NULL DEFAULT 1,
				updated_by_user_id BIGINT UNSIGNED NULL,
				updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				PRIMARY KEY (method_code)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
			`CREATE TABLE IF NOT EXISTS admin_settings_default_tariff (
				vehicle_type_code VARCHAR(64) NOT NULL,
				vehicle_type_label VARCHAR(191) NOT NULL,
				first_hour_amount BIGINT NOT NULL DEFAULT 0,
				next_hour_amount BIGINT NOT NULL DEFAULT 0,
				max_rate_amount BIGINT NOT NULL DEFAULT 0,
				updated_by_user_id BIGINT UNSIGNED NULL,
				updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				PRIMARY KEY (vehicle_type_code)
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

func (r *MySQLRepository) UpdateSettingsOverview(ctx context.Context, adminUserID int64, req dashboard.UpdateSettingsOverviewRequest) (dashboard.SettingsOverview, error) {
	if err := r.ensureAdminTables(ctx); err != nil {
		return dashboard.SettingsOverview{}, err
	}
	if len(req.AlertRuleItems) == 0 {
		return dashboard.SettingsOverview{}, errors.New("alert rule wajib diisi")
	}
	if len(req.DefaultShiftTemplates) == 0 {
		return dashboard.SettingsOverview{}, errors.New("template shift wajib diisi")
	}
	if len(req.DefaultTariffItems) == 0 {
		return dashboard.SettingsOverview{}, errors.New("tarif default wajib diisi")
	}
	if len(req.NotificationItems) == 0 {
		return dashboard.SettingsOverview{}, errors.New("notifikasi wajib diisi")
	}
	if len(req.PaymentMethodItems) == 0 {
		return dashboard.SettingsOverview{}, errors.New("metode pembayaran wajib diisi")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return dashboard.SettingsOverview{}, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM admin_settings_alert_rule`); err != nil {
		return dashboard.SettingsOverview{}, err
	}
	for index, item := range req.AlertRuleItems {
		if strings.TrimSpace(item.Title) == "" || strings.TrimSpace(item.Threshold) == "" || strings.TrimSpace(item.Source) == "" || strings.TrimSpace(item.PIC) == "" {
			return dashboard.SettingsOverview{}, errors.New("setiap rule alert wajib memiliki title, threshold, source, dan PIC")
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
			return dashboard.SettingsOverview{}, err
		}
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM admin_settings_shift_template`); err != nil {
		return dashboard.SettingsOverview{}, err
	}
	for index, item := range req.DefaultShiftTemplates {
		start, end := splitHoursRange(item.Hours)
		if strings.TrimSpace(item.Label) == "" || start == "" || end == "" || strings.TrimSpace(item.UseCase) == "" {
			return dashboard.SettingsOverview{}, errors.New("setiap template shift wajib memiliki label, jam, dan use case")
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
			return dashboard.SettingsOverview{}, err
		}
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM admin_settings_default_tariff`); err != nil {
		return dashboard.SettingsOverview{}, err
	}
	for _, item := range req.DefaultTariffItems {
		vehicleType := strings.TrimSpace(item.VehicleType)
		if vehicleType == "" {
			return dashboard.SettingsOverview{}, errors.New("jenis kendaraan tarif default wajib diisi")
		}
		if item.FirstHour < 0 || item.NextHour < 0 || item.MaxRate < 0 {
			return dashboard.SettingsOverview{}, errors.New("nominal tarif default tidak boleh negatif")
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
			return dashboard.SettingsOverview{}, err
		}
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM admin_settings_notification`); err != nil {
		return dashboard.SettingsOverview{}, err
	}
	for index, item := range req.NotificationItems {
		if strings.TrimSpace(item.Channel) == "" || strings.TrimSpace(item.Trigger) == "" || strings.TrimSpace(item.Response) == "" {
			return dashboard.SettingsOverview{}, errors.New("setiap notifikasi wajib memiliki channel, trigger, dan response")
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
			return dashboard.SettingsOverview{}, err
		}
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM admin_settings_payment_method`); err != nil {
		return dashboard.SettingsOverview{}, err
	}
	for _, item := range req.PaymentMethodItems {
		if strings.TrimSpace(item.Label) == "" {
			return dashboard.SettingsOverview{}, errors.New("label metode pembayaran wajib diisi")
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
			return dashboard.SettingsOverview{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return dashboard.SettingsOverview{}, err
	}
	return r.GetSettingsOverview(ctx)
}

func (r *MySQLRepository) UpdateLocationSettings(ctx context.Context, adminUserID int64, locationID string, req dashboard.UpdateLocationSettingsRequest) (dashboard.ParkingLocation, error) {
	if err := r.ensureAdminTables(ctx); err != nil {
		return dashboard.ParkingLocation{}, err
	}
	locationPK, err := parseIDString(locationID)
	if err != nil {
		return dashboard.ParkingLocation{}, err
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
		return dashboard.ParkingLocation{}, err
	}
	location, err := r.findLocationByStringID(ctx, locationID)
	if err != nil {
		return dashboard.ParkingLocation{}, err
	}
	return location, nil
}

func (r *MySQLRepository) SaveLocationShiftTemplates(ctx context.Context, adminUserID int64, locationID string, shiftTemplates []dashboard.ParkingShiftTemplate) (dashboard.ParkingLocation, error) {
	locationPK, err := parseIDString(locationID)
	if err != nil {
		return dashboard.ParkingLocation{}, err
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return dashboard.ParkingLocation{}, err
	}
	defer tx.Rollback()

	existingRows, err := tx.QueryContext(ctx, `SELECT id FROM parking_shift_template WHERE location_id = ?`, locationPK)
	if err != nil {
		return dashboard.ParkingLocation{}, err
	}
	defer existingRows.Close()
	existing := map[int64]struct{}{}
	for existingRows.Next() {
		var id int64
		if err := existingRows.Scan(&id); err != nil {
			return dashboard.ParkingLocation{}, err
		}
		existing[id] = struct{}{}
	}
	if err := existingRows.Err(); err != nil {
		return dashboard.ParkingLocation{}, err
	}

	used := map[int64]struct{}{}
	for _, shift := range shiftTemplates {
		label := strings.TrimSpace(shift.Label)
		start := strings.TrimSpace(shift.Start)
		end := strings.TrimSpace(shift.End)
		if label == "" || start == "" || end == "" {
			return dashboard.ParkingLocation{}, errors.New("shift template wajib memiliki nama, waktu mulai, dan waktu selesai")
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
					return dashboard.ParkingLocation{}, err
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
			return dashboard.ParkingLocation{}, err
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
			return dashboard.ParkingLocation{}, err
		}
	}

	if _, err := tx.ExecContext(ctx, `UPDATE officer_assignment_current SET assigned_by_user_id = ?, updated_at = CURRENT_TIMESTAMP WHERE location_id = ? AND effective_to IS NULL`, adminUserID, locationPK); err != nil {
		return dashboard.ParkingLocation{}, err
	}

	if err := tx.Commit(); err != nil {
		return dashboard.ParkingLocation{}, err
	}
	location, err := r.findLocationByStringID(ctx, locationID)
	if err != nil {
		return dashboard.ParkingLocation{}, err
	}
	return location, nil
}

func (r *MySQLRepository) UpdateOfficerStatus(ctx context.Context, adminUserID int64, officerID string, status string) (dashboard.ParkingOfficerOption, error) {
	officerPK, err := parseIDString(officerID)
	if err != nil {
		return dashboard.ParkingOfficerOption{}, err
	}
	nextStatus := toOperationalStatus(status)
	var previous sql.NullString
	_ = r.db.QueryRowContext(ctx, `SELECT operational_status FROM officer_assignment_current WHERE officer_user_id = ? AND effective_to IS NULL LIMIT 1`, officerPK).Scan(&previous)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return dashboard.ParkingOfficerOption{}, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `
		UPDATE officer_assignment_current
		SET operational_status = ?, assigned_by_user_id = ?, updated_at = CURRENT_TIMESTAMP
		WHERE officer_user_id = ? AND effective_to IS NULL`,
		nextStatus, adminUserID, officerPK,
	); err != nil {
		return dashboard.ParkingOfficerOption{}, err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO officer_status_history (officer_user_id, old_operational_status, new_operational_status, change_reason, changed_by_user_id, changed_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`,
		officerPK, nullableStringValue(previous), nextStatus, "Status diperbarui dari dashboard admin", adminUserID,
	); err != nil {
		return dashboard.ParkingOfficerOption{}, err
	}
	if err := tx.Commit(); err != nil {
		return dashboard.ParkingOfficerOption{}, err
	}
	return r.findOfficerByStringID(ctx, officerID)
}

func (r *MySQLRepository) ApplyOfficerMutation(ctx context.Context, adminUserID int64, req dashboard.ApplyOfficerMutationRequest) (dashboard.ParkingOfficerOption, error) {
	officerPK, err := parseIDString(req.OfficerID)
	if err != nil {
		return dashboard.ParkingOfficerOption{}, err
	}
	locationPK, err := parseIDString(req.TargetLocationID)
	if err != nil {
		return dashboard.ParkingOfficerOption{}, err
	}
	shiftPK, err := parseIDString(req.TargetShiftID)
	if err != nil {
		return dashboard.ParkingOfficerOption{}, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return dashboard.ParkingOfficerOption{}, err
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
		return dashboard.ParkingOfficerOption{}, err
	}
	if strings.TrimSpace(previousStatus) == "" {
		previousStatus = "off_duty"
	}

	if err := tx.QueryRowContext(ctx, `SELECT zone_id, area_id FROM parking_location WHERE id = ? LIMIT 1`, locationPK).Scan(&zoneID, &areaID); err != nil {
		return dashboard.ParkingOfficerOption{}, err
	}

	if currentAssignmentID > 0 {
		if _, err := tx.ExecContext(ctx, `
			UPDATE officer_assignment_current
			SET location_id = ?, zone_id = ?, area_id = ?, shift_template_id = ?, assigned_by_user_id = ?, updated_at = CURRENT_TIMESTAMP
			WHERE id = ?`,
			locationPK, nullableInt64Value(zoneID), nullableInt64Value(areaID), shiftPK, adminUserID, currentAssignmentID,
		); err != nil {
			return dashboard.ParkingOfficerOption{}, err
		}
	} else {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO officer_assignment_current (
				officer_user_id, location_id, zone_id, area_id, shift_template_id, operational_status, effective_from, effective_to, assigned_by_user_id
			) VALUES (?, ?, ?, ?, ?, 'off_duty', CURRENT_TIMESTAMP, NULL, ?)`,
			officerPK, locationPK, nullableInt64Value(zoneID), nullableInt64Value(areaID), shiftPK, adminUserID,
		); err != nil {
			return dashboard.ParkingOfficerOption{}, err
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
		return dashboard.ParkingOfficerOption{}, err
	}

	if err := tx.Commit(); err != nil {
		return dashboard.ParkingOfficerOption{}, err
	}
	return r.findOfficerByStringID(ctx, req.OfficerID)
}

func (r *MySQLRepository) CreateDisputeCase(ctx context.Context, adminUserID int64, req dashboard.CreateDisputeCaseRequest) (dashboard.DisputeCaseSummary, error) {
	if strings.TrimSpace(req.ReferenceEntityType) == "" || req.ReferenceEntityID <= 0 {
		return dashboard.DisputeCaseSummary{}, errors.New("reference entity dispute wajib diisi")
	}
	if strings.TrimSpace(req.CaseType) == "" {
		return dashboard.DisputeCaseSummary{}, errors.New("case type wajib diisi")
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
		return dashboard.DisputeCaseSummary{}, err
	}
	disputeID, err := result.LastInsertId()
	if err != nil {
		return dashboard.DisputeCaseSummary{}, err
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
		return dashboard.DisputeCaseSummary{}, err
	}
	return r.fetchDisputeCaseSummary(ctx, disputeID)
}

func (r *MySQLRepository) UpdateDisputeCaseStatus(ctx context.Context, adminUserID int64, disputeID string, req dashboard.UpdateDisputeCaseStatusRequest) (dashboard.DisputeCaseSummary, error) {
	disputePK, err := parseIDString(disputeID)
	if err != nil {
		return dashboard.DisputeCaseSummary{}, err
	}
	nextStatus := normalizeDisputeStatus(req.Status)
	if nextStatus == "" {
		return dashboard.DisputeCaseSummary{}, errors.New("status dispute tidak valid")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return dashboard.DisputeCaseSummary{}, err
	}
	defer tx.Rollback()

	var oldStatus string
	if err := tx.QueryRowContext(ctx, `SELECT case_status FROM financial_dispute_case WHERE id = ?`, disputePK).Scan(&oldStatus); err != nil {
		return dashboard.DisputeCaseSummary{}, err
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
		return dashboard.DisputeCaseSummary{}, err
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
		return dashboard.DisputeCaseSummary{}, err
	}
	if err := tx.Commit(); err != nil {
		return dashboard.DisputeCaseSummary{}, err
	}
	return r.fetchDisputeCaseSummary(ctx, disputePK)
}

func (r *MySQLRepository) CreateRefundTransaction(ctx context.Context, adminUserID int64, req dashboard.CreateRefundTransactionRequest) (dashboard.RefundTransactionSummary, error) {
	if strings.TrimSpace(req.ReferenceEntityType) == "" || req.ReferenceEntityID <= 0 {
		return dashboard.RefundTransactionSummary{}, errors.New("reference entity refund wajib diisi")
	}
	if req.RefundAmount <= 0 {
		return dashboard.RefundTransactionSummary{}, errors.New("refund amount harus lebih besar dari nol")
	}
	if strings.TrimSpace(req.RefundReason) == "" {
		return dashboard.RefundTransactionSummary{}, errors.New("refund reason wajib diisi")
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
		return dashboard.RefundTransactionSummary{}, err
	}
	refundID, err := result.LastInsertId()
	if err != nil {
		return dashboard.RefundTransactionSummary{}, err
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
		return dashboard.RefundTransactionSummary{}, err
	}
	return r.fetchRefundTransactionSummary(ctx, refundID)
}

func (r *MySQLRepository) UpdateRefundTransactionStatus(ctx context.Context, adminUserID int64, refundID string, req dashboard.UpdateRefundStatusRequest) (dashboard.RefundTransactionSummary, error) {
	refundPK, err := parseIDString(refundID)
	if err != nil {
		return dashboard.RefundTransactionSummary{}, err
	}
	nextStatus := normalizeRefundStatus(req.Status)
	if nextStatus == "" {
		return dashboard.RefundTransactionSummary{}, errors.New("status refund tidak valid")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return dashboard.RefundTransactionSummary{}, err
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
		return dashboard.RefundTransactionSummary{}, err
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
		return dashboard.RefundTransactionSummary{}, err
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
		return dashboard.RefundTransactionSummary{}, err
	}

	if (nextStatus == "processed" || nextStatus == "settled") && refType == "financial_parking_transaction" {
		var txStatus string
		var finalAmount int64
		if err := tx.QueryRowContext(ctx, `SELECT transaction_status, final_amount FROM financial_parking_transaction WHERE id = ?`, refID).Scan(&txStatus, &finalAmount); err != nil {
			return dashboard.RefundTransactionSummary{}, err
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
			return dashboard.RefundTransactionSummary{}, err
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
			return dashboard.RefundTransactionSummary{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return dashboard.RefundTransactionSummary{}, err
	}
	return r.fetchRefundTransactionSummary(ctx, refundPK)
}

func (r *MySQLRepository) CreateClosingBatch(ctx context.Context, adminUserID int64, req dashboard.CreateClosingBatchRequest) (dashboard.ClosingBatchSummary, error) {
	if req.LocationID <= 0 {
		return dashboard.ClosingBatchSummary{}, errors.New("location id closing wajib diisi")
	}
	closingDate, err := time.Parse("2006-01-02", strings.TrimSpace(req.ClosingDate))
	if err != nil {
		return dashboard.ClosingBatchSummary{}, errors.New("closing date harus berformat YYYY-MM-DD")
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
		return dashboard.ClosingBatchSummary{}, err
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
		return dashboard.ClosingBatchSummary{}, err
	}

	const openingBalance, topupAmount, adjustmentAmount int64 = 0, 0, 0
	expected := openingBalance + cashSales + cashlessSales + topupAmount + adjustmentAmount - refundAmount
	actual := req.ActualClosingBalanceAmount
	variance := actual - expected

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return dashboard.ClosingBatchSummary{}, err
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
		return dashboard.ClosingBatchSummary{}, err
	}
	closingID, err := result.LastInsertId()
	if err != nil {
		return dashboard.ClosingBatchSummary{}, err
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
		return dashboard.ClosingBatchSummary{}, err
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
		return dashboard.ClosingBatchSummary{}, err
	}
	items := make([]closingItemSeed, 0)
	for transactionRows.Next() {
		var id, amount int64
		if err := transactionRows.Scan(&id, &amount); err != nil {
			transactionRows.Close()
			return dashboard.ClosingBatchSummary{}, err
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
		return dashboard.ClosingBatchSummary{}, err
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
		return dashboard.ClosingBatchSummary{}, err
	}
	for refundRows.Next() {
		var id, amount int64
		if err := refundRows.Scan(&id, &amount); err != nil {
			refundRows.Close()
			return dashboard.ClosingBatchSummary{}, err
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
		return dashboard.ClosingBatchSummary{}, err
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
			return dashboard.ClosingBatchSummary{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return dashboard.ClosingBatchSummary{}, err
	}
	return r.fetchClosingBatchSummary(ctx, closingID)
}

func (r *MySQLRepository) UpdateClosingBatchStatus(ctx context.Context, adminUserID int64, closingID string, req dashboard.UpdateClosingStatusRequest) (dashboard.ClosingBatchSummary, error) {
	closingPK, err := parseIDString(closingID)
	if err != nil {
		return dashboard.ClosingBatchSummary{}, err
	}
	nextStatus := normalizeClosingStatus(req.Status)
	if nextStatus == "" {
		return dashboard.ClosingBatchSummary{}, errors.New("status closing tidak valid")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return dashboard.ClosingBatchSummary{}, err
	}
	defer tx.Rollback()

	var oldStatus string
	if err := tx.QueryRowContext(ctx, `SELECT closing_status FROM location_daily_closing_batch WHERE id = ?`, closingPK).Scan(&oldStatus); err != nil {
		return dashboard.ClosingBatchSummary{}, err
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
		return dashboard.ClosingBatchSummary{}, err
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
		return dashboard.ClosingBatchSummary{}, err
	}
	if err := tx.Commit(); err != nil {
		return dashboard.ClosingBatchSummary{}, err
	}
	return r.fetchClosingBatchSummary(ctx, closingPK)
}

func (r *MySQLRepository) findLocationByStringID(ctx context.Context, locationID string) (dashboard.ParkingLocation, error) {
	locations, err := r.listParkingLocations(ctx, 1000)
	if err != nil {
		return dashboard.ParkingLocation{}, err
	}
	for _, location := range locations {
		if location.ID == strings.TrimSpace(locationID) {
			return location, nil
		}
	}
	return dashboard.ParkingLocation{}, sql.ErrNoRows
}

func (r *MySQLRepository) findOfficerByStringID(ctx context.Context, officerID string) (dashboard.ParkingOfficerOption, error) {
	officers, err := r.listOfficerOptions(ctx)
	if err != nil {
		return dashboard.ParkingOfficerOption{}, err
	}
	for _, officer := range officers {
		if officer.ID == strings.TrimSpace(officerID) {
			return officer, nil
		}
	}
	return dashboard.ParkingOfficerOption{}, sql.ErrNoRows
}

func (r *MySQLRepository) fetchDisputeCaseSummary(ctx context.Context, disputeID int64) (dashboard.DisputeCaseSummary, error) {
	var item dashboard.DisputeCaseSummary
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
		return dashboard.DisputeCaseSummary{}, err
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

func (r *MySQLRepository) fetchRefundTransactionSummary(ctx context.Context, refundID int64) (dashboard.RefundTransactionSummary, error) {
	var item dashboard.RefundTransactionSummary
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
		return dashboard.RefundTransactionSummary{}, err
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

func (r *MySQLRepository) fetchClosingBatchSummary(ctx context.Context, closingID int64) (dashboard.ClosingBatchSummary, error) {
	var item dashboard.ClosingBatchSummary
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
		return dashboard.ClosingBatchSummary{}, err
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

func splitHoursRange(raw string) (string, string) {
	parts := strings.Split(raw, "-")
	if len(parts) != 2 {
		return "", ""
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

func parseIDString(value string) (int64, error) {
	parsed, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil || parsed <= 0 {
		return 0, errors.New("id tidak valid")
	}
	return parsed, nil
}

func shiftCodeFromTemplate(shift dashboard.ParkingShiftTemplate) string {
	label := strings.ToUpper(strings.TrimSpace(shift.Label))
	label = strings.ReplaceAll(label, " ", "_")
	label = strings.ReplaceAll(label, "-", "_")
	if label == "" {
		label = "SHIFT"
	}
	return fmt.Sprintf("%s_%s_%s", label, strings.ReplaceAll(shift.Start, ":", ""), strings.ReplaceAll(shift.End, ":", ""))
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

func nextBusinessCode(prefix string) string {
	buffer := make([]byte, 3)
	if _, err := rand.Read(buffer); err != nil {
		return fmt.Sprintf("%s-%s-%d", prefix, time.Now().Format("20060102150405"), time.Now().UnixNano())
	}
	return fmt.Sprintf("%s-%s-%s", prefix, time.Now().Format("20060102150405"), strings.ToUpper(fmt.Sprintf("%x", buffer)))
}

func nullableZeroInt64(value int64) any {
	if value <= 0 {
		return nil
	}
	return value
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

func maxInt64(values ...int64) int64 {
	var max int64
	for i, value := range values {
		if i == 0 || value > max {
			max = value
		}
	}
	return max
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
