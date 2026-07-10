package laporan

import (
	"context"
	"database/sql"
	"fmt"

	model "modulegue/internal/domain/mobile/model/laporan"
	"modulegue/internal/domain/mobile/repository"
)

type LaporanRepositoryImpl struct {
	db *sql.DB
}

func NewLaporanRepositoryImpl(db *sql.DB) repository.LaporanRepository {
	return &LaporanRepositoryImpl{db: db}
}

func (r *LaporanRepositoryImpl) GetLaporan(ctx context.Context, filter model.LaporanRequestModel) (*model.LaporanModel, error) {
	query := `
		SELECT
			DATE(?) AS tanggalAwal,
			DATE(?) AS tanggalAkhir,

			COALESCE(SUM(total_transaction), 0) AS totalTransaksi,
			COALESCE(SUM(total_income), 0) AS totalPendapatan,
			COALESCE(SUM(total_jukir_share), 0) AS totalPendapatanJukir,

			COALESCE(SUM(motor_count), 0) AS totalMotor,
			COALESCE(SUM(car_count), 0) AS totalMobil,
			COALESCE(SUM(qris_count), 0) AS totalQris,
			COALESCE(SUM(cash_count), 0) AS totalCash

		FROM summary_officer_daily_report
		WHERE officer_user_id = ?
		AND report_date >= DATE(?)
		AND report_date <= DATE(?);
	`

	var result model.LaporanModel
	if err := r.db.QueryRowContext(
		ctx,
		query,
		filter.StartDate, filter.EndDate, filter.UserID,
		filter.StartDate, filter.EndDate,
	).Scan(
		&result.TanggalAwal,
		&result.TanggalAkhir,
		&result.TotalTransaksi,
		&result.TotalPendapatanJukir,
		new(int64),
		new(int64),
		new(int64),
		new(int64),
		new(int64),
	); err != nil {
		return nil, fmt.Errorf("get laporan: %w", err)
	}

	items, err := r.GetItem(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get item: %w", err)
	}
	result.PendapatanPerTanggal = items

	return &result, nil
}

func (r *LaporanRepositoryImpl) GetItem(ctx context.Context, filter model.LaporanRequestModel) (*[]model.LaporanItem, error) {
	query := `
		WITH RECURSIVE date_range AS (
			SELECT DATE(?) AS report_date

			UNION ALL

			SELECT DATE_ADD(report_date, INTERVAL 1 DAY)
			FROM date_range
			WHERE report_date < DATE(?)
		)

		SELECT
			dr.report_date AS tanggal,

			COALESCE(sodr.total_transaction, 0) AS totalTransaksi,
			COALESCE(sodr.total_income, 0) AS totalPendapatan,
			COALESCE(sodr.total_jukir_share, 0) AS totalPendapatanJukir,

			COALESCE(sodr.motor_count, 0) AS motorCount,
			COALESCE(sodr.car_count, 0) AS carCount,
			COALESCE(sodr.qris_count, 0) AS qrisCount,
			COALESCE(sodr.cash_count, 0) AS cashCount

		FROM date_range dr

		LEFT JOIN summary_officer_daily_report sodr
			ON sodr.report_date = dr.report_date
		AND sodr.officer_user_id = ?

		ORDER BY dr.report_date ASC;
	`

	rows, err := r.db.QueryContext(ctx, query, filter.StartDate, filter.EndDate, filter.UserID)
	if err != nil {
		return nil, fmt.Errorf("get laporan item: %w", err)
	}
	defer rows.Close()

	items := make([]model.LaporanItem, 0)
	for rows.Next() {
		var item model.LaporanItem
		if err := rows.Scan(
			&item.Tanggal,
			&item.TotalTransaksi,
			new(int64),
			&item.TotalPendapatanJukir,
			&item.MotorCount,
			&item.CarCount,
			&item.QrisCount,
			&item.CashCount,
		); err != nil {
			return nil, fmt.Errorf("scan laporan item: %w", err)
		}
		items = append(items, item)
	}

	return &items, nil
}
