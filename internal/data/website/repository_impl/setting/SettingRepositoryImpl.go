package setting

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"modulegue/core/hash"
	"modulegue/core/utils"
	model "modulegue/internal/domain/web/model/setting"
	repository "modulegue/internal/domain/web/repository"
)

type SettingRepositoryImpl struct {
	db *sql.DB
}

func NewSettingRepositoryImpl(db *sql.DB) repository.SettingRepository {
	return &SettingRepositoryImpl{db: db}
}

func (r *SettingRepositoryImpl) AddParlok(ctx context.Context, req model.AddParlokRequestModel) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin add parlok transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	const validateQuery = `
SELECT
    lz.id AS zoneId,
    mp.id AS provinceId,
    mr.id AS regencyId,
    md.id AS districtId,
    mv.id AS villageId

FROM location_zone lz

JOIN master_province mp
    ON mp.id = ?

JOIN master_regency mr
    ON mr.id = ?
   AND mr.province_id = mp.id

JOIN master_district md
    ON md.id = ?
   AND md.regency_id = mr.id

JOIN master_village mv
    ON mv.id = ?
   AND mv.district_id = md.id

WHERE lz.id = ?
  AND lz.is_active = 1

LIMIT 1;
`

	var (
		zoneID     int64
		provinceID int64
		regencyID  int64
		districtID int64
		villageID  int64
	)
	if err := tx.QueryRowContext(ctx, validateQuery, req.IDProv, req.IDKab, req.IDKec, req.IDDes, req.IDZona).Scan(
		&zoneID,
		&provinceID,
		&regencyID,
		&districtID,
		&villageID,
	); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data zona atau wilayah tidak valid")
		}
		return fmt.Errorf("validate add parlok: %w", err)
	}

	const insertParlokQuery = `
INSERT INTO location_parking (
    zone_id,
    location_name,
    address,

    province_id,
    regency_id,
    district_id,
    village_id,

    min_latitude,
    max_latitude,
    min_longitude,
    max_longitude,
    center_latitude,
    center_longitude,

    is_active,
    created_at,
    updated_at
)
VALUES (
    ?,
    ?,
    ?,

    ?,
    ?,
    ?,
    ?,

    ?,
    ?,
    ?,
    ?,
    ?,
    ?,

    1,
    NOW(),
    NOW()
);
`
	result, err := tx.ExecContext(
		ctx,
		insertParlokQuery,
		req.IDZona,
		req.NamaParlok,
		req.JalanParlok,
		req.IDProv,
		req.IDKab,
		req.IDKec,
		req.IDDes,
		req.LatMinArea,
		req.LatMaxArea,
		req.LngMinArea,
		req.LngMaxArea,
		req.CenterAreaX,
		req.CenterAreaY,
	)
	if err != nil {
		return fmt.Errorf("insert location parking: %w", err)
	}

	locationID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get location parking id: %w", err)
	}

	const insertAreaQuery = `
INSERT INTO location_area (
    location_id,
    area_name,
    is_active,
    created_at,
    updated_at
)
VALUES (
    ?,
    'Area Utama',
    1,
    NOW(),
    NOW()
);
`
	if _, err := tx.ExecContext(ctx, insertAreaQuery, locationID); err != nil {
		return fmt.Errorf("insert default area: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit add parlok: %w", err)
	}
	tx = nil
	return nil
}

