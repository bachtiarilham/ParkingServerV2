package riwayat

import (
	"context"
	"database/sql"
	"fmt"

	"modulegue/core/utils"
	model "modulegue/internal/domain/mobile/model/riwayat"
	"modulegue/internal/domain/mobile/repository"
)

type RiwayatRepositoryImpl struct {
	db *sql.DB
}

type topUpRow struct {
	SectionDate        sql.NullString
	TopUpTransactionID sql.NullInt64
	TopUpCode          sql.NullString
	UserID             sql.NullInt64
	WalletID           sql.NullInt64
	PaymentMethodID    sql.NullInt64
	PaymentMethodCode  sql.NullString
	PaymentMethodName  sql.NullString
	Amount             sql.NullInt64
	AdminFee           sql.NullInt64
	TotalAmount        sql.NullInt64
	TransactionStatus  sql.NullString
	ExternalReference  sql.NullString
	ProviderName       sql.NullString
	CreatedAt          sql.NullTime
	ExpiredAt          sql.NullTime
	PaidAt             sql.NullTime
	CompletedAt        sql.NullTime
	FailedReason       sql.NullString
}

func mapTopUpRowToItem(row topUpRow) model.TopUpItemModel {
	return model.TopUpItemModel{
		Code:              utils.NullStringValue(row.TopUpCode),
		PaymentMethodName: utils.NullStringValue(row.PaymentMethodName),
		TransactionStatus: utils.NullStringValue(row.TransactionStatus),
		ProviderName:      utils.NullStringValue(row.ProviderName),
		Time:              utils.FormatRiwayatTime(row.PaidAt, row.CompletedAt, row.CreatedAt),
		Amount:            utils.NullInt64Value(row.Amount),
	}
}

func NewRiwayatRepositoryImpl(db *sql.DB) repository.RiwayatRepository {
	return &RiwayatRepositoryImpl{db: db}
}

func (r *RiwayatRepositoryImpl) GetRiwayat(ctx context.Context, req model.RiwayatRequestModel) (*model.RiwayatModel, error) {
	parkir, err := r.GetParkirRiwayat(ctx, req)
	if err != nil {
		return nil, err
	}

	transaksi, err := r.GetTransaksiRiwayat(ctx, req)
	if err != nil {
		return nil, err
	}

	return &model.RiwayatModel{
		ParkirSections: parkir,
		TopUpSections:  transaksi,
	}, nil
}

