package helper

import (
	"context"
	"database/sql"
	"fmt"

	model "modulegue/internal/domain/mobile/model/helper"
	"modulegue/internal/domain/mobile/repository"
)

type GetLocationRepositoryImpl struct {
	db *sql.DB
}

func NewGetLocationRepositoryImpl(db *sql.DB) repository.HelperRepository {
	return &GetLocationRepositoryImpl{db: db}
}

func (r *GetLocationRepositoryImpl) GetLokasi(ctx context.Context) (*model.LokasiModel, error) {
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
