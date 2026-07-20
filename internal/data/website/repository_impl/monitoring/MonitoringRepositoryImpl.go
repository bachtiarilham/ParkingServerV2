package monitoring

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"modulegue/core/utils"
	model "modulegue/internal/domain/web/model/monitoring"
)

type MonitoringRepositoryImpl struct {
	db *sql.DB
}

func NewMonitoringRepositoryImpl(db *sql.DB) *MonitoringRepositoryImpl {
	return &MonitoringRepositoryImpl{db: db}
}

func (r *MonitoringRepositoryImpl) GetMonitoring(ctx context.Context, reqModel model.MonitoringRequestModel) (*model.MonitoringResponseModel, error) {
	parlok, err := r.getMonitoringParlok(ctx, reqModel)
	if err != nil {
		return nil, err
	}

	transaksi, err := r.getMonitoringTransaksi(ctx, reqModel)
	if err != nil {
		return nil, err
	}

	return &model.MonitoringResponseModel{
		Parlok:    parlok,
		Transaksi: transaksi,
	}, nil
}

type monitoringParlokRow struct {
	NamaParlok      string
	IDZona          int
	NamaZona        string
	LatMin          string
	LatMax          string
	LngMin          string
	LngMax          string
	CenterX         string
	CenterY         string
	Altitude        string
	PendapatanMotor int
	PendapatanMobil int
	TotalPendapatan int
}

