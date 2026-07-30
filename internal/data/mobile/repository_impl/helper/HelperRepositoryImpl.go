package helper

import (
	"context"
	"database/sql"
	"fmt"

	"modulegue/core/utils"
	model "modulegue/internal/domain/mobile/model/helper"
	"modulegue/internal/domain/mobile/repository"
)

type HelperRepositoryImpl struct {
	db *sql.DB
}

func NewHelperRepositoryImpl(db *sql.DB) repository.HelperRepository {
	return &HelperRepositoryImpl{db: db}
}

func (r *HelperRepositoryImpl) GetLokasi(ctx context.Context, userId int64) (*model.LokasiModel, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT
			lz.id AS zoneId,
			lz.zone_name AS zoneName,

			lp.id AS locationId,
			lp.location_name AS locationName,
			lp.address AS address,

			la.id AS areaId,
			la.area_name AS areaName,

			CASE
				WHEN la.id = aoc.area_id THEN 1
				ELSE 0
			END AS isCurrentArea

		FROM assignment_officer_current aoc

		JOIN location_parking lp
			ON lp.id = aoc.location_id
		AND lp.is_active = 1

		JOIN location_zone lz
			ON lz.id = lp.zone_id
		AND lz.is_active = 1

		JOIN location_area la
			ON la.location_id = lp.id
		AND la.is_active = 1

		WHERE aoc.officer_user_id = ?

		ORDER BY
			isCurrentArea DESC,
			la.area_name ASC;
		`,
		userId,
	)
	if err != nil {
		return nil, fmt.Errorf("get lokasi: %w", err)
	}
	defer rows.Close()

	items := make([]model.LokasiItem, 0)

	for rows.Next() {
		var item model.LokasiItem
		if err := rows.Scan(
			&item.ZonaId,
			&item.NamaZona,
			&item.LokasiId,
			&item.NamaLokasi,
			&item.Address,
			&item.AreaId,
			&item.NamaArea,
			&item.IsCurrentArea,
		); err != nil {
			return nil, fmt.Errorf("scan lokasi: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate lokasi: %w", err)
	}

	return &model.LokasiModel{
		LokasiItem: &items,
	}, nil
}

func (r *HelperRepositoryImpl) GetTarif(ctx context.Context, userId int64) (*model.TarifModel, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT
			mvt.id AS vehicle_type_id,
			mvt.vehicle_type_code,
			mvt.vehicle_type_name,
			lt.tariff_amount

		FROM assignment_officer_current aoc

		JOIN location_tariff lt
			ON lt.location_id = aoc.location_id

		JOIN master_vehicle_type mvt
			ON mvt.id = lt.vehicle_type_id

		WHERE aoc.officer_user_id = ?
		AND lt.is_active = 1
		AND mvt.is_active = 1
		AND (lt.effective_from IS NULL OR lt.effective_from <= NOW())
		AND (lt.effective_to IS NULL OR lt.effective_to >= NOW())

		ORDER BY mvt.sort_order ASC, mvt.id ASC;
		`,
		userId,
	)
	if err != nil {
		return nil, fmt.Errorf("get tarif: %w", err)
	}
	defer rows.Close()

	items := make([]model.TarifItem, 0)
	for rows.Next() {
		var (
			vehicleTypeID   int64
			vehicleTypeCode string
			vehicleTypeName string
			tariffAmount    int64
		)

		if err := rows.Scan(
			&vehicleTypeID,
			&vehicleTypeCode,
			&vehicleTypeName,
			&tariffAmount,
		); err != nil {
			return nil, fmt.Errorf("scan tarif: %w", err)
		}

		items = append(items, model.TarifItem{
			KendaraanId:   vehicleTypeID,
			KendaraanKode: vehicleTypeCode,
			KendaraanNama: vehicleTypeName,
			Nominal:       tariffAmount,
		})
	}

	result := &model.TarifModel{
		TarifItem: &items,
	}

	return result, nil
}

func (r *HelperRepositoryImpl) GetNominalPayment(ctx context.Context) (*model.NominalPaymentModel, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT
			id AS optionId,
			nominal_amount AS nominalAmount,
			label AS label
		FROM payment_topup_option
		WHERE is_active = 1
		ORDER BY sort_order ASC, nominal_amount ASC;
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("get nominal tarif: %w", err)
	}
	defer rows.Close()

	items := make([]model.NominalItemModel, 0)
	for rows.Next() {
		var (
			optionID      int64
			nominalAmount int64
			label         sql.NullString
		)

		if err := rows.Scan(
			&optionID,
			&nominalAmount,
			&label,
		); err != nil {
			return nil, fmt.Errorf("scan tarif: %w", err)
		}

		items = append(items, model.NominalItemModel{
			OptionID:      optionID,
			NominalAmount: nominalAmount,
			Label:         utils.NullStringValue(label),
		})
	}

	rows1, err := r.db.QueryContext(
		ctx,
		`
		SELECT
			id AS paymentMethodId,
			payment_method_code AS paymentMethodCode,
			payment_method_name AS paymentMethodName,
			logo_url AS logoUrl
		FROM master_payment_method
		WHERE is_active = 1 
		AND payment_method_code NOT IN ('QRIS', 'CASH') -- Tambahkan baris ini
		ORDER BY id ASC;
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("Get Metode pembayaran: %w", err)
	}
	defer rows1.Close()

	result := &model.NominalPaymentModel{
		Nominal: &items,
	}

	return result, nil
}

func (r *HelperRepositoryImpl) GetPaymentMethod(ctx context.Context) (*model.PaymentMethodModel, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT
			id AS paymentMethodId,
			payment_method_code AS paymentMethodCode,
			payment_method_name AS paymentMethodName,
			logo_url AS logoUrl
		FROM master_payment_method
		WHERE is_active = 1 
		AND payment_method_code NOT IN ('QRIS', 'CASH', 'MANUAL') -- Tambahkan baris ini
		ORDER BY id ASC;
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("Get Metode pembayaran: %w", err)
	}
	defer rows.Close()

	items := make([]model.MetodeItemModel, 0)
	for rows.Next() {
		var (
			PaymentMethodId int64
			NamaPayment     sql.NullString
			CodePayment     sql.NullString
			LogoPayment     sql.NullString
		)

		if err := rows.Scan(
			&PaymentMethodId,
			&NamaPayment,
			&CodePayment,
			&LogoPayment,
		); err != nil {
			return nil, fmt.Errorf("scan metode pembayaran: %w", err)
		}

		items = append(items, model.MetodeItemModel{
			PaymentMethodId: PaymentMethodId,
			NamaPayment:     utils.NullStringValue(NamaPayment),
			CodePayment:     utils.NullStringValue(CodePayment),
			LogoPayment:     utils.NullStringValue(LogoPayment),
		})

	}

	result := &model.PaymentMethodModel{
		MetodeBayar: &items,
	}

	return result, nil
}
