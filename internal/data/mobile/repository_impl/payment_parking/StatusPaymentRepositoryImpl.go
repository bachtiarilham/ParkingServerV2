package payment

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"modulegue/core/utils"
	model "modulegue/internal/domain/mobile/model/payment_parking"
	"modulegue/internal/domain/mobile/repository"
)

type StatusPaymentRepositoryImpl struct {
	db *sql.DB
}

func NewStatusPaymentRepositoryImpl(db *sql.DB) repository.StatusPaymentRepository {
	return &StatusPaymentRepositoryImpl{db: db}
}

func (r *StatusPaymentRepositoryImpl) GetPembayaranStatus(ctx context.Context, sessionId string) (*model.PostPaymentParkingResponseModel, error) {
	sessionID, err := strconv.ParseInt(strings.TrimSpace(sessionId), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}

	query := `
		SELECT
			ps.id AS sessionId,
			ps.session_code AS sessionCode,
			ps.transaction_code AS transactionCode,

			ps.plate_number AS plateNumber,

			mvt.vehicle_type_code AS vehicleTypeCode,
			mvt.vehicle_type_name AS vehicleTypeName,

			lp.id AS locationId,
			lp.location_name AS locationName,

			la.id AS areaId,
			la.area_name AS areaName,

			ps.amount AS amount,

			mps.status_code AS parkingStatusCode,
			mps.status_name AS parkingStatusName,

			mpys.status_code AS paymentStatusCode,
			mpys.status_name AS paymentStatusName,

			pt.payment_code AS paymentCode,
			pt.failed_reason AS failedReason,

			pr.receipt_number AS receiptNumber,

			ps.started_at AS startedAt,
			ps.paid_at AS paidAt,
			ps.qr_expired_at AS qrExpiredAt

		FROM parking_session ps

		JOIN master_vehicle_type mvt
			ON mvt.id = ps.vehicle_type_id

		JOIN location_parking lp
			ON lp.id = ps.location_id

		LEFT JOIN location_area la
			ON la.id = ps.area_id

		JOIN master_parking_status mps
			ON mps.id = ps.parking_status_id

		JOIN master_payment_status mpys
			ON mpys.id = ps.payment_status_id

		LEFT JOIN payment_transaction pt
			ON pt.parking_session_id = ps.id

		LEFT JOIN parking_receipt pr
			ON pr.parking_session_id = ps.id

		WHERE ps.id = ?

		LIMIT 1;
		`

	var (
		resp            model.PostPaymentParkingResponseModel
		areaID          sql.NullInt64
		receiptNumber   sql.NullString
		sessionCode     sql.NullString
		transactionCode sql.NullString
		plateNumber     sql.NullString
		vehicleCode     sql.NullString
		vehicleName     sql.NullString
		locationName    sql.NullString
		failedReason    sql.NullString
		paymentCode     sql.NullString
		startedAt       sql.NullTime
		paidAt          sql.NullTime
		qrExpiredAt     sql.NullTime
	)

	if err := r.db.QueryRowContext(ctx, query, sessionID).Scan(
		&resp.SessionId,
		&sessionCode,
		&transactionCode,
		&plateNumber,
		&vehicleCode,
		&vehicleName,
		&resp.LocationId,
		&locationName,
		&areaID,
		&resp.AreaName,
		&resp.Amount,
		&resp.ParkingStatusCode,
		&resp.ParkingStatusName,
		&resp.PaymentStatusCode,
		&resp.PaymentStatusName,
		&paymentCode,
		&failedReason,
		&receiptNumber,
		&startedAt,
		&paidAt,
		&qrExpiredAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get pembayaran status: %w", err)
	}

	resp.SessionCode = utils.NullStringValue(sessionCode)
	resp.TransactionCode = utils.NullStringValue(transactionCode)
	resp.PlateNumber = utils.NullStringValue(plateNumber)
	resp.VehicleTypeCode = utils.NullStringValue(vehicleCode)
	resp.VehicleTypeName = utils.NullStringValue(vehicleName)
	resp.LocationName = utils.NullStringValue(locationName)
	resp.AreaId = utils.NullInt64Value(areaID)
	resp.PaymentCode = utils.NullStringValue(paymentCode)
	resp.FailedReason = utils.NullStringValue(failedReason)
	resp.ReceiptNumber = parseReceiptNumber(receiptNumber)
	resp.StartedAt = utils.NullTimeValue(startedAt)
	resp.PaidAt = utils.NullTimeValue(paidAt)
	resp.QrExpiredAt = utils.NullTimeValue(qrExpiredAt)

	return &resp, nil
}

func parseReceiptNumber(v sql.NullString) int64 {
	if !v.Valid {
		return 0
	}
	n, err := strconv.ParseInt(strings.TrimSpace(v.String), 10, 64)
	if err != nil {
		return 0
	}
	return n
}
