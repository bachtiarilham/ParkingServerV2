package payment

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"modulegue/core/utils"
	model "modulegue/internal/domain/mobile/model/payment"
	"modulegue/internal/domain/mobile/repository"
)

type PaymentRepositoryImpl struct {
	db *sql.DB
}

func NewPaymentRepositoryImpl(db *sql.DB) repository.PaymentRepository {
	return &PaymentRepositoryImpl{db: db}
}

func (r *PaymentRepositoryImpl) PostPaymentParking(ctx context.Context, req model.PostPaymentParkingRequestModel) (*model.PaymentBusinessModel, error) {
	query := `
	SELECT
		ps.id AS sessionId,
		ps.session_code AS sessionCode,
		ps.transaction_code AS transactionCode,

		ps.plate_number AS plateNumber,

		mvt.id AS vehicleTypeId,
		mvt.vehicle_type_code AS vehicleTypeCode,
		mvt.vehicle_type_name AS vehicleTypeName,

		lp.id AS locationId,
		lp.location_name AS locationName,
		lp.address AS locationAddress,

		la.id AS areaId,
		la.area_name AS areaName,

		ps.amount AS amount,
		ps.qr_expired_at AS qrExpiredAt,

		pt.payment_code AS paymentCode,

		mps.status_code AS parkingStatusCode,
		mpy.status_code AS paymentStatusCode

	FROM parking_session ps

	JOIN master_vehicle_type mvt
		ON mvt.id = ps.vehicle_type_id

	JOIN location_parking lp
		ON lp.id = ps.location_id

	LEFT JOIN location_area la
		ON la.id = ps.area_id

	JOIN payment_transaction pt
		ON pt.parking_session_id = ps.id

	JOIN master_parking_status mps
		ON mps.id = ps.parking_status_id

	JOIN master_payment_status mpy
		ON mpy.id = ps.payment_status_id

	WHERE ps.session_code = ?
	AND mps.status_code = 'WAITING_PAYMENT'
	AND mpy.status_code = 'PENDING'
	AND ps.qr_expired_at > NOW()

	LIMIT 1;
	`

	var (
		businessModel   model.PaymentBusinessModel
		qrExpiredAt     sql.NullTime
		sessionCode     sql.NullString
		transactionCode sql.NullString
		plateNumber     sql.NullString
		vehicleTypeID   sql.NullInt64
		vehicleCode     sql.NullString
		vehicleName     sql.NullString
		locationID      sql.NullInt64
		locationName    sql.NullString
		locationAddress sql.NullString
		areaID          sql.NullInt64
		areaName        sql.NullString
		amount          sql.NullInt64
		paymentCode     sql.NullString
		parkingStatus   sql.NullString
		paymentStatus   sql.NullString
	)

	if err := r.db.QueryRowContext(ctx, query, strings.TrimSpace(req.SessionCode)).Scan(
		&businessModel.SessionId,
		&sessionCode,
		&transactionCode,
		&plateNumber,
		&vehicleTypeID,
		&vehicleCode,
		&vehicleName,
		&locationID,
		&locationName,
		&locationAddress,
		&areaID,
		&areaName,
		&amount,
		&qrExpiredAt,
		&paymentCode,
		&parkingStatus,
		&paymentStatus,
	); err != nil {
		if err == sql.ErrNoRows {
			alreadyPaid, stateErr := r.isPaymentAlreadyCompleted(ctx, req.SessionCode)
			if stateErr != nil {
				return nil, stateErr
			}
			if alreadyPaid {
				return nil, model.ErrPaymentAlreadyCompleted
			}
			return nil, fmt.Errorf("kebutuhan payment tidak ditemukan")
		}
		return nil, fmt.Errorf("post payment parking: %w", err)
	}

	businessModel.SessionCode = utils.NullStringValue(sessionCode)
	businessModel.TransactionCode = utils.NullStringValue(transactionCode)
	businessModel.PlateNumber = utils.NullStringValue(plateNumber)
	businessModel.VehicleTypeId = utils.NullInt64Value(vehicleTypeID)
	businessModel.VehicleTypeCode = utils.NullStringValue(vehicleCode)
	businessModel.VehicleTypeName = utils.NullStringValue(vehicleName)
	businessModel.LocationId = utils.NullInt64Value(locationID)
	businessModel.LocationName = utils.NullStringValue(locationName)
	businessModel.Address = utils.NullStringValue(locationAddress)
	businessModel.AreaId = utils.NullInt64Value(areaID)
	businessModel.AreaName = utils.NullStringValue(areaName)
	businessModel.Amount = utils.NullInt64Value(amount)
	businessModel.QrExpiredAt = utils.NullTimeValue(qrExpiredAt)
	businessModel.PaymentCode = utils.NullStringValue(paymentCode)
	businessModel.ParkingStatusId = utils.NullStringValue(parkingStatus)
	businessModel.PaymentStatusId = utils.NullStringValue(paymentStatus)
	businessModel.CustomerUserId = req.CustomerUserId

	return &businessModel, nil
}

