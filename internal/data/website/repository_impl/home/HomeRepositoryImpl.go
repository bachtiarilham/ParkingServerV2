package home

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"modulegue/core/utils"
	modelauth "modulegue/internal/domain/shared/model/auth"
	homemodel "modulegue/internal/domain/web/model/home"
)

type HomeRepositoryImpl struct {
	db *sql.DB
}

func NewHomeRepositoryImpl(db *sql.DB) *HomeRepositoryImpl {
	return &HomeRepositoryImpl{db: db}
}

func (r *HomeRepositoryImpl) Login(ctx context.Context, reqModel modelauth.LoginRequestModel) (*modelauth.TokenSetModel, *homemodel.HomeResponseModel, error) {
	identity := strings.TrimSpace(reqModel.Identity)
	if identity == "" {
		return nil, nil, fmt.Errorf("identity is required")
	}

	var (
		userID       int64
		roleID       int64
		passwordHash string
	)

	const query = `
	SELECT
		ui.id AS user_id,
		ui.role_id AS role_id,
		ua.password_hash AS password_hash
	FROM user_identity ui
	JOIN user_auth ua
		ON ua.user_id = ui.id
	JOIN master_role mr
		ON mr.id = ui.role_id
	WHERE (
			ui.username = ?
			OR ui.email = ?
			OR ui.phone_number = ?
		)
	AND ui.status = 'ACTIVE'
	AND mr.is_active = 1
	AND (
			ua.locked_until IS NULL
			OR ua.locked_until <= NOW()
		)
	LIMIT 1;
	`

	if err := r.db.QueryRowContext(ctx, query, identity, identity, identity).Scan(&userID, &roleID, &passwordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, fmt.Errorf("invalid credentials")
		}
		return nil, nil, fmt.Errorf("login web: %w", err)
	}

	homeResp, err := r.GetHome(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	homeResp.UserId = userID
	homeResp.RoleId = roleID
	homeResp.PasswordHash = passwordHash

	return &modelauth.TokenSetModel{}, homeResp, nil
}

func (r *HomeRepositoryImpl) GetHome(ctx context.Context, userID int64) (*homemodel.HomeResponseModel, error) {
	user, err := r.getUserMain(ctx, userID)
	if err != nil {
		return nil, err
	}

	grafikPerDay, err := r.getGrafikPerDay(ctx)
	if err != nil {
		return nil, err
	}

	transaksiAkhir, err := r.getTransaksiAkhir(ctx)
	if err != nil {
		return nil, err
	}

	grafikPerMinggu, err := r.getGrafikPerMinggu(ctx)
	if err != nil {
		return nil, err
	}

	listProvinsi, err := r.getListProvinsi(ctx)
	if err != nil {
		return nil, err
	}

	listKabupaten, err := r.getListKabupaten(ctx)
	if err != nil {
		return nil, err
	}

	listKecamatan, err := r.getListKecamatan(ctx)
	if err != nil {
		return nil, err
	}

	listDesa, err := r.getListDesa(ctx)
	if err != nil {
		return nil, err
	}

	listShift, err := r.getListShift(ctx)
	if err != nil {
		return nil, err
	}

	petugasLapangan, err := r.getPetugasLapangan(ctx)
	if err != nil {
		return nil, err
	}

	statsGlobal, err := r.getStatsGlobal(ctx, userID)
	if err != nil {
		return nil, err
	}

	listZona, err := r.getListZona(ctx)
	if err != nil {
		return nil, err
	}

	listRole, err := r.getListRole(ctx)
	if err != nil {
		return nil, err
	}

	return &homemodel.HomeResponseModel{
		NIK:             user.NIK,
		NoTelp:          user.NoTelp,
		GrafikPerDay:    grafikPerDay,
		TransaksiAkhir:  transaksiAkhir,
		GrafikPerMinggu: grafikPerMinggu,
		ListProvinsi:    listProvinsi,
		ListShift:       listShift,
		Alamat:          user.Alamat,
		ListKecamatan:   listKecamatan,
		PetugasLapangan: petugasLapangan,
		ID:              int(user.ID),
		StatsGlobal:     statsGlobal,
		Token:           "",
		ListKabupaten:   listKabupaten,
		Nama:            user.Nama,
		Foto:            user.Foto,
		ListZona:        listZona,
		ListRole:        listRole,
		ListDesa:        listDesa,
	}, nil
}

