package parking

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"modulegue/core/utils"
	model "modulegue/internal/domain/mobile/model/parking"
	"modulegue/internal/domain/mobile/repository"
)

type ParkingRepositoryImpl struct {
	db *sql.DB
}

func NewParkingRepositoryImpl(db *sql.DB) repository.ParkingRepository {
	return &ParkingRepositoryImpl{db: db}
}

func (r *ParkingRepositoryImpl) PostParking(ctx context.Context, req model.PostParkingRequestModel) (*model.ParkingBusinessModel, error) {
	query :=
		`
	SELECT
		aoc.officer_user_id AS officerUserId,

		lp.id AS locationId,
		lp.location_name AS locationName,
		lp.address AS address,

		lp.zone_id AS zoneId,
		lz.zone_name AS zoneName,

		la.id AS areaId,
		la.area_name AS areaName,

		mvt.id AS vehicleTypeId,
		mvt.vehicle_type_code AS vehicleTypeCode,
		mvt.vehicle_type_name AS vehicleTypeName,

		lt.tariff_amount AS amount,

		(
			SELECT id
			FROM master_parking_status
			WHERE status_code = 'WAITING_PAYMENT'
			LIMIT 1
		) AS parkingStatusId,

		(
			SELECT id
			FROM master_payment_status
			WHERE status_code = 'PENDING'
			LIMIT 1
		) AS paymentStatusId

	FROM assignment_officer_current aoc

	JOIN location_parking lp
		ON lp.id = aoc.location_id
	AND lp.is_active = 1

	JOIN location_zone lz
		ON lz.id = lp.zone_id
	AND lz.is_active = 1

	JOIN location_area la
		ON la.id = COALESCE(?, aoc.area_id)
	AND la.location_id = lp.id
	AND la.is_active = 1

	JOIN master_vehicle_type mvt
		ON mvt.vehicle_type_code = ?
	AND mvt.is_active = 1

	JOIN location_tariff lt
		ON lt.location_id = lp.id
	AND lt.vehicle_type_id = mvt.id
	AND lt.is_active = 1
	AND (lt.effective_from IS NULL OR lt.effective_from <= NOW())
	AND (lt.effective_to IS NULL OR lt.effective_to >= NOW())

	WHERE aoc.officer_user_id = ?

	LIMIT 1;
	`

	var (
		businessModel model.ParkingBusinessModel
		parkingStatus sql.NullInt64
		paymentStatus sql.NullInt64
	)

	err := r.db.QueryRowContext(
		ctx,
		query,
		utils.NullInt64Param(req.SelectedAreaId),
		req.VehicleTypeCode,
		req.OfficerUserId,
	).Scan(
		&businessModel.OfficerUserId,
		&businessModel.LocationId,
		&businessModel.LocationName,
		&businessModel.Address,
		&businessModel.ZoneId,
		&businessModel.ZoneName,
		&businessModel.AreaId,
		&businessModel.AreaName,
		&businessModel.VehicleTypeId,
		&businessModel.VehicleTypeCode,
		&businessModel.VehicleTypeName,
		&businessModel.Amount,
		&parkingStatus,
		&paymentStatus,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("kebutuhan parking tidak ditemukan")
		}
		return nil, fmt.Errorf("post parking: %w", err)
	}
	return &businessModel, nil
}