func (r *PaymentRepositoryImpl) isPaymentAlreadyCompleted(ctx context.Context, sessionCode string) (bool, error) {
	query := `
	SELECT
		COALESCE(mps.status_code, '') AS parkingStatusCode,
		COALESCE(mpy.status_code, '') AS paymentStatusCode,
		ps.paid_at
	FROM parking_session ps

	JOIN master_parking_status mps
		ON mps.id = ps.parking_status_id

	JOIN master_payment_status mpy
		ON mpy.id = ps.payment_status_id

	WHERE ps.session_code = ?
	LIMIT 1;
	`

	var (
		parkingStatus sql.NullString
		paymentStatus sql.NullString
		paidAt        sql.NullTime
	)

	if err := r.db.QueryRowContext(ctx, query, strings.TrimSpace(sessionCode)).Scan(&parkingStatus, &paymentStatus, &paidAt); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("check payment already completed: %w", err)
	}

	if utils.NullStringValue(parkingStatus) == "PAID" || utils.NullStringValue(paymentStatus) == "SUCCESS" || paidAt.Valid {
		return true, nil
	}

	return false, nil
}

func (r *PaymentRepositoryImpl) BindCustomerToSessionAndTransaction(ctx context.Context, req model.PaymentBusinessModel) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin bind customer tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	updateSessionQuery := `
	UPDATE parking_session
	SET
		customer_user_id = COALESCE(customer_user_id, ?),
		updated_at = NOW()
	WHERE id = ?
	AND (
			customer_user_id IS NULL
			OR customer_user_id = ?
	);
	`

	if _, err = tx.ExecContext(ctx, updateSessionQuery, req.CustomerUserId, req.SessionId, req.CustomerUserId); err != nil {
		return fmt.Errorf("bind customer parking session: %w", err)
	}

	updatePaymentQuery := `
	UPDATE payment_transaction
	SET
		user_id = COALESCE(user_id, ?),
		updated_at = NOW()
	WHERE parking_session_id = ?
	AND (
			user_id IS NULL
			OR user_id = ?
	);
	`

	if _, err = tx.ExecContext(ctx, updatePaymentQuery, req.CustomerUserId, req.SessionId, req.CustomerUserId); err != nil {
		return fmt.Errorf("bind customer payment transaction: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit bind customer tx: %w", err)
	}

	return nil
}

func (r *PaymentRepositoryImpl) UpdatePaymentTransactionSuccess(ctx context.Context, req model.PaymentBusinessModel) error {
	query := `
	UPDATE payment_transaction pt

	JOIN master_payment_status mps
		ON mps.status_code = 'SUCCESS'

	SET
		pt.payment_status_id = mps.id,
		pt.paid_at = NOW(),
		pt.failed_reason = NULL,
		pt.updated_at = NOW()

	WHERE pt.payment_code = ?
	AND pt.paid_at IS NULL;
	`

	if _, err := r.db.ExecContext(ctx, query, req.PaymentCode); err != nil {
		return fmt.Errorf("update payment transaction success: %w", err)
	}

	return nil
}

