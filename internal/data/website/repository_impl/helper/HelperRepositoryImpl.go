package helper

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"modulegue/core/utils"
	model "modulegue/internal/domain/web/model/helper"
	"modulegue/internal/domain/web/repository"
)

type HelperRepositoryImpl struct {
	db *sql.DB
}

func NewHelperRepositoryImpl(db *sql.DB) repository.HelperRepository {
	return &HelperRepositoryImpl{db: db}
}

func (r *HelperRepositoryImpl) GetLokasi(ctx context.Context, reqModel model.GetLokasiRequestModel) (*model.GetLokasiResponseModel, error) {
	const query = `
SELECT
    lp.id AS ID,
    COALESCE(lp.location_name, '') AS NamaParlok,
    COALESCE(lp.address, '') AS JalanParlok
FROM location_parking lp
WHERE lp.is_active = 1
ORDER BY lp.location_name ASC;
`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get lokasi: %w", err)
	}
	defer rows.Close()

	items := make([]model.LokasiItemModel, 0)
	for rows.Next() {
		var (
			id          sql.NullInt64
			namaParlok  sql.NullString
			jalanParlok sql.NullString
		)

		if err := rows.Scan(&id, &namaParlok, &jalanParlok); err != nil {
			return nil, fmt.Errorf("scan lokasi: %w", err)
		}

		items = append(items, model.LokasiItemModel{
			ID:          int(utils.NullInt64Value(id)),
			NamaParlok:  utils.NullStringValue(namaParlok),
			JalanParlok: utils.NullStringValue(jalanParlok),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate lokasi: %w", err)
	}

	return &model.GetLokasiResponseModel{LokasiItem: &items}, nil
}

func (r *HelperRepositoryImpl) GetTarif(ctx context.Context, reqModel model.GetTarifRequestModel) (*model.GetTarifResponseModel, error) {
	idLokasi, err := strconv.Atoi(reqModel.IDLokasi)
	if err != nil {
		return nil, fmt.Errorf("id_lokasi tidak valid: %w", err)
	}

	const query = `
SELECT
    CASE
        WHEN mvt.vehicle_type_code IN ('MOTOR', 'MOTORCYCLE') THEN 'Motor'
        WHEN mvt.vehicle_type_code IN ('MOBIL', 'CAR') THEN 'Mobil'
        ELSE mvt.vehicle_type_name
    END AS ket_tarif,

    COALESCE(lt.id, 0) AS id,
    COALESCE(lt.tariff_amount, 0) AS tarif


FROM master_vehicle_type mvt

LEFT JOIN location_tariff lt
    ON lt.vehicle_type_id = mvt.id
   AND lt.location_id = ?
   AND lt.is_active = 1

WHERE mvt.is_active = 1
  AND mvt.vehicle_type_code IN ('MOTOR', 'MOTORCYCLE', 'MOBIL', 'CAR')

ORDER BY
    CASE
        WHEN mvt.vehicle_type_code IN ('MOTOR', 'MOTORCYCLE') THEN 1
        WHEN mvt.vehicle_type_code IN ('MOBIL', 'CAR') THEN 2
        ELSE 99
    END;
`

	rows, err := r.db.QueryContext(ctx, query, idLokasi)
	if err != nil {
		return nil, fmt.Errorf("get tarif by lokasi: %w", err)
	}
	defer rows.Close()

	items := make([]model.TarifItemModel, 0)
	for rows.Next() {
		var (
			ketTarif sql.NullString
			id       sql.NullInt64
			tarif    sql.NullInt64
		)

		if err := rows.Scan(&ketTarif, &id, &tarif); err != nil {
			return nil, fmt.Errorf("scan tarif: %w", err)
		}

		items = append(items, model.TarifItemModel{
			KetTarif: utils.NullStringValue(ketTarif),
			ID:       int(utils.NullInt64Value(id)),
			Tarif:    int(utils.NullInt64Value(tarif)),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tarif: %w", err)
	}

	return &model.GetTarifResponseModel{Tarif: &items}, nil
}