func (r *ParkingRepositoryImpl) InsertParkingSession(ctx context.Context, req model.ParkingBusinessModel) (sessionId int64, err error) {
	query := `
	INSERT INTO parking_session (
		session_code,
		transaction_code,

		customer_user_id,
		officer_user_id,

		vehicle_id,
		vehicle_type_id,
		plate_number,

		zone_id,
		location_id,
		area_id,

		parking_status_id,
		payment_status_id,

		amount,
		qr_string,
		qr_expired_at,

		started_at,
		paid_at,
		finished_at,
		cancelled_at,

		created_at,
		updated_at
	)
	VALUES (
		?,
		?,

		NULL,
		?,

		NULL,
		?,
		UPPER(?),

		?,
		?,
		?,

		?,
		?,

		?,
		NULL,
		DATE_ADD(NOW(), INTERVAL 24 HOUR),

		NOW(),
		NULL,
		NULL,
		NULL,

		NOW(),
		NOW()
	);
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		req.SessionCode,
		req.TransactionCode,
		req.OfficerUserId,
		req.VehicleTypeId,
		req.PlateNumber,
		req.ZoneId,
		req.LocationId,
		req.AreaId,
		req.ParkingStatusId,
		req.PaymentStatusId,
		req.Amount,
	)
	if err != nil {
		return 0, fmt.Errorf("insert parking session: %w", err)
	}

	sessionId, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get parking session id: %w", err)
	}

	return sessionId, nil
}

func (r *ParkingRepositoryImpl) InsertQrisString(ctx context.Context, req model.ParkingBusinessModel) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin qris transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	updateQuery := `
	UPDATE parking_session
	SET
		qr_string = ?,
		qr_expired_at = DATE_ADD(NOW(), INTERVAL 15 MINUTE),
		updated_at = NOW()
	WHERE id = ?
	AND officer_user_id = ?;
	`

	if _, err = tx.ExecContext(ctx, updateQuery, req.QrisString, req.SessionId, req.OfficerUserId); err != nil {
		return fmt.Errorf("update parking session qris: %w", err)
	}

	insertQuery := `
	INSERT INTO payment_transaction (
		payment_code,
		parking_session_id,
		user_id,

		payment_method_id,
		payment_status_id,

		amount,
		external_reference,
		provider_name,

		paid_at,
		expired_at,
		failed_reason,

		created_at,
		updated_at
	)
	VALUES (
		?,
		?,
		NULL,

		(
			SELECT id
			FROM master_payment_method
			WHERE payment_method_code = 'QRIS'
			LIMIT 1
		),

		(
			SELECT id
			FROM master_payment_status
			WHERE status_code = 'PENDING'
			LIMIT 1
		),

		?,
		?,
		?,

		NULL,
		DATE_ADD(NOW(), INTERVAL 24 HOUR),
		NULL,

		NOW(),
		NOW()
	);
	`

	if _, err = tx.ExecContext(
		ctx,
		insertQuery,
		req.PaymentCode,
		req.SessionId,
		req.Amount,
		req.ExternalReference,
		req.ProviderName,
	); err != nil {
		return fmt.Errorf("insert payment transaction: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit qris transaction: %w", err)
	}

	return nil
}

func (r *ParkingRepositoryImpl) ReturnToApp(ctx context.Context, req model.ParkingBusinessModel) (*model.PostParkingResponseModel, error) {
	query := `
	SELECT
		ps.id AS sessionId,
		ps.session_code AS sessionCode,
		ps.transaction_code AS transactionCode,

		ps.plate_number AS plateNumber,

		mvt.vehicle_type_code AS vehicleTypeCode,
		mvt.vehicle_type_name AS vehicleTypeName,

		lz.id AS zoneId,
		lz.zone_name AS zoneName,

		lp.id AS locationId,
		lp.location_name AS locationName,
		lp.address AS address,

		la.id AS areaId,
		la.area_name AS areaName,

		ps.amount AS amount,
		ps.qr_string AS qrString,
		ps.qr_expired_at AS qrExpiredAt,

		pt.payment_code AS paymentCode,

		mps.status_code AS paymentStatusCode,
		mps.status_name AS paymentStatusName

	FROM parking_session ps

	JOIN master_vehicle_type mvt
		ON mvt.id = ps.vehicle_type_id

	JOIN location_zone lz
    	ON lz.id = ps.zone_id

	JOIN location_parking lp
		ON lp.id = ps.location_id

	LEFT JOIN location_area la
		ON la.id = ps.area_id

	JOIN payment_transaction pt
		ON pt.parking_session_id = ps.id

	JOIN master_payment_status mps
		ON mps.id = ps.payment_status_id

	WHERE ps.id = ?
	AND ps.officer_user_id = ?

	LIMIT 1;
	`

	var (
		resp                model.PostParkingResponseModel
		qrExpiredAt         time.Time
		qrExpiredAtNullable sql.NullTime
	)

	if err := r.db.QueryRowContext(ctx, query, req.SessionId, req.OfficerUserId).Scan(
		&resp.SessionId,
		&resp.SessionCode,
		&resp.TransactionCode,
		&resp.PlateNumber,
		&resp.VehicleTypeCode,
		&resp.VehicleTypeName,
		&resp.ZoneId,
		&resp.ZoneName,
		&resp.LocationId,
		&resp.LocationName,
		&resp.Address,
		&resp.AreaId,
		&resp.AreaName,
		&resp.Amount,
		&resp.QrString,
		&qrExpiredAtNullable,
		&resp.PaymentCode,
		&resp.PaymentStatusCode,
		&resp.PaymentStatusName,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("return parking response not found")
		}
		return nil, fmt.Errorf("return to app: %w", err)
	}

	if qrExpiredAtNullable.Valid {
		qrExpiredAt = qrExpiredAtNullable.Time
	}
	resp.QrExpiredAt = qrExpiredAt

	return &resp, nil
}