func (r *MonitoringRepositoryImpl) getMonitoringParlok(ctx context.Context, reqModel model.MonitoringRequestModel) (*[]model.ParlokItemModel, error) {
	const query = `
WITH filter_params AS (
    SELECT
        DATE(?) AS tgl_awal,
        DATE(?) AS tgl_akhir,
        ? AS id_lokasi,
        COALESCE(NULLIF(TRIM(?), ''), '') AS nama_petugas
)

SELECT
    COALESCE(lp.location_name, '') AS nama_parlok,

    COALESCE(lz.id, 0) AS id_zona,
    COALESCE(lz.zone_name, '') AS nama_zona,

    COALESCE(CAST(lp.min_latitude AS CHAR), '') AS lat_min,
    COALESCE(CAST(lp.max_latitude AS CHAR), '') AS lat_max,
    COALESCE(CAST(lp.min_longitude AS CHAR), '') AS lng_min,
    COALESCE(CAST(lp.max_longitude AS CHAR), '') AS lng_max,
    COALESCE(CAST(lp.center_latitude AS CHAR), '') AS center_x,
    COALESCE(CAST(lp.center_longitude AS CHAR), '') AS center_y,

    '' AS altitude,

    COALESCE(SUM(
        CASE
            WHEN mvt.vehicle_type_code IN ('MOTOR', 'MOTORCYCLE')
            THEN fpt.final_amount
            ELSE 0
        END
    ), 0) AS pendapatan_motor,

    COALESCE(SUM(
        CASE
            WHEN mvt.vehicle_type_code IN ('MOBIL', 'CAR')
            THEN fpt.final_amount
            ELSE 0
        END
    ), 0) AS pendapatan_mobil,

    COALESCE(SUM(fpt.final_amount), 0) AS total_pendapatan

FROM location_parking lp

JOIN location_zone lz
    ON lz.id = lp.zone_id

CROSS JOIN filter_params fp

LEFT JOIN financial_parking_transaction fpt
    ON fpt.location_id = lp.id
   AND fpt.transaction_status = 'SUCCESS'
   AND fpt.paid_at >= fp.tgl_awal
   AND fpt.paid_at < DATE_ADD(fp.tgl_akhir, INTERVAL 1 DAY)
   AND (
        fp.nama_petugas = ''
        OR EXISTS (
            SELECT 1
            FROM user_identity uif
            WHERE uif.id = fpt.jukir_user_id
              AND uif.full_name LIKE CONCAT('%', fp.nama_petugas, '%')
        )
   )

LEFT JOIN master_vehicle_type mvt
    ON mvt.id = fpt.vehicle_type_id

WHERE lp.is_active = 1
  AND (
        fp.id_lokasi = 0
        OR lp.id = fp.id_lokasi
  )

GROUP BY
    lp.id,
    lp.location_name,
    lz.id,
    lz.zone_name,
    lp.min_latitude,
    lp.max_latitude,
    lp.min_longitude,
    lp.max_longitude,
    lp.center_latitude,
    lp.center_longitude

ORDER BY
    lp.location_name ASC;
`

	rows, err := r.db.QueryContext(ctx, query, reqModel.TglAwal, reqModel.TglAkhir, reqModel.IDLokasi, reqModel.NamaPetugas)
	if err != nil {
		return nil, fmt.Errorf("get monitoring parlok: %w", err)
	}
	defer rows.Close()

	items := make([]model.ParlokItemModel, 0)
	for rows.Next() {
		var (
			namaParlok      sql.NullString
			idZona          sql.NullInt64
			namaZona        sql.NullString
			latMin          sql.NullString
			latMax          sql.NullString
			lngMin          sql.NullString
			lngMax          sql.NullString
			centerX         sql.NullString
			centerY         sql.NullString
			altitude        sql.NullString
			pendapatanMotor sql.NullInt64
			pendapatanMobil sql.NullInt64
			totalPendapatan sql.NullInt64
		)

		if err := rows.Scan(
			&namaParlok,
			&idZona,
			&namaZona,
			&latMin,
			&latMax,
			&lngMin,
			&lngMax,
			&centerX,
			&centerY,
			&altitude,
			&pendapatanMotor,
			&pendapatanMobil,
			&totalPendapatan,
		); err != nil {
			return nil, fmt.Errorf("scan monitoring parlok: %w", err)
		}

		items = append(items, model.ParlokItemModel{
			NamaParlok:      utils.NullStringValue(namaParlok),
			IDZona:          int(utils.NullInt64Value(idZona)),
			NamaZona:        utils.NullStringValue(namaZona),
			LatMin:          utils.NullStringValue(latMin),
			LatMax:          utils.NullStringValue(latMax),
			LngMin:          utils.NullStringValue(lngMin),
			LngMax:          utils.NullStringValue(lngMax),
			CenterX:         utils.NullStringValue(centerX),
			CenterY:         utils.NullStringValue(centerY),
			Altitude:        strings.TrimSpace(utils.NullStringValue(altitude)),
			PendapatanMotor: int(utils.NullInt64Value(pendapatanMotor)),
			PendapatanMobil: int(utils.NullInt64Value(pendapatanMobil)),
			TotalPendapatan: int(utils.NullInt64Value(totalPendapatan)),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate monitoring parlok: %w", err)
	}

	return &items, nil
}

func (r *MonitoringRepositoryImpl) getMonitoringTransaksi(ctx context.Context, reqModel model.MonitoringRequestModel) (*[]model.TransaksiItemModel, error) {
	const query = `
WITH filter_params AS (
    SELECT
        DATE(?) AS tgl_awal,
        DATE(?) AS tgl_akhir,
        ? AS id_lokasi,
        COALESCE(NULLIF(TRIM(?), ''), '') AS nama_petugas
)

SELECT
    COALESCE(ui.full_name, '') AS nama_jukir,
    COALESCE(lp.location_name, '') AS parlok,
    COALESCE(lz.zone_name, '') AS zona,

    COALESCE(fpt.plate_number, '') AS plat,

    DATE_FORMAT(fpt.paid_at, '%d/%m/%Y %H:%i:%s') AS waktu,

    COALESCE(mvt.vehicle_type_name, '') AS kendaraan,
    COALESCE(mpm.payment_method_name, '') AS pembayaran,

    COALESCE(fpt.final_amount, 0) AS tarif

FROM financial_parking_transaction fpt

CROSS JOIN filter_params fp

LEFT JOIN user_identity ui
    ON ui.id = fpt.jukir_user_id

LEFT JOIN location_parking lp
    ON lp.id = fpt.location_id

LEFT JOIN location_zone lz
    ON lz.id = fpt.zone_id

LEFT JOIN master_vehicle_type mvt
    ON mvt.id = fpt.vehicle_type_id

LEFT JOIN master_payment_method mpm
    ON mpm.id = fpt.payment_method_id

WHERE fpt.transaction_status = 'SUCCESS'
  AND fpt.paid_at >= fp.tgl_awal
  AND fpt.paid_at < DATE_ADD(fp.tgl_akhir, INTERVAL 1 DAY)

  AND (
        fp.id_lokasi = 0
        OR fpt.location_id = fp.id_lokasi
  )

  AND (
        fp.nama_petugas = ''
        OR ui.full_name LIKE CONCAT('%', fp.nama_petugas, '%')
  )

ORDER BY
    fpt.paid_at DESC,
    fpt.id DESC;
`

	rows, err := r.db.QueryContext(ctx, query, reqModel.TglAwal, reqModel.TglAkhir, reqModel.IDLokasi, reqModel.NamaPetugas)
	if err != nil {
		return nil, fmt.Errorf("get monitoring transaksi: %w", err)
	}
	defer rows.Close()

	items := make([]model.TransaksiItemModel, 0)
	for rows.Next() {
		var (
			namaJukir  sql.NullString
			parlok     sql.NullString
			zona       sql.NullString
			plat       sql.NullString
			waktu      sql.NullString
			kendaraan  sql.NullString
			pembayaran sql.NullString
			tarif      sql.NullInt64
		)

		if err := rows.Scan(
			&namaJukir,
			&parlok,
			&zona,
			&plat,
			&waktu,
			&kendaraan,
			&pembayaran,
			&tarif,
		); err != nil {
			return nil, fmt.Errorf("scan monitoring transaksi: %w", err)
		}

		items = append(items, model.TransaksiItemModel{
			NamaJukir:  utils.NullStringValue(namaJukir),
			Parlok:     utils.NullStringValue(parlok),
			Zona:       utils.NullStringValue(zona),
			Plat:       utils.NullStringValue(plat),
			Waktu:      utils.NullStringValue(waktu),
			Kendaraan:  utils.NullStringValue(kendaraan),
			Pembayaran: utils.NullStringValue(pembayaran),
			Tarif:      int(utils.NullInt64Value(tarif)),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate monitoring transaksi: %w", err)
	}

	return &items, nil
}
