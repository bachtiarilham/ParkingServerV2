package petugas

import (
	"context"
	"database/sql"
	"fmt"

	"modulegue/core/utils"
	model "modulegue/internal/domain/web/model/petugas"
)

type PetugasRepositoryImpl struct {
	db *sql.DB
}

func NewPetugasRepositoryImpl(db *sql.DB) *PetugasRepositoryImpl {
	return &PetugasRepositoryImpl{db: db}
}

func (r *PetugasRepositoryImpl) GetPetugas(ctx context.Context, reqModel model.PetugasRequestModel) (*model.PetugasResponseModel, error) {
	const query = `
SELECT
    ui.id AS id,
    COALESCE(ui.full_name, '') AS nama,

    COUNT(fpt.id) AS jml_transaksi,

    COALESCE(SUM(fpt.final_amount), 0) AS total_pendapatan,

    CASE
        WHEN ui.status = 'ACTIVE'
         AND (
                aoc.effective_from IS NULL
                OR aoc.effective_from <= NOW()
             )
         AND (
                aoc.effective_to IS NULL
                OR aoc.effective_to >= NOW()
             )
        THEN 1
        ELSE 0
    END AS is_aktif,

    COALESCE(lp.location_name, '') AS parlok

FROM assignment_officer_current aoc

JOIN user_identity ui
    ON ui.id = aoc.officer_user_id

JOIN master_role mr
    ON mr.id = ui.role_id

JOIN location_parking lp
    ON lp.id = aoc.location_id

LEFT JOIN financial_parking_transaction fpt
    ON fpt.jukir_user_id = ui.id
   AND fpt.location_id = lp.id
   AND fpt.transaction_status = 'SUCCESS'

WHERE mr.role_code IN ('OFFICER', 'JUKIR', 'PETUGAS')
  AND (
        ? = 0
        OR lp.id = ?
  )

GROUP BY
    ui.id,
    ui.full_name,
    ui.status,
    aoc.effective_from,
    aoc.effective_to,
    lp.location_name

ORDER BY
    lp.location_name ASC,
    ui.full_name ASC;
`

	rows, err := r.db.QueryContext(ctx, query, reqModel.IDLokasi, reqModel.IDLokasi)
	if err != nil {
		return nil, fmt.Errorf("get petugas: %w", err)
	}
	defer rows.Close()

	items := make([]model.PetugasItemModel, 0)
	for rows.Next() {
		var (
			id              sql.NullInt64
			nama            sql.NullString
			jmlTransaksi    sql.NullInt64
			totalPendapatan sql.NullInt64
			isAktif         sql.NullInt64
			parlok          sql.NullString
		)

		if err := rows.Scan(&id, &nama, &jmlTransaksi, &totalPendapatan, &isAktif, &parlok); err != nil {
			return nil, fmt.Errorf("scan petugas: %w", err)
		}

		items = append(items, model.PetugasItemModel{
			ID:              int(utils.NullInt64Value(id)),
			Nama:            utils.NullStringValue(nama),
			JmlTransaksi:    int(utils.NullInt64Value(jmlTransaksi)),
			TotalPendapatan: int(utils.NullInt64Value(totalPendapatan)),
			IsAktif:         utils.NullInt64Value(isAktif) == 1,
			Parlok:          utils.NullStringValue(parlok),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate petugas: %w", err)
	}

	return &model.PetugasResponseModel{
		Petugas: &items,
	}, nil
}