func (r *SettingRepositoryImpl) RegisterJukir(ctx context.Context, req model.RegisterRequestModel) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin register jukir transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	const checkDuplicateQuery = `
SELECT
    CASE WHEN EXISTS (
        SELECT 1 FROM user_identity
        WHERE nik = ?
        LIMIT 1
    ) THEN 1 ELSE 0 END AS nikUsed,

    CASE WHEN EXISTS (
        SELECT 1 FROM user_identity
        WHERE phone_number = ?
        LIMIT 1
    ) THEN 1 ELSE 0 END AS phoneUsed,

    CASE WHEN EXISTS (
        SELECT 1 FROM user_identity
        WHERE email = ?
        LIMIT 1
    ) THEN 1 ELSE 0 END AS emailUsed,

    CASE WHEN EXISTS (
        SELECT 1 FROM user_identity
        WHERE username = ?
        LIMIT 1
    ) THEN 1 ELSE 0 END AS usernameUsed;
`

	var nikUsed, phoneUsed, emailUsed, usernameUsed int
	if err := tx.QueryRowContext(ctx, checkDuplicateQuery, req.NIK, req.NoTelp, req.Email, req.Username).Scan(
		&nikUsed,
		&phoneUsed,
		&emailUsed,
		&usernameUsed,
	); err != nil {
		return fmt.Errorf("check duplicate user: %w", err)
	}
	if nikUsed == 1 || phoneUsed == 1 || emailUsed == 1 || usernameUsed == 1 {
		return errors.New("data user sudah digunakan")
	}

	passwordHash, err := hash.Hash(req.Password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	const insertUserQuery = `
INSERT INTO user_identity (
    nik,
    full_name,
    phone_number,
    email,
    username,
    role_id,
    photo_url,
    status,
    is_verified,
    created_at,
    updated_at
)
VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    'ACTIVE',
    1,
    NOW(),
    NOW()
);
`
	result, err := tx.ExecContext(
		ctx,
		insertUserQuery,
		req.NIK,
		req.Nama,
		req.NoTelp,
		req.Email,
		req.Username,
		req.IDRole,
		req.Foto,
	)
	if err != nil {
		return fmt.Errorf("insert user identity: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get user id: %w", err)
	}

	const insertAuthQuery = `
INSERT INTO user_auth (
    user_id,
    password_hash,
    failed_login_count,
    locked_until,
    created_at,
    updated_at
)
VALUES (
    ?,
    ?,
    0,
    NULL,
    NOW(),
    NOW()
);
`
	if _, err := tx.ExecContext(ctx, insertAuthQuery, userID, passwordHash); err != nil {
		return fmt.Errorf("insert user auth: %w", err)
	}

	const insertAddressQuery = `
INSERT INTO user_address (
    user_id,
    detail_address,
    is_primary,
    created_at,
    updated_at
)
VALUES (
    ?,
    ?,
    1,
    NOW(),
    NOW()
);
`
	if _, err := tx.ExecContext(ctx, insertAddressQuery, userID, req.Alamat); err != nil {
		return fmt.Errorf("insert user address: %w", err)
	}

	const insertWalletQuery = `
INSERT INTO wallet_account (
    user_id,
    wallet_type_id,
    current_balance_amount,
    is_active,
    created_at,
    updated_at
)
SELECT
    ?,
    mwt.id,
    0,
    1,
    NOW(),
    NOW()
FROM master_wallet_type mwt
WHERE mwt.wallet_type_code = 'BALANCE'
LIMIT 1;
`
	if _, err := tx.ExecContext(ctx, insertWalletQuery, userID); err != nil {
		return fmt.Errorf("insert wallet account: %w", err)
	}

	const insertSummaryQuery = `
INSERT INTO summary_user_home (
    user_id,
    role_id,
    full_name,
    username,
    email,
    phone_number,
    photo_url,
    saldo,
    updated_at
)
SELECT
    ui.id,
    ui.role_id,
    ui.full_name,
    ui.username,
    ui.email,
    ui.phone_number,
    ui.photo_url,
    0,
    NOW()
FROM user_identity ui
WHERE ui.id = ?;
`
	if _, err := tx.ExecContext(ctx, insertSummaryQuery, userID); err != nil {
		return fmt.Errorf("insert summary user home: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit register jukir: %w", err)
	}
	tx = nil
	return nil
}

