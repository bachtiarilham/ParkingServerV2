package sync

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"modulegue/internal/domain/mobile/repository"
)

type SyncRepoImpl struct {
	db *sql.DB
}

func NewSyncRepoImpl(db *sql.DB) repository.SyncRepo {
	return &SyncRepoImpl{db: db}
}

func (r *SyncRepoImpl) SyncAfterMembership(ctx context.Context, userId int64, packageId int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin membership callback tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	buatschquery := `
		SELECT 
			mp.package_name AS membership_package_name,
			mu.expired_at AS membership_expired_at
		FROM membership_user mu
		JOIN membership_package mp 
			ON mu.package_id = mp.id
		WHERE mu.user_id = ? 
		  AND mu.package_id = ?
		ORDER BY mu.expired_at DESC
		LIMIT 1
	`
	var packageName string
	var expiredAt time.Time
	err = tx.QueryRowContext(ctx, buatschquery, userId, packageId).Scan(&packageName, &expiredAt)
	if err != nil {
		return err
	}

	updateSummary := `
	UPDATE summary_customer_home 
	SET membership_package_name = ?, membership_expired_at = ?
	WHERE user_id = ?;
	`
	if _, err = tx.ExecContext(ctx, updateSummary, packageName, expiredAt, userId); err != nil {
		return fmt.Errorf("update payment transaction: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit membership callback tx: %w", err)
	}

	return nil
}

func (r *SyncRepoImpl) SyncAfterParkir(ctx context.Context, userId int64, refId int64, amount int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin membership callback tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 1. Deklarasikan variabel untuk menampung 9 kolom hasil SELECT
	var (
		locationID                     int64
		areaID                         int64
		zoneID                         int64
		vehicleTypeCode                string
		partnerShare                   int64
		companyShare                   int64
		govShare                       int64
		paymentMethodCode              string
		isMotor, isCar, isQRIS, isCash int
	)

	ambilLokasiQuery := `
	SELECT 
		ps.location_id,
		ps.area_id,
		ps.zone_id,
		mvt.vehicle_type_code,
		pt.amount,
		pt.partner_share,
		pt.company_share,
		pt.gov_share,
		mpm.payment_method_code
	FROM parking_session ps
	JOIN location_parking lp 
		ON ps.location_id = lp.id
	JOIN master_vehicle_type mvt 
		ON ps.vehicle_type_id = mvt.id
	JOIN payment_transaction pt 
		ON pt.reference_id = ps.id 
		AND pt.payment_type = 'PARKING'
	JOIN master_payment_method mpm 
		ON pt.payment_method_id = mpm.id
	WHERE ps.id = ? 
	AND ps.parking_status_id = 3 
	AND pt.transaction_status = 'SUCCESS' 
	LIMIT 1;
	`

	// 2. Scan ke 9 variabel sesuai urutan SELECT (dari atas ke bawah)
	err = tx.QueryRowContext(ctx, ambilLokasiQuery, refId).Scan(
		&locationID,
		&areaID,
		&zoneID,
		&vehicleTypeCode,
		&amount,
		&partnerShare,
		&companyShare,
		&govShare,
		&paymentMethodCode,
	)
	if err != nil {
		return fmt.Errorf("ambil detail parking session: %w", err)
	}

	if vehicleTypeCode == "MOTOR" {
		isMotor = 1
		isCar = 0
	} else {
		isMotor = 0
		isCar = 1
	}
	switch paymentMethodCode {
	case "QRIS":
		isQRIS = 1
		isCash = 0
	case "CASH":
		isQRIS = 0
		isCash = 1
	default:
		// Jika pembayaran wallet/manual
		isQRIS = 0
		isCash = 0
	}

	summarylocationquery :=
		`
		INSERT INTO summary_location_daily_report (
			report_date, 
			location_id, 
			area_id, 
			zone_id, 
			total_transaction, 
			total_income,
			total_jukir_share, 
			total_company_share, 
			motor_count, 
			car_count
		)
		VALUES (
			CURRENT_DATE(), 
			?, -- location_id
			?, -- area_id
			?, -- zone_id
			1, -- total_transaction awal
			?, -- total_income (nominal pembayaran)
			?, -- total_jukir_share (bagian jukir)
			?, -- total_company_share (bagian platform/pemda)
			?, -- motor_count (isi 1 jika kendaraan = motor, 0 jika mobil)
			?  -- car_count (isi 1 jika kendaraan = mobil, 0 jika motor)
		)
		ON DUPLICATE KEY UPDATE 
			total_transaction   = total_transaction + 1,
			total_income        = total_income + VALUES(total_income),
			total_jukir_share   = total_jukir_share + VALUES(total_jukir_share),
			total_company_share = total_company_share + VALUES(total_company_share),
			motor_count         = motor_count + VALUES(motor_count),
			car_count           = car_count + VALUES(car_count);
		`

	if _, err = tx.ExecContext(ctx, summarylocationquery, locationID, areaID, zoneID, amount, partnerShare, companyShare, isMotor, isCar); err != nil {
		return fmt.Errorf("update summary location: %w", err)
	}

	summaryofficequery :=
		`
	INSERT INTO summary_officer_daily_report (
		report_date, 
		officer_user_id, 
		location_id, 
		area_id, 
		zone_id, 
		total_transaction, 
		total_jukir_share, 
		motor_count, 
		car_count, 
		qris_count, 
		cash_count
	)
	VALUES (
		CURRENT_DATE(), 
		?, -- officer_user_id (ID Jukir)
		?, -- location_id
		?, -- area_id
		?, -- zone_id
		1, -- total_transaction awal
		?, -- total_jukir_share (bagian jukir)
		?, -- motor_count (isi 1 jika motor, 0 jika mobil)
		?, -- car_count (isi 1 jika mobil, 0 jika motor)
		?, -- qris_count (isi 1 jika bayar pakai QRIS, 0 jika cash/wallet)
		?  -- cash_count (isi 1 jika bayar pakai Cash, 0 jika QRIS/wallet)
	)
	ON DUPLICATE KEY UPDATE 
		total_transaction   = total_transaction + 1,
		total_jukir_share   = total_jukir_share + VALUES(total_jukir_share),
		motor_count         = motor_count + VALUES(motor_count),
		car_count           = car_count + VALUES(car_count),
		qris_count          = qris_count + VALUES(qris_count),
		cash_count          = cash_count + VALUES(cash_count);
	`
	if _, err = tx.ExecContext(ctx, summaryofficequery, userId, locationID, areaID, zoneID, partnerShare, isMotor, isCar, isQRIS, isCash); err != nil {
		return fmt.Errorf("update summary location: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit membership callback tx: %w", err)
	}
	return nil
}