type homeUserRow struct {
	ID     int64
	NIK    string
	NoTelp string
	Nama   string
	Foto   string
	Alamat string
}

func (r *HomeRepositoryImpl) getUserMain(ctx context.Context, userID int64) (*homeUserRow, error) {
	const query = `
SELECT
    ui.id AS id,
    ui.nik AS nik,
    ui.phone_number AS no_telp,
    ui.full_name AS nama,
    ui.photo_url AS foto,
    ua.detail_address AS alamat
FROM user_identity ui
LEFT JOIN user_address ua
	
    ON ua.user_id = ui.id
WHERE ui.id = ?
LIMIT 1;
`

	var row homeUserRow
	var (
		nik    sql.NullString
		noTelp sql.NullString
		nama   sql.NullString
		foto   sql.NullString
		alamat sql.NullString
	)
	if err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&row.ID,
		&nik,
		&noTelp,
		&nama,
		&foto,
		&alamat,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user home utama tidak ditemukan")
		}
		return nil, fmt.Errorf("get user home utama: %w", err)
	}

	row.NIK = utils.NullStringValue(nik)
	row.NoTelp = utils.NullStringValue(noTelp)
	row.Nama = utils.NullStringValue(nama)
	row.Foto = utils.NullStringValue(foto)
	row.Alamat = utils.NullStringValue(alamat)

	return &row, nil
}

func (r *HomeRepositoryImpl) getGrafikPerDay(ctx context.Context) (*[]homemodel.GrafikPerDayModel, error) {
	const query = `
WITH RECURSIVE date_range AS (
    SELECT DATE_SUB(CURDATE(), INTERVAL 6 DAY) AS report_date

    UNION ALL

    SELECT DATE_ADD(report_date, INTERVAL 1 DAY)
    FROM date_range
    WHERE report_date < CURDATE()
)
SELECT
    DATE_FORMAT(dr.report_date, '%d/%m/%Y') AS tanggal,

    COALESCE(SUM(
        CASE
            WHEN mvt.vehicle_type_code IN ('MOTOR', 'MOTORCYCLE')
            THEN 1 ELSE 0
        END
    ), 0) AS jml_motor,

    COALESCE(SUM(
        CASE
            WHEN mvt.vehicle_type_code IN ('MOBIL', 'CAR')
            THEN 1 ELSE 0
        END
    ), 0) AS jml_mobil,

    COALESCE(COUNT(fpt.id), 0) AS total

FROM date_range dr

LEFT JOIN financial_parking_transaction fpt
    ON DATE(fpt.paid_at) = dr.report_date
   AND fpt.transaction_status = 'SUCCESS'

LEFT JOIN master_vehicle_type mvt
    ON mvt.id = fpt.vehicle_type_id

GROUP BY dr.report_date

ORDER BY dr.report_date ASC;
`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get grafik per day: %w", err)
	}
	defer rows.Close()

	items := make([]homemodel.GrafikPerDayModel, 0)
	for rows.Next() {
		var item homemodel.GrafikPerDayModel
		var (
			tanggal  sql.NullString
			jmlMotor sql.NullInt64
			jmlMobil sql.NullInt64
			total    sql.NullInt64
		)
		if err := rows.Scan(&tanggal, &jmlMotor, &jmlMobil, &total); err != nil {
			return nil, fmt.Errorf("scan grafik per day: %w", err)
		}
		item.Tanggal = utils.NullStringValue(tanggal)
		item.JmlMotor = int(utils.NullInt64Value(jmlMotor))
		item.JmlMobil = int(utils.NullInt64Value(jmlMobil))
		item.Total = int(utils.NullInt64Value(total))
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate grafik per day: %w", err)
	}
	return &items, nil
}

