package filterpencarian

import (
	"context"
	"database/sql"
	"fmt"

	req "modulegue/internal/domain/mobile/model/filter_pencarian"
	laporanresp "modulegue/internal/domain/mobile/model/laporan"
	riwayatresp "modulegue/internal/domain/mobile/model/riwayat"
	"modulegue/internal/domain/mobile/repository"
)

type FilterPencarianRepoImpl struct {
	db *sql.DB
}

func NewFilterPencarianRepoImpl(db *sql.DB) repository.FilterPencarianRepository {
	return &FilterPencarianRepoImpl{db: db}
}

func (r *FilterPencarianRepoImpl) GetRiwayatParkir(ctx context.Context, req req.FilterPencarianModel) (*riwayatresp.RiwayatParkirModel, error) {
	summaryQuery := `
SELECT 
    COALESCE(SUM(TIMESTAMPDIFF(MINUTE, ps.started_at, COALESCE(ps.finished_at, NOW()))), 0) AS total_duration_minutes,
    COALESCE(SUM(pt.amount), 0) AS total_amount,
    COUNT(CASE WHEN ps.parking_status_id = 3 THEN 1 END) AS completed_count
FROM parking_session ps
LEFT JOIN payment_transaction pt ON pt.reference_id = ps.id AND pt.payment_type = 'PARKING' AND pt.transaction_status = 'SUCCESS'
WHERE ps.customer_user_id = ?;
`
	var summary riwayatresp.ParkingSummaryModel
	err := r.db.QueryRowContext(ctx, summaryQuery, req.UserID).Scan(
		&summary.TotalDurationMinutes,
		&summary.TotalAmount,
		&summary.CompletedCount,
	)
	if err != nil {
		return nil, fmt.Errorf("get riwayat parkir summary: %w", err)
	}

	// Go formatting logic for duration text
	summary.TotalDurationText = fmt.Sprintf("%d Jam %d Mnt", summary.TotalDurationMinutes/60, summary.TotalDurationMinutes%60)

	itemsQuery := `
	SELECT 
		CAST(ps.id AS CHAR) AS id,
		pt.transaction_code AS ticket_no,
		lp.location_name AS location_name,
		'ON_STREET' AS parking_type,
		ps.plate_number AS license_plate,
		mvt.vehicle_type_name AS vehicle_type,
		ps.started_at AS check_in_time,
		ps.finished_at AS check_out_time,
		TIMESTAMPDIFF(MINUTE, ps.started_at, COALESCE(ps.finished_at, NOW())) AS duration_minutes, -- hitung menit untuk diproses ke teks di Go
		ps.amount AS amount,
		CASE 
			WHEN ps.parking_status_id = 3 THEN 'SELESAI'
			WHEN ps.parking_status_id = 5 THEN 'BATAL'
			ELSE 'BERLANGSUNG'
		END AS status
	FROM parking_session ps
	JOIN payment_transaction pt ON pt.reference_id = ps.id AND pt.payment_type = 'PARKING' AND pt.transaction_status = 'SUCCESS'
	JOIN location_parking lp ON ps.location_id = lp.id
	JOIN master_vehicle_type mvt ON ps.vehicle_type_id = mvt.id
	WHERE ps.customer_user_id = ?
	ORDER BY ps.started_at DESC;
`
	rows, err := r.db.QueryContext(ctx, itemsQuery, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("get riwayat parkir items: %w", err)
	}
	defer rows.Close()

	items := make([]riwayatresp.RiwayatParkirItemModel, 0)
	for rows.Next() {
		var item riwayatresp.RiwayatParkirItemModel
		var checkOutTime sql.NullTime
		var durationMinutes int

		if err := rows.Scan(
			&item.ID,
			&item.TicketNo,
			&item.LocationName,
			&item.ParkingType,
			&item.LicensePlate,
			&item.VehicleType,
			&item.CheckInTime,
			&checkOutTime,
			&durationMinutes,
			&item.Amount,
			&item.Status,
		); err != nil {
			return nil, fmt.Errorf("scan riwayat parkir item: %w", err)
		}

		if checkOutTime.Valid {
			item.CheckOutTime = &checkOutTime.Time
		}

		// Go formatting logic for item duration text
		item.DurationText = fmt.Sprintf("%d Jam %d Menit", durationMinutes/60, durationMinutes%60)
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate riwayat parkir items: %w", err)
	}

	return &riwayatresp.RiwayatParkirModel{
		Summary: summary,
		Items:   items,
	}, nil
}

