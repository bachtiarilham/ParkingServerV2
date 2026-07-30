package parking

import (
	"context"
	"database/sql"
	"fmt"

	model "modulegue/internal/domain/mobile/model/parking"
	"modulegue/internal/domain/mobile/repository"

	"strings"

	uuid "github.com/google/uuid"
)

type ParkingRepositoryImpl struct {
	db *sql.DB
}

func NewParkingRepositoryImpl(db *sql.DB) repository.ParkingRepository {
	return &ParkingRepositoryImpl{db: db}
}

func (r *ParkingRepositoryImpl) GetParkingMetadata(ctx context.Context, req model.PostParkingRequestModel) (*model.ParkingBusinessModel, error) {
	query :=
		`
		SELECT 
		aoc.location_id AS LocationId,
		aoc.zone_id AS ZoneId,
		lt.tariff_amount AS Amount
	FROM assignment_officer_current aoc
	JOIN location_tariff lt 
		ON aoc.location_id = lt.location_id
	JOIN master_vehicle_type mvt 
		ON lt.vehicle_type_id = mvt.id
	WHERE aoc.officer_user_id = ? 
	AND mvt.vehicle_type_code = ?
	AND lt.is_active = 1;
	`

	var (
		businessModel model.ParkingBusinessModel
	)

	err := r.db.QueryRowContext(
		ctx,
		query,
		req.OfficerUserId,
		req.VehicleTypeCode,
	).Scan(
		&businessModel.LocationId,
		&businessModel.ZoneId,
		&businessModel.Amount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("kebutuhan parking tidak ditemukan")
		}
		return nil, fmt.Errorf("post parking: %w", err)
	}
	return &businessModel, nil
}

func (r *ParkingRepositoryImpl) CreateParkingTransaction(ctx context.Context, req1 model.PostParkingRequestModel, req model.ParkingBusinessModel) (*model.PostParkingResponseModel, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin parking transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 1. Insert initial parking session (with temporary session code)
	insertSessionQuery := `
		INSERT INTO parking_session (
		session_code,
		customer_user_id,
		officer_user_id,
		vehicle_id,
		vehicle_type_id,
		plate_number,
		zone_id,
		location_id,
		area_id,
		parking_status_id,
		amount,
		started_at,
		finished_at,
		cancelled_at,
		created_at,
		updated_at
	)
	VALUES (
		?,         -- 1. session_code (Wajib diisi kode unik, misal: 'SESSION-20260727-001')
		NULL,         -- 2. customer_user_id (Bisa diisi ID user atau nil/NULL)
		?,         -- 3. officer_user_id (ID Petugas/Jukir)
		NULL,         -- 4. vehicle_id (Bisa diisi ID kendaraan terdaftar atau nil/NULL)
		?,         -- 5. vehicle_type_id (1 = Motor, 2 = Mobil)
		UPPER(?),  -- 6. plate_number (Nomor plat otomatis kapital)
		?,         -- 7. zone_id
		?,         -- 8. location_id
		?,         -- 9. area_id
		1,         -- 10. parking_status_id (Status awal parkir)
		?,         -- 11. amount (Nominal tarif dasar)
		NOW(),     -- 12. started_at (Jam mulai parkir otomatis waktu saat ini)
		NULL,      -- 13. finished_at (NULL karena sesi baru mulai)
		NULL,      -- 14. cancelled_at (NULL)
		NOW(),     -- 15. created_at
		NOW()      -- 16. updated_at
	);
	`
	sessionCode := fmt.Sprintf("SESS-PRK-%s", uuid.NewString())
	var vehicle_id int64
	switch strings.ToLower(req1.VehicleTypeCode) {
	case "motor":
		vehicle_id = 1
	case "car":
		vehicle_id = 2
	}

	res, err := tx.ExecContext(
		ctx,
		insertSessionQuery,
		sessionCode,         // 1. -> session_code
		req1.OfficerUserId,  // 2. -> officer_user_id (req.TransactionCode DIBUANG)
		vehicle_id,          // 4. -> vehicle_type_id
		req1.PlateNumber,    // 5. -> plate_number (akan di-UPPER di SQL)
		req.ZoneId,          // 6. -> zone_id
		req.LocationId,      // 7. -> location_id
		req1.SelectedAreaId, // 8. -> area_id
		req.Amount,          // 10. -> amount
	)
	if err != nil {
		return nil, fmt.Errorf("insert parking session: %w", err)
	}

	sessionId, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("get session id: %w", err)
	}

	// 2. Retrieve final response data
	selectQuery := `
	SELECT
		ps.session_code AS sessionCode,
		ps.plate_number AS plateNumber,
		ps.started_at AS waktu
	FROM parking_session ps
	WHERE ps.id = ? AND ps.officer_user_id = ?
	LIMIT 1;
	`

	var (
		resp model.PostParkingResponseModel
	)
	if err = tx.QueryRowContext(ctx, selectQuery, sessionId, req1.OfficerUserId).Scan(
		&resp.SessionCode,
		&resp.PlateNumber,
		&resp.Waktu,
	); err != nil {
		return nil, fmt.Errorf("retrieve final parking details: %w", err)
	}

	// 5. Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit parking transaction: %w", err)
	}

	return &resp, nil
}