func (r *PaymentRepositoryImpl) UpdateParkingSessionSuccess(ctx context.Context, req model.PaymentBusinessModel) error {
	query := `
	UPDATE parking_session ps

	JOIN payment_transaction pt
		ON pt.parking_session_id = ps.id

	JOIN master_parking_status mps
		ON mps.status_code = 'PAID'

	JOIN master_payment_status mpys
		ON mpys.status_code = 'SUCCESS'

	SET
		ps.parking_status_id = mps.id,
		ps.payment_status_id = mpys.id,
		ps.paid_at = COALESCE(pt.paid_at, NOW()),
		ps.updated_at = NOW()

	WHERE pt.payment_code = ?;
	`

	if _, err := r.db.ExecContext(ctx, query, req.PaymentCode); err != nil {
		return fmt.Errorf("update parking session success: %w", err)
	}

	return nil
}

func (r *PaymentRepositoryImpl) BuatParkingReceipt(ctx context.Context, req model.PaymentBusinessModel) error {
	query := `
	INSERT INTO parking_receipt (
		parking_session_id,
		receipt_number,
		transaction_code,
		plate_number,
		vehicle_type_id,
		location_id,
		amount,
		payment_method_id,
		paid_at,
		created_at
	)
	SELECT
		ps.id,
		CONCAT('RCPT-', ps.transaction_code),
		ps.transaction_code,
		ps.plate_number,
		ps.vehicle_type_id,
		ps.location_id,
		ps.amount,
		pt.payment_method_id,
		ps.paid_at,
		NOW()

	FROM parking_session ps

	JOIN payment_transaction pt
		ON pt.parking_session_id = ps.id

	WHERE pt.payment_code = ?

	ON DUPLICATE KEY UPDATE
		paid_at = VALUES(paid_at);
	`

	if _, err := r.db.ExecContext(ctx, query, req.PaymentCode); err != nil {
		return fmt.Errorf("buat parking receipt: %w", err)
	}

	return nil
}

func (r *PaymentRepositoryImpl) BuatFinancialParkingTransaction(ctx context.Context, req model.PaymentBusinessModel) error {
	query := `
	INSERT INTO financial_parking_transaction (
		transaction_code,
		session_id,

		customer_user_id,
		jukir_user_id,
		officer_user_id,

		location_id,
		area_id,
		zone_id,

		vehicle_type_id,
		plate_number,

		operation_type,
		payment_method_id,
		transaction_status,

		base_amount,
		discount_amount,
		final_amount,

		company_share,
		jukir_share,

		tax_amount,
		fee_amount,

		occurred_at,
		paid_at,

		created_at,
		updated_at
	)
	SELECT
		ps.transaction_code,
		ps.id,

		ps.customer_user_id,
		ps.officer_user_id,
		ps.officer_user_id,

		ps.location_id,
		ps.area_id,
		lp.zone_id,

		ps.vehicle_type_id,
		ps.plate_number,

		'ON_STREET',
		pt.payment_method_id,
		'SUCCESS',

		ps.amount,
		0,
		ps.amount,

		ps.amount - FLOOR(
			ps.amount * COALESCE((
				SELECT CAST(setting_value AS UNSIGNED)
				FROM config_parking_setting
				WHERE setting_key = 'DEFAULT_JUKIR_SHARE_PERCENT'
				LIMIT 1
			), 0) / 100
		),

		FLOOR(
			ps.amount * COALESCE((
				SELECT CAST(setting_value AS UNSIGNED)
				FROM config_parking_setting
				WHERE setting_key = 'DEFAULT_JUKIR_SHARE_PERCENT'
				LIMIT 1
			), 0) / 100
		),

		0,
		0,

		ps.started_at,
		ps.paid_at,

		NOW(),
		NOW()

	FROM parking_session ps

	JOIN payment_transaction pt
		ON pt.parking_session_id = ps.id

	JOIN location_parking lp
		ON lp.id = ps.location_id

	WHERE pt.payment_code = ?

	ON DUPLICATE KEY UPDATE
		paid_at = VALUES(paid_at),
		updated_at = NOW();
	`

	if _, err := r.db.ExecContext(ctx, query, req.PaymentCode); err != nil {
		return fmt.Errorf("buat financial parking transaction: %w", err)
	}

	return nil
}