func (r *RiwayatRepositoryImpl) GetParkirRiwayat(ctx context.Context, req model.RiwayatRequestModel) ([]model.RiwayatSectionModel, error) {
	query := `
		SELECT
			DATE(fpt.paid_at) AS tanggal,

			fpt.id AS transactionId,
			fpt.transaction_code AS transactionCode,

			fpt.session_id AS sessionId,

			fpt.plate_number AS plateNumber,

			mvt.id AS vehicleTypeId,
			mvt.vehicle_type_code AS vehicleTypeCode,
			mvt.vehicle_type_name AS vehicleTypeName,

			mpm.id AS paymentMethodId,
			mpm.payment_method_code AS paymentMethodCode,
			mpm.payment_method_name AS paymentMethodName,

			lp.id AS locationId,
			lp.location_name AS locationName,
			lp.address AS locationAddress,

			la.id AS areaId,
			la.area_name AS areaName,

			lz.id AS zoneId,
			lz.zone_name AS zoneName,

			fpt.base_amount AS baseAmount,
			fpt.discount_amount AS discountAmount,
			fpt.final_amount AS finalAmount,
			fpt.company_share AS companyShare,
			fpt.jukir_share AS jukirShare,
			fpt.tax_amount AS taxAmount,
			fpt.fee_amount AS feeAmount,

			fpt.transaction_status AS transactionStatus,
			fpt.operation_type AS operationType,

			fpt.occurred_at AS occurredAt,
			fpt.paid_at AS paidAt,
			fpt.created_at AS createdAt

		FROM financial_parking_transaction fpt

		JOIN master_vehicle_type mvt
			ON mvt.id = fpt.vehicle_type_id

		JOIN master_payment_method mpm
			ON mpm.id = fpt.payment_method_id

		JOIN location_parking lp
			ON lp.id = fpt.location_id

		LEFT JOIN location_area la
			ON la.id = fpt.area_id

		LEFT JOIN location_zone lz
			ON lz.id = fpt.zone_id

		WHERE
			(
				fpt.customer_user_id = ?
				OR fpt.jukir_user_id = ?
				OR fpt.officer_user_id = ?
			)

			AND fpt.transaction_status = 'SUCCESS'

			AND fpt.paid_at >= DATE(?)
			AND fpt.paid_at < DATE_ADD(DATE(?), INTERVAL 1 DAY)

		ORDER BY
			DATE(fpt.paid_at) DESC,
			fpt.paid_at DESC;
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		req.UserID, req.UserID, req.UserID,
		req.StartDate, req.EndDate,
	)
	if err != nil {
		return nil, fmt.Errorf("get riwayat: %w", err)
	}
	defer rows.Close()

	sections := make([]model.RiwayatSectionModel, 0)
	sectionIndexByDate := map[string]int{}

	for rows.Next() {
		var row riwayatRow
		if err := rows.Scan(
			&row.SectionDate,
			&row.TransactionID,
			&row.TransactionCode,
			&row.SessionID,
			&row.PlateNumber,
			&row.VehicleTypeID,
			&row.VehicleTypeCode,
			&row.VehicleTypeName,
			&row.PaymentMethodID,
			&row.PaymentMethodCode,
			&row.PaymentMethodName,
			&row.LocationID,
			&row.LocationName,
			&row.LocationAddress,
			&row.AreaID,
			&row.AreaName,
			&row.ZoneID,
			&row.ZoneName,
			&row.BaseAmount,
			&row.DiscountAmount,
			&row.FinalAmount,
			&row.CompanyShare,
			&row.JukirShare,
			&row.TaxAmount,
			&row.FeeAmount,
			&row.TransactionStatus,
			&row.OperationType,
			&row.OccurredAt,
			&row.PaidAt,
			&row.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan riwayat item: %w", err)
		}

		dateKey := mapRiwayatSectionDate(row)
		idx, exists := sectionIndexByDate[dateKey]
		if !exists {
			sections = append(sections, model.RiwayatSectionModel{
				Date:  dateKey,
				Items: []model.RiwayatItemModel{},
			})
			idx = len(sections) - 1
			sectionIndexByDate[dateKey] = idx
		}

		sections[idx].Items = append(sections[idx].Items, mapRiwayatRowToItem(row))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate riwayat: %w", err)
	}

	return sections, nil
}

func (r *RiwayatRepositoryImpl) GetTransaksiRiwayat(ctx context.Context, req model.RiwayatRequestModel) ([]model.TopUpSectionModel, error) {
	query := `
			SELECT
		DATE(COALESCE(ptt.paid_at, ptt.completed_at, ptt.created_at)) AS tanggal,

		ptt.id AS topupTransactionId,
		ptt.topup_code AS topupCode,
		ptt.user_id AS userId,
		ptt.wallet_id AS walletId,

		mpm.id AS paymentMethodId,
		mpm.payment_method_code AS paymentMethodCode,
		mpm.payment_method_name AS paymentMethodName,

		ptt.amount AS amount,
		ptt.admin_fee AS adminFee,
		ptt.total_amount AS totalAmount,

		ptt.transaction_status AS transactionStatus,
		ptt.external_reference AS externalReference,
		ptt.provider_name AS providerName,

		ptt.created_at AS createdAt,
		ptt.expired_at AS expiredAt,
		ptt.paid_at AS paidAt,
		ptt.completed_at AS completedAt,
		ptt.failed_reason AS failedReason

	FROM payment_topup_transaction ptt

	JOIN master_payment_method mpm
		ON mpm.id = ptt.payment_method_id

	WHERE
		COALESCE(ptt.paid_at, ptt.completed_at, ptt.created_at) >= DATE(?)
		AND COALESCE(ptt.paid_at, ptt.completed_at, ptt.created_at) < DATE_ADD(DATE(?), INTERVAL 1 DAY)

	ORDER BY
		DATE(COALESCE(ptt.paid_at, ptt.completed_at, ptt.created_at)) DESC,
		ptt.created_at DESC;
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		req.StartDate, req.EndDate,
	)
	if err != nil {
		return nil, fmt.Errorf("get riwayat: %w", err)
	}
	defer rows.Close()

	sections := make([]model.TopUpSectionModel, 0)
	sectionIndexByDate := map[string]int{}

	for rows.Next() {
		var row topUpRow
		if err := rows.Scan(
			&row.SectionDate,
			&row.TopUpTransactionID,
			&row.TopUpCode,
			&row.UserID,
			&row.WalletID,
			&row.PaymentMethodID,
			&row.PaymentMethodCode,
			&row.PaymentMethodName,
			&row.Amount,
			&row.AdminFee,
			&row.TotalAmount,
			&row.TransactionStatus,
			&row.ExternalReference,
			&row.ProviderName,
			&row.CreatedAt,
			&row.ExpiredAt,
			&row.PaidAt,
			&row.CompletedAt,
			&row.FailedReason,
		); err != nil {
			return nil, fmt.Errorf("scan riwayat item: %w", err)
		}

		dateKey := utils.NullStringValue(row.SectionDate)
		idx, exists := sectionIndexByDate[dateKey]
		if !exists {
			sections = append(sections, model.TopUpSectionModel{
				Date:  dateKey,
				Items: []model.TopUpItemModel{},
			})
			idx = len(sections) - 1
			sectionIndexByDate[dateKey] = idx
		}

		sections[idx].Items = append(sections[idx].Items, mapTopUpRowToItem(row))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate riwayat: %w", err)
	}

	return sections, nil
}

