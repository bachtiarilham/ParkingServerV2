package riwayat

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	model "modulegue/internal/domain/mobile/model/riwayat"
	"modulegue/internal/domain/mobile/repository"
)

type RiwayatRepositoryImpl struct {
	db *sql.DB
}

func NewRiwayatRepositoryImpl(db *sql.DB) repository.RiwayatRepository {
	return &RiwayatRepositoryImpl{db: db}
}

func (r *RiwayatRepositoryImpl) GetRiwayat(ctx context.Context, req model.RiwayatRequestModel) (*model.RiwayatModel, error) {
	startDate := strings.TrimSpace(req.StartDate)
	endDate := strings.TrimSpace(req.EndDate)
	statusFilter := ""

	if startDate == "" {
		startDate = time.Now().Format("2006-01-02")
	}
	if endDate == "" {
		endDate = startDate
	}

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

			AND (
				? = 'SEMUA'
				OR mpm.payment_method_code = ?
			)

			AND (
				? = 'SEMUA'
				OR mvt.vehicle_type_code = ?
			)

			AND (
				? = 0
				OR fpt.location_id = ?
			)

		ORDER BY
			DATE(fpt.paid_at) DESC,
			fpt.paid_at DESC;
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		req.UserID, req.UserID, req.UserID,
		startDate, endDate,
		req.PaymentCode, req.PaymentCode,
		req.VehicleCode, req.VehicleCode,
		req.LokasiCode, req.LokasiCode,
		statusFilter, statusFilter,
		req.PaymentCode, req.PaymentCode,
		req.UserID, req.UserID, req.UserID, req.UserID,
	)
	if err != nil {
		return nil, fmt.Errorf("get riwayat: %w", err)
	}
	defer rows.Close()

	sections := []model.RiwayatSectionModel{}
	sectionIndexByDate := map[string]int{}

	for rows.Next() {
		var (
			sectionDate string
			code        string
			plateNumber string
			vehicleType string
			itemTime    string
			amount      int64
			isEntryInt  int
		)

		if err := rows.Scan(&sectionDate, &code, &plateNumber, &vehicleType, &itemTime, &amount, &isEntryInt); err != nil {
			return nil, fmt.Errorf("scan riwayat item: %w", err)
		}

		idx, exists := sectionIndexByDate[sectionDate]
		if !exists {
			dateCopy := sectionDate
			sections = append(sections, model.RiwayatSectionModel{
				Date:  &dateCopy,
				Items: []model.RiwayatItemModel{},
			})
			idx = len(sections) - 1
			sectionIndexByDate[sectionDate] = idx
		}

		codeCopy := code
		plateCopy := plateNumber
		vehicleCopy := vehicleType
		timeCopy := itemTime
		amountCopy := amount
		isEntry := isEntryInt == 1

		sections[idx].Items = append(sections[idx].Items, model.RiwayatItemModel{
			Code:        &codeCopy,
			PlateNumber: &plateCopy,
			VehicleType: &vehicleCopy,
			Time:        &timeCopy,
			Amount:      &amountCopy,
			IsEntry:     &isEntry,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate riwayat: %w", err)
	}

	return &model.RiwayatModel{
		Sections: sections,
	}, nil
}