func (r *PaymentRepositoryImpl) InsertNotifikasiSuccess(ctx context.Context, req model.PaymentBusinessModel) error {
	query := `
	INSERT INTO notification_user (
		user_id,
		notification_type_id,
		title,
		body,
		data_json,
		is_read,
		read_at,
		created_at
	)
	SELECT
		target.user_id,
		(
			SELECT id
			FROM master_notification_type
			WHERE notification_type_code = 'PAYMENT'
			LIMIT 1
		),
		'Pembayaran Parkir Berhasil',
		CONCAT('Pembayaran parkir ', ps.plate_number, ' sebesar Rp ', ps.amount, ' berhasil.'),
		JSON_OBJECT(
			'sessionId', ps.id,
			'transactionCode', ps.transaction_code,
			'status', 'SUCCESS',
			'amount', ps.amount
		),
		0,
		NULL,
		NOW()

	FROM parking_session ps

	JOIN payment_transaction pt
		ON pt.parking_session_id = ps.id

	JOIN (
		SELECT ps1.customer_user_id AS user_id, ps1.id AS session_id
		FROM parking_session ps1
		WHERE ps1.customer_user_id IS NOT NULL

		UNION

		SELECT ps2.officer_user_id AS user_id, ps2.id AS session_id
		FROM parking_session ps2
	) target
		ON target.session_id = ps.id

	WHERE pt.payment_code = ?;
	`

	if _, err := r.db.ExecContext(ctx, query, req.PaymentCode); err != nil {
		return fmt.Errorf("insert notifikasi success: %w", err)
	}

	return nil
}

func (r *PaymentRepositoryImpl) UpdatePaymentTransactionFailed(ctx context.Context, req model.PaymentBusinessModel) error {
	query := `
	UPDATE payment_transaction pt

	JOIN master_payment_status mps
		ON mps.status_code = 'FAILED'

	SET
		pt.payment_status_id = mps.id,
		pt.failed_reason = ?,
		pt.updated_at = NOW()

	WHERE pt.payment_code = ?
	AND pt.paid_at IS NULL;
	`

	if _, err := r.db.ExecContext(ctx, query, req.FailedReason, req.PaymentCode); err != nil {
		return fmt.Errorf("update payment transaction failed: %w", err)
	}

	return nil
}

func (r *PaymentRepositoryImpl) UpdateParkingSessionFailed(ctx context.Context, req model.PaymentBusinessModel) error {
	query := `
	UPDATE parking_session ps

	JOIN payment_transaction pt
		ON pt.parking_session_id = ps.id

	JOIN master_payment_status mps
		ON mps.status_code = 'FAILED'

	SET
		ps.payment_status_id = mps.id,
		ps.updated_at = NOW()

	WHERE pt.payment_code = ?;
	`

	if _, err := r.db.ExecContext(ctx, query, req.PaymentCode); err != nil {
		return fmt.Errorf("update parking session failed: %w", err)
	}

	return nil
}

func (r *PaymentRepositoryImpl) InsertNotifikasiFailed(ctx context.Context, req model.PaymentBusinessModel) error {
	query := `
	INSERT INTO notification_user (
		user_id,
		notification_type_id,
		title,
		body,
		data_json,
		is_read,
		read_at,
		created_at
	)
	SELECT
		target.user_id,
		(
			SELECT id
			FROM master_notification_type
			WHERE notification_type_code = 'PAYMENT'
			LIMIT 1
		),
		'Pembayaran Parkir Gagal',
		CONCAT('Pembayaran parkir ', ps.plate_number, ' gagal.'),
		JSON_OBJECT(
			'sessionId', ps.id,
			'transactionCode', ps.transaction_code,
			'status', 'FAILED',
			'amount', ps.amount
		),
		0,
		NULL,
		NOW()

	FROM parking_session ps

	JOIN payment_transaction pt
		ON pt.parking_session_id = ps.id

	JOIN (
		SELECT ps1.customer_user_id AS user_id, ps1.id AS session_id
		FROM parking_session ps1
		WHERE ps1.customer_user_id IS NOT NULL

		UNION

		SELECT ps2.officer_user_id AS user_id, ps2.id AS session_id
		FROM parking_session ps2
	) target
		ON target.session_id = ps.id

	WHERE pt.payment_code = ?;
	`

	if _, err := r.db.ExecContext(ctx, query, req.PaymentCode); err != nil {
		return fmt.Errorf("insert notifikasi failed: %w", err)
	}

	return nil
}