func (r *RiwayatRepositoryImpl) GetParkirDetil(ctx context.Context, req model.DetilParkirRequestModel) (*model.DetilParkirModel, error) {
	query := `
	SELECT
		DATE(fpt.paid_at) AS tanggal,

		fpt.id AS transactionId,
		fpt.transaction_code AS transactionCode,
		fpt.session_id AS sessionId,
		fpt.plate_number AS plateNumber,

		mvt.id AS vehicleTypeId,
		mvt.vehicle_type_code AS vehicleTypeCode,
		mvt.vehicle_type_name AS vehicleTypeName,

		mpm.id AS paymentMethodId,
		mpm.payment_method_code AS paymentMethodCode,
		mpm.payment_method_name AS paymentMethodName,

		lp.id AS locationId,
		lp.location_name AS locationName,
		lp.address AS locationAddress,

		COALESCE(la.id, 0) AS areaId,
		COALESCE(la.area_name, '') AS areaName,

		COALESCE(lz.id, 0) AS zoneId,
		COALESCE(lz.zone_name, '') AS zoneName,

		fpt.base_amount AS baseAmount,
		fpt.discount_amount AS discountAmount,
		fpt.final_amount AS finalAmount,
		fpt.company_share AS companyShare,
		fpt.jukir_share AS jukirShare,
		fpt.tax_amount AS taxAmount,
		fpt.fee_amount AS feeAmount,

		fpt.transaction_status AS transactionStatus,
		fpt.operation_type AS operationType,

		fpt.occurred_at AS occurredAt,
		fpt.paid_at AS paidAt,
		fpt.created_at AS createdAt

	FROM financial_parking_transaction fpt

	JOIN master_vehicle_type mvt
		ON mvt.id = fpt.vehicle_type_id

	JOIN master_payment_method mpm
		ON mpm.id = fpt.payment_method_id

	JOIN location_parking lp
		ON lp.id = fpt.location_id

	LEFT JOIN location_area la
		ON la.id = fpt.area_id

	LEFT JOIN location_zone lz
		ON lz.id = fpt.zone_id

	WHERE
		fpt.transaction_code = ?;
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		req.TransactionCode,
	)
	if err != nil {
		return nil, fmt.Errorf("get riwayat: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var row riwayatRow
		if err := rows.Scan(
			&row.SectionDate,
			&row.TransactionID,
			&row.TransactionCode,
			&row.SessionID,
			&row.PlateNumber,
			&row.VehicleTypeID,
			&row.VehicleTypeCode,
			&row.VehicleTypeName,
			&row.PaymentMethodID,
			&row.PaymentMethodCode,
			&row.PaymentMethodName,
			&row.LocationID,
			&row.LocationName,
			&row.LocationAddress,
			&row.AreaID,
			&row.AreaName,
			&row.ZoneID,
			&row.ZoneName,
			&row.BaseAmount,
			&row.DiscountAmount,
			&row.FinalAmount,
			&row.CompanyShare,
			&row.JukirShare,
			&row.TaxAmount,
			&row.FeeAmount,
			&row.TransactionStatus,
			&row.OperationType,
			&row.OccurredAt,
			&row.PaidAt,
			&row.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan detil parkir: %w", err)
		}

		return &model.DetilParkirModel{
			Tanggal:           utils.NullStringValue(row.SectionDate),
			TransactionID:     utils.NullInt64Value(row.TransactionID),
			TransactionCode:   utils.NullStringValue(row.TransactionCode),
			SessionID:         utils.NullInt64Value(row.SessionID),
			PlateNumber:       utils.NullStringValue(row.PlateNumber),
			VehicleTypeID:     utils.NullInt64Value(row.VehicleTypeID),
			VehicleTypeCode:   utils.NullStringValue(row.VehicleTypeCode),
			VehicleTypeName:   utils.NullStringValue(row.VehicleTypeName),
			PaymentMethodID:   utils.NullInt64Value(row.PaymentMethodID),
			PaymentMethodCode: utils.NullStringValue(row.PaymentMethodCode),
			PaymentMethodName: utils.NullStringValue(row.PaymentMethodName),
			LocationID:        utils.NullInt64Value(row.LocationID),
			LocationName:      utils.NullStringValue(row.LocationName),
			LocationAddress:   utils.NullStringValue(row.LocationAddress),
			AreaID:            utils.NullInt64Value(row.AreaID),
			AreaName:          utils.NullStringValue(row.AreaName),
			ZoneID:            utils.NullInt64Value(row.ZoneID),
			ZoneName:          utils.NullStringValue(row.ZoneName),
			BaseAmount:        utils.NullInt64Value(row.BaseAmount),
			DiscountAmount:    utils.NullInt64Value(row.DiscountAmount),
			FinalAmount:       utils.NullInt64Value(row.FinalAmount),
			CompanyShare:      utils.NullInt64Value(row.CompanyShare),
			JukirShare:        utils.NullInt64Value(row.JukirShare),
			TaxAmount:         utils.NullInt64Value(row.TaxAmount),
			FeeAmount:         utils.NullInt64Value(row.FeeAmount),
			TransactionStatus: utils.NullStringValue(row.TransactionStatus),
			OperationType:     utils.NullStringValue(row.OperationType),
			OccurredAt:        utils.FormatRiwayatTime(row.OccurredAt),
			PaidAt:            utils.FormatRiwayatTime(row.PaidAt),
			CreatedAt:         utils.FormatRiwayatTime(row.CreatedAt),
		}, nil
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate riwayat: %w", err)
	}

	return nil, sql.ErrNoRows
}

func (r *RiwayatRepositoryImpl) GetTransaksiDetil(ctx context.Context, req model.DetilTransaksiRequestModel) (*model.DetilTransaksiModel, error) {
	query := `
	SELECT
		DATE(COALESCE(ptt.paid_at, ptt.completed_at, ptt.created_at)) AS tanggal,

		ptt.id AS topupTransactionId,
		ptt.topup_code AS topupCode,
		ptt.user_id AS userId,
		ptt.wallet_id AS walletId,

		mpm.id AS paymentMethodId,
		mpm.payment_method_code AS paymentMethodCode,
		mpm.payment_method_name AS paymentMethodName,

		ptt.amount AS amount,
		ptt.admin_fee AS adminFee,
		ptt.total_amount AS totalAmount,

		ptt.transaction_status AS transactionStatus,
		ptt.external_reference AS externalReference,
		ptt.provider_name AS providerName,

		ptt.created_at AS createdAt,
		ptt.expired_at AS expiredAt,
		ptt.paid_at AS paidAt,
		ptt.completed_at AS completedAt,
		ptt.failed_reason AS failedReason

	FROM payment_topup_transaction ptt

	JOIN master_payment_method mpm
		ON mpm.id = ptt.payment_method_id

	WHERE
		ptt.topup_code = ?;
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		req.TopUpCode,
	)
	if err != nil {
		return nil, fmt.Errorf("get riwayat: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var row topUpRow
		if err := rows.Scan(
			&row.SectionDate,
			&row.TopUpTransactionID,
			&row.TopUpCode,
			&row.UserID,
			&row.WalletID,
			&row.PaymentMethodID,
			&row.PaymentMethodCode,
			&row.PaymentMethodName,
			&row.Amount,
			&row.AdminFee,
			&row.TotalAmount,
			&row.TransactionStatus,
			&row.ExternalReference,
			&row.ProviderName,
			&row.CreatedAt,
			&row.ExpiredAt,
			&row.PaidAt,
			&row.CompletedAt,
			&row.FailedReason,
		); err != nil {
			return nil, fmt.Errorf("scan detil transaksi: %w", err)
		}

		return &model.DetilTransaksiModel{
			Tanggal:            utils.NullStringValue(row.SectionDate),
			TopUpTransactionID: utils.NullInt64Value(row.TopUpTransactionID),
			TopUpCode:          utils.NullStringValue(row.TopUpCode),
			UserID:             utils.NullInt64Value(row.UserID),
			WalletID:           utils.NullInt64Value(row.WalletID),
			PaymentMethodID:    utils.NullInt64Value(row.PaymentMethodID),
			PaymentMethodCode:  utils.NullStringValue(row.PaymentMethodCode),
			PaymentMethodName:  utils.NullStringValue(row.PaymentMethodName),
			Amount:             utils.NullInt64Value(row.Amount),
			AdminFee:           utils.NullInt64Value(row.AdminFee),
			TotalAmount:        utils.NullInt64Value(row.TotalAmount),
			TransactionStatus:  utils.NullStringValue(row.TransactionStatus),
			ExternalReference:  utils.NullStringValue(row.ExternalReference),
			ProviderName:       utils.NullStringValue(row.ProviderName),
			CreatedAt:          utils.FormatRiwayatTime(row.CreatedAt),
			ExpiredAt:          utils.FormatRiwayatTime(row.ExpiredAt),
			PaidAt:             utils.FormatRiwayatTime(row.PaidAt),
			CompletedAt:        utils.FormatRiwayatTime(row.CompletedAt),
			FailedReason:       utils.NullStringValue(row.FailedReason),
		}, nil
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate riwayat: %w", err)
	}

	return nil, sql.ErrNoRows
}