func (r *SettingRepositoryImpl) SaveSchedule(ctx context.Context, req model.SaveScheduleRequestModel) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin save schedule transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	const validateQuery = `
SELECT
    ui.id AS officerUserId,
    lp.id AS locationId,
    lp.zone_id AS zoneId,
    la.id AS areaId,
    ls.id AS locationShiftId

FROM user_identity ui

JOIN master_role mr
    ON mr.id = ui.role_id

JOIN location_parking lp
    ON lp.id = ?
   AND lp.zone_id = ?
   AND lp.is_active = 1

JOIN location_area la
    ON la.id = ?
   AND la.location_id = lp.id
   AND la.is_active = 1

JOIN location_shift ls
    ON ls.id = ?
   AND ls.location_id = lp.id
   AND ls.is_active = 1

WHERE ui.id = ?
  AND ui.status = 'ACTIVE'
  AND mr.role_code IN ('OFFICER', 'JUKIR', 'PETUGAS')

LIMIT 1;
`

	var (
		officerUserID   int64
		locationID      int64
		zoneID          int64
		areaID          int64
		locationShiftID int64
	)
	if err := tx.QueryRowContext(ctx, validateQuery, req.IDLokasi, req.IDZona, req.IDArea, req.IDShift, req.IDUser).Scan(
		&officerUserID,
		&locationID,
		&zoneID,
		&areaID,
		&locationShiftID,
	); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("data schedule tidak valid")
		}
		return fmt.Errorf("validate save schedule: %w", err)
	}

	if req.ID == 0 {
		const insertQuery = `
INSERT INTO assignment_officer (
    officer_user_id,
    location_id,
    area_id,
    zone_id,
    location_shift_id,
    effective_from,
    effective_to,
    status,
    created_by,
    created_at,
    updated_at
)
VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    DATE(?),
    DATE(?),
    'ACTIVE',
    NULL,
    NOW(),
    NOW()
);
`
		if _, err := tx.ExecContext(ctx, insertQuery, req.IDUser, req.IDLokasi, req.IDArea, req.IDZona, req.IDShift, req.DateAwal, req.DateAkhir); err != nil {
			return fmt.Errorf("insert assignment officer: %w", err)
		}
	} else {
		const updateQuery = `
UPDATE assignment_officer
SET
    officer_user_id = ?,
    location_id = ?,
    area_id = ?,
    zone_id = ?,
    location_shift_id = ?,
    effective_from = DATE(?),
    effective_to = DATE(?),
    status = 'ACTIVE',
    updated_at = NOW()
WHERE id = ?;
`
		if _, err := tx.ExecContext(ctx, updateQuery, req.IDUser, req.IDLokasi, req.IDArea, req.IDZona, req.IDShift, req.DateAwal, req.DateAkhir, req.ID); err != nil {
			return fmt.Errorf("update assignment officer: %w", err)
		}
	}

	const upsertCurrentQuery = `
INSERT INTO assignment_officer_current (
    officer_user_id,
    location_id,
    area_id,
    zone_id,
    location_shift_id,
    effective_from,
    effective_to,
    updated_at
)
VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    DATE(?),
    DATE(?),
    NOW()
)
ON DUPLICATE KEY UPDATE
    location_id = VALUES(location_id),
    area_id = VALUES(area_id),
    zone_id = VALUES(zone_id),
    location_shift_id = VALUES(location_shift_id),
    effective_from = VALUES(effective_from),
    effective_to = VALUES(effective_to),
    updated_at = NOW();
`
	if _, err := tx.ExecContext(ctx, upsertCurrentQuery, req.IDUser, req.IDLokasi, req.IDArea, req.IDZona, req.IDShift, req.DateAwal, req.DateAkhir); err != nil {
		return fmt.Errorf("upsert assignment officer current: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit save schedule: %w", err)
	}
	tx = nil
	return nil
}

func (r *SettingRepositoryImpl) SaveTarif(ctx context.Context, req model.SaveTarifRequestModel) error {
	const query = `
INSERT INTO location_tariff (
    location_id,
    vehicle_type_id,
    base_tariff,
    is_active,
    created_at,
    updated_at
)
SELECT
    ?,
    mvt.id,
    ?,
    1,
    NOW(),
    NOW()
FROM master_vehicle_type mvt
WHERE
    (
        UPPER(?) IN ('MOTOR', 'MOTORCYCLE')
        AND mvt.vehicle_type_code IN ('MOTOR', 'MOTORCYCLE')
    )
    OR
    (
        UPPER(?) IN ('MOBIL', 'CAR')
        AND mvt.vehicle_type_code IN ('MOBIL', 'CAR')
    )
LIMIT 1
ON DUPLICATE KEY UPDATE
    base_tariff = VALUES(base_tariff),
    is_active = 1,
    updated_at = NOW();
`
	result, err := r.db.ExecContext(ctx, query, req.IDLokasi, req.Tarif, req.KeteranganTarif, req.KeteranganTarif)
	if err != nil {
		return fmt.Errorf("save tarif: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("save tarif rows affected: %w", err)
	}
	if affected == 0 {
		return errors.New("keterangan tarif tidak valid")
	}
	return nil
}

func (r *SettingRepositoryImpl) ShowSelectedJukir(ctx context.Context, req model.ShowSelectedJukirRequestModel) (*model.ShowSelectedJukirResponseModel, error) {
	const query = `
SELECT
    COALESCE(ui.nik, '') AS nik,
    COALESCE(ui.username, '') AS username,
    COALESCE(ui.phone_number, '') AS notelp,

    0 AS saldo_min,

    CASE
        WHEN ui.status = 'ACTIVE' THEN 1
        ELSE 0
    END AS idstatus,

    ui.role_id AS idrole,

    COALESCE(ua.detail_address, '') AS alamat,

    ui.id AS id,
    COALESCE(ui.email, '') AS email,

    '' AS password,

    COALESCE(wa.current_balance_amount, 0) AS saldo,

    COALESCE(ui.full_name, '') AS nama,
    COALESCE(ui.photo_url, '') AS foto

FROM user_identity ui

LEFT JOIN user_address ua
    ON ua.user_id = ui.id
   AND ua.is_primary = 1

LEFT JOIN wallet_account wa
    ON wa.user_id = ui.id
   AND wa.is_active = 1

LEFT JOIN master_wallet_type mwt
    ON mwt.id = wa.wallet_type_id
   AND mwt.wallet_type_code = 'BALANCE'

WHERE ui.id = ?

LIMIT 1;
`

	var (
		nik      sql.NullString
		username sql.NullString
		noTelp   sql.NullString
		idStatus sql.NullInt64
		idRole   sql.NullInt64
		alamat   sql.NullString
		id       sql.NullInt64
		email    sql.NullString
		saldo    sql.NullInt64
		nama     sql.NullString
		foto     sql.NullString
	)

	if err := r.db.QueryRowContext(ctx, query, req.ID).Scan(
		&nik,
		&username,
		&noTelp,
		new(sql.NullInt64),
		&idStatus,
		&idRole,
		&alamat,
		&id,
		&email,
		new(string),
		&saldo,
		&nama,
		&foto,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("data jukir tidak ditemukan")
		}
		return nil, fmt.Errorf("show selected jukir: %w", err)
	}

	return &model.ShowSelectedJukirResponseModel{
		NIK:      utils.NullStringValue(nik),
		Username: utils.NullStringValue(username),
		NoTelp:   utils.NullStringValue(noTelp),
		SaldoMin: 0,
		IDStatus: int(utils.NullInt64Value(idStatus)),
		IDRole:   int(utils.NullInt64Value(idRole)),
		Alamat:   utils.NullStringValue(alamat),
		ID:       int(utils.NullInt64Value(id)),
		Email:    utils.NullStringValue(email),
		Password: "",
		Saldo:    int(utils.NullInt64Value(saldo)),
		Nama:     utils.NullStringValue(nama),
		Foto:     utils.NullStringValue(foto),
	}, nil
}