func (r *FilterPencarianRepoImpl) GetRiwayatMembership(ctx context.Context, req req.FilterPencarianModel) (*riwayatresp.RiwayatMembershipModel, error) {
	summaryQuery := `
SELECT 
    COALESCE(mp.package_name, 'Belum Berlangganan') AS package_name,
    IF(mu.id IS NOT NULL AND mu.status = 'ACTIVE' AND mu.expired_at > NOW(), TRUE, FALSE) AS is_active,
    mu.expired_at AS active_until,
    COALESCE(mp.price, 0) AS next_billing_amount,
    FALSE AS is_auto_renew
FROM user_identity ui
LEFT JOIN membership_user mu ON mu.user_id = ui.id AND mu.status = 'ACTIVE' AND mu.expired_at > NOW()
LEFT JOIN membership_package mp ON mp.id = mu.package_id
WHERE ui.id = ?
ORDER BY mu.expired_at DESC
LIMIT 1;
`
	var summary riwayatresp.MembershipSummaryModel
	var activeUntil sql.NullTime
	err := r.db.QueryRowContext(ctx, summaryQuery, req.UserID).Scan(
		&summary.PackageName,
		&summary.IsActive,
		&activeUntil,
		&summary.NextBillingAmount,
		&summary.IsAutoRenew,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			summary = riwayatresp.MembershipSummaryModel{
				PackageName: "Belum Berlangganan",
				IsActive:    false,
			}
		} else {
			return nil, fmt.Errorf("get riwayat membership summary: %w", err)
		}
	}
	if activeUntil.Valid {
		summary.ActiveUntil = &activeUntil.Time
	}

	itemsQuery := `
SELECT 
    CAST(pt.id AS CHAR) AS id,
    pt.transaction_code AS invoice_no,
    CONCAT(mp.package_name, ' - ', mp.duration_days, ' Hari') AS package_name,
    mu.activated_at AS period_start,
    mu.expired_at AS period_end,
    pt.amount AS amount,
    pt.paid_at AS paid_at,
    CASE 
        WHEN pt.transaction_status = 'SUCCESS' THEN 'DIBAYAR'
        WHEN pt.transaction_status = 'CANCELLED' THEN 'BATAL'
        ELSE 'GAGAL'
    END AS status
FROM payment_transaction pt
JOIN membership_user mu ON pt.reference_id = mu.id
JOIN membership_package mp ON mu.package_id = mp.id
WHERE pt.user_id = ? AND pt.payment_type = 'MEMBERSHIP'
ORDER BY pt.created_at DESC;
`
	rows, err := r.db.QueryContext(ctx, itemsQuery, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("get riwayat membership items: %w", err)
	}
	defer rows.Close()

	items := make([]riwayatresp.RiwayatMembershipItemModel, 0)
	for rows.Next() {
		var item riwayatresp.RiwayatMembershipItemModel
		if err := rows.Scan(
			&item.ID,
			&item.InvoiceNo,
			&item.PackageName,
			&item.PeriodStart,
			&item.PeriodEnd,
			&item.Amount,
			&item.PaidAt,
			&item.Status,
		); err != nil {
			return nil, fmt.Errorf("scan riwayat membership item: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate riwayat membership items: %w", err)
	}

	return &riwayatresp.RiwayatMembershipModel{
		Summary: summary,
		Items:   items,
	}, nil
}