func (r *HomeRepositoryImpl) getTransaksiAkhir(ctx context.Context) (*[]homemodel.TransaksiAkhirModel, error) {
	const query = `
SELECT
    COALESCE(fpt.plate_number, '') AS plat_nomor,
    COALESCE(ui.full_name, '') AS nama_jukir,
    DATE_FORMAT(fpt.paid_at, '%d/%m/%Y %H:%i:%s') AS waktu,
    COALESCE(mvt.vehicle_type_name, '') AS kendaraan,
    COALESCE(fpt.final_amount, 0) AS tarif

FROM financial_parking_transaction fpt

LEFT JOIN user_identity ui
    ON ui.id = fpt.jukir_user_id

LEFT JOIN master_vehicle_type mvt
    ON mvt.id = fpt.vehicle_type_id

WHERE fpt.transaction_status = 'SUCCESS'

ORDER BY fpt.paid_at DESC

LIMIT 5;
`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get transaksi akhir: %w", err)
	}
	defer rows.Close()

	items := make([]homemodel.TransaksiAkhirModel, 0)
	for rows.Next() {
		var item homemodel.TransaksiAkhirModel
		var (
			platNomor sql.NullString
			namaJukir sql.NullString
			waktu     sql.NullString
			kendaraan sql.NullString
			tarif     sql.NullInt64
		)
		if err := rows.Scan(&platNomor, &namaJukir, &waktu, &kendaraan, &tarif); err != nil {
			return nil, fmt.Errorf("scan transaksi akhir: %w", err)
		}
		item.PlatNomor = utils.NullStringValue(platNomor)
		item.NamaJukir = utils.NullStringValue(namaJukir)
		item.Waktu = utils.NullStringValue(waktu)
		item.Kendaraan = utils.NullStringValue(kendaraan)
		item.Tarif = int(utils.NullInt64Value(tarif))
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate transaksi akhir: %w", err)
	}
	return &items, nil
}

func (r *HomeRepositoryImpl) getGrafikPerMinggu(ctx context.Context) (*homemodel.GrafikPerMingguModel, error) {
	const query = `
SELECT
    COALESCE(SUM(
        CASE
            WHEN mvt.vehicle_type_code IN ('MOTOR', 'MOTORCYCLE')
            THEN 1 ELSE 0
        END
    ), 0) AS motor,

    COALESCE(SUM(
        CASE
            WHEN mvt.vehicle_type_code IN ('MOBIL', 'CAR')
            THEN 1 ELSE 0
        END
    ), 0) AS mobil,

    COUNT(fpt.id) AS total

FROM financial_parking_transaction fpt

JOIN master_vehicle_type mvt
    ON mvt.id = fpt.vehicle_type_id

WHERE fpt.transaction_status = 'SUCCESS'
  AND YEARWEEK(fpt.paid_at, 1) = YEARWEEK(CURDATE(), 1);
`

	var item homemodel.GrafikPerMingguModel
	var (
		motor sql.NullInt64
		mobil sql.NullInt64
		total sql.NullInt64
	)
	if err := r.db.QueryRowContext(ctx, query).Scan(&motor, &mobil, &total); err != nil {
		return nil, fmt.Errorf("get grafik per minggu: %w", err)
	}
	item.Motor = int(utils.NullInt64Value(motor))
	item.Mobil = int(utils.NullInt64Value(mobil))
	item.Total = int(utils.NullInt64Value(total))
	return &item, nil
}

func (r *HomeRepositoryImpl) getListProvinsi(ctx context.Context) (*[]homemodel.ProvinsiModel, error) {
	const query = `
SELECT
    id AS id,
    province_name AS nama_prov
FROM master_province
ORDER BY province_name ASC;
`

	return r.scanProvinceLike(ctx, query, "get list provinsi", "scan list provinsi")
}

func (r *HomeRepositoryImpl) getListKabupaten(ctx context.Context) (*[]homemodel.KabupatenModel, error) {
	const query = `
SELECT
    id AS id,
    regency_name AS nama_kab
FROM master_regency
ORDER BY regency_name ASC;
`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get list kabupaten: %w", err)
	}
	defer rows.Close()

	items := make([]homemodel.KabupatenModel, 0)
	for rows.Next() {
		var item homemodel.KabupatenModel
		var (
			id      sql.NullInt64
			namaKab sql.NullString
		)
		if err := rows.Scan(&id, &namaKab); err != nil {
			return nil, fmt.Errorf("scan list kabupaten: %w", err)
		}
		item.ID = int(utils.NullInt64Value(id))
		item.NamaKab = utils.NullStringValue(namaKab)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate list kabupaten: %w", err)
	}
	return &items, nil
}

