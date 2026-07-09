package helper

import (
	"context"
	"database/sql"
	"fmt"

	model "modulegue/internal/domain/mobile/model/helper"
	"modulegue/internal/domain/mobile/repository"
)

type HelperRepositoryImpl struct {
	db *sql.DB
}

func NewHelperRepositoryImpl(db *sql.DB) repository.HelperRepository {
	return &HelperRepositoryImpl{db: db}
}

func (r *HelperRepositoryImpl) GetLokasi(ctx context.Context) (*model.LokasiModel, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT location_name
		FROM parking_location
		WHERE is_active = 1
		ORDER BY location_name ASC
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("get lokasi: %w", err)
	}
	defer rows.Close()

	result := &model.LokasiModel{
		Lokasi: []string{},
	}

	for rows.Next() {
		var lokasi string
		if err := rows.Scan(&lokasi); err != nil {
			return nil, fmt.Errorf("scan lokasi: %w", err)
		}
		result.Lokasi = append(result.Lokasi, lokasi)
	}

	return result, nil
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