func (r *SettingRepositoryImpl) UpdateProfil(ctx context.Context, req model.UpdateProfilRequestModel) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin update profil transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	statusValue := normalizeStatus(req.IDStatus)
	fotoValue := nullableString(req.Foto)

	const updateUserQuery = `
UPDATE user_identity
SET
    full_name = COALESCE(?, full_name),
    username = COALESCE(?, username),
    phone_number = COALESCE(?, phone_number),
    status = COALESCE(?, status),
    photo_url = COALESCE(?, photo_url),
    updated_at = NOW()
WHERE id = ?;
`
	if _, err := tx.ExecContext(ctx, updateUserQuery, nullableString(req.Nama), nullableString(req.Username), nullableString(req.NoTelp), statusValue, fotoValue, req.IDJukir); err != nil {
		return fmt.Errorf("update user identity: %w", err)
	}

	if req.Alamat != nil {
		const upsertAddressQuery = `
INSERT INTO user_address (
    user_id,
    detail_address,
    is_primary,
    created_at,
    updated_at
)
VALUES (
    ?,
    ?,
    1,
    NOW(),
    NOW()
)
ON DUPLICATE KEY UPDATE
    detail_address = VALUES(detail_address),
    updated_at = NOW();
`
		if _, err := tx.ExecContext(ctx, upsertAddressQuery, req.IDJukir, nullableString(req.Alamat)); err != nil {
			return fmt.Errorf("upsert user address: %w", err)
		}
	}

	if req.Password != nil {
		password := strings.TrimSpace(*req.Password)
		if password != "" {
			passwordHash, err := hash.Hash(password)
			if err != nil {
				return fmt.Errorf("hash password: %w", err)
			}
			const updatePasswordQuery = `
UPDATE user_auth
SET
    password_hash = ?,
    failed_login_count = 0,
    locked_until = NULL,
    updated_at = NOW()
WHERE user_id = ?;
`
			if _, err := tx.ExecContext(ctx, updatePasswordQuery, passwordHash, req.IDJukir); err != nil {
				return fmt.Errorf("update password: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit update profil: %w", err)
	}
	tx = nil
	return nil
}

func nullableString(src *string) any {
	if src == nil {
		return nil
	}
	value := strings.TrimSpace(*src)
	if value == "" {
		return nil
	}
	return value
}

func normalizeStatus(src *string) any {
	if src == nil {
		return nil
	}
	value := strings.TrimSpace(*src)
	if value == "" {
		return nil
	}
	switch strings.ToUpper(value) {
	case "1", "ACTIVE":
		return "ACTIVE"
	case "0", "INACTIVE":
		return "INACTIVE"
	default:
		return value
	}
}