func (r *HomeRepositoryImpl) getListKecamatan(ctx context.Context) (*[]homemodel.KecamatanModel, error) {
	const query = `
SELECT
    id AS id,
    district_name AS nama_kec
FROM master_district
ORDER BY district_name ASC;
`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get list kecamatan: %w", err)
	}
	defer rows.Close()

	items := make([]homemodel.KecamatanModel, 0)
	for rows.Next() {
		var item homemodel.KecamatanModel
		var (
			id      sql.NullInt64
			namaKec sql.NullString
		)
		if err := rows.Scan(&id, &namaKec); err != nil {
			return nil, fmt.Errorf("scan list kecamatan: %w", err)
		}
		item.ID = int(utils.NullInt64Value(id))
		item.NamaKec = utils.NullStringValue(namaKec)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate list kecamatan: %w", err)
	}
	return &items, nil
}

func (r *HomeRepositoryImpl) getListDesa(ctx context.Context) (*[]homemodel.DesaModel, error) {
	const query = `
SELECT
    id AS id,
    village_name AS nama_desa
FROM master_village
ORDER BY village_name ASC;
`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get list desa: %w", err)
	}
	defer rows.Close()

	items := make([]homemodel.DesaModel, 0)
	for rows.Next() {
		var item homemodel.DesaModel
		var (
			id       sql.NullInt64
			namaDesa sql.NullString
		)
		if err := rows.Scan(&id, &namaDesa); err != nil {
			return nil, fmt.Errorf("scan list desa: %w", err)
		}
		item.ID = int(utils.NullInt64Value(id))
		item.NamaDesa = utils.NullStringValue(namaDesa)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate list desa: %w", err)
	}
	return &items, nil
}

func (r *HomeRepositoryImpl) getListShift(ctx context.Context) (*[]homemodel.ShiftModel, error) {
	const query = `
SELECT
    ls.id AS id,
    ls.shift_name AS nama_shift,
    TIME_FORMAT(ls.start_time, '%H:%i:%s') AS jam_masuk,
    TIME_FORMAT(ls.end_time, '%H:%i:%s') AS jam_keluar,
    lp.location_name AS nama_parlok

FROM location_shift ls

JOIN location_parking lp
    ON lp.id = ls.location_id

WHERE ls.is_active = 1
  AND lp.is_active = 1

ORDER BY
    lp.location_name ASC,
    ls.start_time ASC;
`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get list shift: %w", err)
	}
	defer rows.Close()

	items := make([]homemodel.ShiftModel, 0)
	for rows.Next() {
		var item homemodel.ShiftModel
		var (
			id         sql.NullInt64
			namaShift  sql.NullString
			jamMasuk   sql.NullString
			jamKeluar  sql.NullString
			namaParlok sql.NullString
		)
		if err := rows.Scan(&id, &namaShift, &jamMasuk, &jamKeluar, &namaParlok); err != nil {
			return nil, fmt.Errorf("scan list shift: %w", err)
		}
		item.ID = int(utils.NullInt64Value(id))
		item.NamaShift = utils.NullStringValue(namaShift)
		item.JamMasuk = utils.NullStringValue(jamMasuk)
		item.JamKeluar = utils.NullStringValue(jamKeluar)
		item.NamaParlok = utils.NullStringValue(namaParlok)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate list shift: %w", err)
	}
	return &items, nil
}

func (r *HomeRepositoryImpl) getPetugasLapangan(ctx context.Context) (*[]homemodel.PetugasLapanganModel, error) {
	const query = `
SELECT
    COALESCE(lp.location_name, '') AS parlok,
    COALESCE(ui.full_name, '') AS nama,
    COALESCE(ui.photo_url, '') AS foto

FROM assignment_officer_current aoc

JOIN user_identity ui
    ON ui.id = aoc.officer_user_id

JOIN master_role mr
    ON mr.id = ui.role_id

LEFT JOIN location_parking lp
    ON lp.id = aoc.location_id

WHERE ui.status = 'ACTIVE'
  AND mr.role_code IN ('OFFICER', 'JUKIR', 'PETUGAS')

ORDER BY
    lp.location_name ASC,
    ui.full_name ASC;
`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get petugas lapangan: %w", err)
	}
	defer rows.Close()

	items := make([]homemodel.PetugasLapanganModel, 0)
	for rows.Next() {
		var item homemodel.PetugasLapanganModel
		var (
			parlok sql.NullString
			nama   sql.NullString
			foto   sql.NullString
		)
		if err := rows.Scan(&parlok, &nama, &foto); err != nil {
			return nil, fmt.Errorf("scan petugas lapangan: %w", err)
		}
		item.Parlok = utils.NullStringValue(parlok)
		item.Nama = utils.NullStringValue(nama)
		item.Foto = utils.NullStringValue(foto)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate petugas lapangan: %w", err)
	}
	return &items, nil
}