func (r *FilterPencarianRepoImpl) GetRiwayatTransaksi(ctx context.Context, req req.FilterPencarianModel) (*riwayatresp.RiwayatTransaksiModel, error) {
	summaryQuery := `
SELECT 
    COALESCE(SUM(CASE WHEN wh.mutation_type = 'CREDIT' THEN wh.amount ELSE 0 END), 0) AS total_income,
    COALESCE(SUM(CASE WHEN wh.mutation_type = 'DEBIT' THEN wh.amount ELSE 0 END), 0) AS total_expense,
    COALESCE(wa.current_balance_amount, 0) AS current_balance
FROM wallet_account wa
LEFT JOIN wallet_history wh ON wh.wallet_id = wa.id
WHERE wa.user_id = ? AND wa.wallet_type_id = 1;
`
	var summary riwayatresp.WalletSummaryModel
	err := r.db.QueryRowContext(ctx, summaryQuery, req.UserID).Scan(
		&summary.TotalIncome,
		&summary.TotalExpense,
		&summary.CurrentBalance,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			summary = riwayatresp.WalletSummaryModel{}
		} else {
			return nil, fmt.Errorf("get riwayat transaksi summary: %w", err)
		}
	}

	itemsQuery := `
SELECT 
    CAST(wh.id AS CHAR) AS id,
    COALESCE(pt.transaction_code, CONCAT('TRX-MUT-', wh.id)) AS reference_no,
    wh.description AS title,
    CASE 
        WHEN wh.reference_type = 'TOPUP' THEN 'TOP_UP'
        WHEN wh.reference_type = 'TRANSFER' THEN 'TRANSFER'
        ELSE 'PAYMENT'
    END AS transaction_type,
    CASE 
        WHEN wh.mutation_type = 'CREDIT' THEN 'IN'
        ELSE 'OUT'
    END AS flow,
    wh.amount AS amount,
    'BERHASIL' AS status, -- Setiap catatan di wallet_history selalu sukses/berhasil
    wh.created_at AS created_at
FROM wallet_history wh
LEFT JOIN payment_transaction pt ON wh.reference_id = pt.reference_id AND pt.payment_type = wh.reference_type
WHERE wh.user_id = ?
ORDER BY wh.created_at DESC;
`
	rows, err := r.db.QueryContext(ctx, itemsQuery, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("get riwayat transaksi items: %w", err)
	}
	defer rows.Close()

	items := make([]riwayatresp.RiwayatTransaksiItemModel, 0)
	for rows.Next() {
		var item riwayatresp.RiwayatTransaksiItemModel
		if err := rows.Scan(
			&item.ID,
			&item.ReferenceNo,
			&item.Title,
			&item.TransactionType,
			&item.Flow,
			&item.Amount,
			&item.Status,
			&item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan riwayat transaksi item: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate riwayat transaksi items: %w", err)
	}

	return &riwayatresp.RiwayatTransaksiModel{
		Summary: summary,
		Items:   items,
	}, nil
}

func (r *FilterPencarianRepoImpl) GetLaporan(ctx context.Context, req req.FilterPencarianModel) (*laporanresp.LaporanModel, error) {
	items, err := r.GetItem(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get items for laporan: %w", err)
	}

	var (
		totalTransaksi       int64
		totalPendapatanJukir int64
	)

	// Golang calculation for total summaries
	for _, item := range *items {
		totalTransaksi += item.TotalTransaksi
		totalPendapatanJukir += item.TotalPendapatanJukir
	}

	result := &laporanresp.LaporanModel{
		TanggalAwal:          req.StartDate,
		TanggalAkhir:         req.EndDate,
		TotalTransaksi:       totalTransaksi,
		TotalPendapatanJukir: totalPendapatanJukir,
		PendapatanPerTanggal: items,
	}

	return result, nil
}

func (r *FilterPencarianRepoImpl) GetItem(ctx context.Context, filter req.FilterPencarianModel) (*[]laporanresp.LaporanItem, error) {
	query := `
SELECT 
    sodr.report_date AS tanggal,
    SUM(sodr.total_transaction) AS total_transaksi,
    SUM(sodr.total_jukir_share) AS total_pendapatan_jukir,
    SUM(sodr.motor_count) AS motor_count,
    SUM(sodr.car_count) AS car_count,
    SUM(sodr.qris_count) AS qris_count,
    SUM(sodr.cash_count) AS cash_count
FROM summary_officer_daily_report sodr
WHERE sodr.officer_user_id = ? 
  AND sodr.report_date >= DATE(?) 
  AND sodr.report_date <= DATE(?)
GROUP BY sodr.report_date
ORDER BY sodr.report_date ASC;
`
	rows, err := r.db.QueryContext(ctx, query, filter.UserID, filter.StartDate, filter.EndDate)
	if err != nil {
		return nil, fmt.Errorf("get laporan items: %w", err)
	}
	defer rows.Close()

	items := make([]laporanresp.LaporanItem, 0)
	for rows.Next() {
		var item laporanresp.LaporanItem
		if err := rows.Scan(
			&item.Tanggal,
			&item.TotalTransaksi,
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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate laporan items: %w", err)
	}

	return &items, nil
}