func (r *HomeRepositoryImpl) getStatsGlobal(ctx context.Context, userID int64) (*homemodel.StatsGlobalModel, error) {
	const query = `
SELECT
    (
        SELECT COUNT(*)
        FROM financial_parking_transaction fpt
        WHERE fpt.transaction_status = 'SUCCESS'
    ) AS jml_transaksi,

    (
        SELECT COALESCE(SUM(fpt.final_amount), 0)
        FROM financial_parking_transaction fpt
        WHERE fpt.transaction_status = 'SUCCESS'
    ) AS total_pad,

    (
        SELECT COUNT(*)
        FROM user_identity ui
        JOIN master_role mr
            ON mr.id = ui.role_id
        WHERE ui.status = 'ACTIVE'
          AND mr.role_code IN ('OFFICER', 'JUKIR', 'PETUGAS')
    ) AS jml_petugas,

    (
        SELECT COALESCE(wa.current_balance_amount, 0)
        FROM wallet_account wa
        JOIN master_wallet_type mwt
            ON mwt.id = wa.wallet_type_id
        WHERE wa.user_id = ?
          AND mwt.wallet_type_code = 'BALANCE'
          AND wa.is_active = 1
        LIMIT 1
    ) AS saldo;
`

	var item homemodel.StatsGlobalModel
	var saldo sql.NullInt64
	if err := r.db.QueryRowContext(ctx, query, userID).Scan(&item.JmlTransaksi, &item.TotalPAD, &item.JmlPetugas, &saldo); err != nil {
		if err == sql.ErrNoRows {
			return &homemodel.StatsGlobalModel{}, nil
		}
		return nil, fmt.Errorf("get stats global: %w", err)
	}
	item.Saldo = int(utils.NullInt64Value(saldo))
	return &item, nil
}

func (r *HomeRepositoryImpl) getListZona(ctx context.Context) (*[]homemodel.ZonaModel, error) {
	const query = `
SELECT
    id AS id,
    zone_name AS keterangan
FROM location_zone
WHERE is_active = 1
ORDER BY zone_name ASC;
`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get list zona: %w", err)
	}
	defer rows.Close()

	items := make([]homemodel.ZonaModel, 0)
	for rows.Next() {
		var item homemodel.ZonaModel
		var (
			id         sql.NullInt64
			keterangan sql.NullString
		)
		if err := rows.Scan(&id, &keterangan); err != nil {
			return nil, fmt.Errorf("scan list zona: %w", err)
		}
		item.ID = int(utils.NullInt64Value(id))
		item.Keterangan = utils.NullStringValue(keterangan)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate list zona: %w", err)
	}
	return &items, nil
}

func (r *HomeRepositoryImpl) getListRole(ctx context.Context) (*[]homemodel.RoleModel, error) {
	const query = `
SELECT
    id AS id,
    role_name AS nama_role
FROM master_role
WHERE is_active = 1
ORDER BY id ASC;
`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get list role: %w", err)
	}
	defer rows.Close()

	items := make([]homemodel.RoleModel, 0)
	for rows.Next() {
		var item homemodel.RoleModel
		var (
			id       sql.NullInt64
			namaRole sql.NullString
		)
		if err := rows.Scan(&id, &namaRole); err != nil {
			return nil, fmt.Errorf("scan list role: %w", err)
		}
		item.ID = int(utils.NullInt64Value(id))
		item.NamaRole = utils.NullStringValue(namaRole)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate list role: %w", err)
	}
	return &items, nil
}

func (r *HomeRepositoryImpl) scanProvinceLike(ctx context.Context, query, getLabel, scanLabel string) (*[]homemodel.ProvinsiModel, error) {
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", getLabel, err)
	}
	defer rows.Close()

	items := make([]homemodel.ProvinsiModel, 0)
	for rows.Next() {
		var item homemodel.ProvinsiModel
		var (
			id       sql.NullInt64
			namaProv sql.NullString
		)
		if err := rows.Scan(&id, &namaProv); err != nil {
			return nil, fmt.Errorf("%s: %w", scanLabel, err)
		}
		item.ID = int(utils.NullInt64Value(id))
		item.NamaProv = utils.NullStringValue(namaProv)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate list provinsi: %w", err)
	}
	return &items, nil
}
