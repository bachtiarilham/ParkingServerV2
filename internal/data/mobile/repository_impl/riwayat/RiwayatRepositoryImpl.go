package riwayat

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	model "modulegue/internal/domain/mobile/model/riwayat"
	"modulegue/internal/domain/mobile/repository"
)

type RiwayatRepositoryImpl struct {
	db *sql.DB
}

func NewRiwayatRepositoryImpl(db *sql.DB) repository.RiwayatRepository {
	return &RiwayatRepositoryImpl{db: db}
}

func (r *RiwayatRepositoryImpl) GetRiwayat(ctx context.Context, req model.RiwayatRequestModel) (*model.RiwayatModel, error) {
	startDate := strings.TrimSpace(req.StartDate)
	endDate := strings.TrimSpace(req.EndDate)
	lokasi := strings.TrimSpace(req.Lokasi)
	statusFilter := ""
	paymentFilter := strings.TrimSpace(req.Payment)
	vehicleFilter := strings.TrimSpace(req.Vehicle)

	if startDate == "" {
		startDate = time.Now().Format("2006-01-02")
	}
	if endDate == "" {
		endDate = startDate
	}

	query := `
		SELECT
			DATE_FORMAT(COALESCE(fpt.paid_at, fpt.occurred_at), '%Y-%m-%d') AS section_date,
			fpt.transaction_code,
			COALESCE(fpt.plate_number, ps.plate_number, '') AS plate_number,
			COALESCE(vt.vehicle_type_name, '') AS vehicle_type,
			DATE_FORMAT(COALESCE(fpt.paid_at, fpt.occurred_at), '%Y-%m-%d %H:%i:%s') AS item_time,
			fpt.final_amount,
			CASE
				WHEN UPPER(fpt.operation_type) IN ('ENTRY', 'MASUK', 'CHECKIN') THEN 1
				ELSE 0
			END AS is_entry
		FROM financial_parking_transaction fpt
		LEFT JOIN parking_session ps ON ps.id = fpt.session_id
		LEFT JOIN vehicle_type vt ON vt.id = fpt.vehicle_type_id
		LEFT JOIN parking_location pl ON pl.id = fpt.location_id
		WHERE DATE(fpt.occurred_at) BETWEEN ? AND ?
		  AND (? = '' OR pl.location_name = ? OR pl.location_code = ?)
		  AND (? = '' OR UPPER(COALESCE(fpt.transaction_status, '')) = UPPER(?))
		  AND (? = '' OR UPPER(COALESCE(CAST(fpt.payment_method AS CHAR), '')) = UPPER(?))
		  AND (? = '' OR UPPER(vt.vehicle_type_name) = UPPER(?) OR UPPER(vt.vehicle_type_code) = UPPER(?))
		  AND (? = 0 OR fpt.customer_user_id = ? OR fpt.jukir_user_id = ? OR fpt.officer_user_id = ?)
		  AND UPPER(COALESCE(fpt.transaction_status, '')) <> 'VOID'
		ORDER BY COALESCE(fpt.paid_at, fpt.occurred_at) DESC
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		startDate, endDate,
		lokasi, lokasi, lokasi,
		statusFilter, statusFilter,
		paymentFilter, paymentFilter,
		vehicleFilter, vehicleFilter, vehicleFilter,
		req.UserID, req.UserID, req.UserID, req.UserID,
	)
	if err != nil {
		return nil, fmt.Errorf("get riwayat: %w", err)
	}
	defer rows.Close()

	sections := []model.RiwayatSectionModel{}
	sectionIndexByDate := map[string]int{}

	for rows.Next() {
		var (
			sectionDate string
			code        string
			plateNumber string
			vehicleType string
			itemTime    string
			amount      int64
			isEntryInt  int
		)

		if err := rows.Scan(&sectionDate, &code, &plateNumber, &vehicleType, &itemTime, &amount, &isEntryInt); err != nil {
			return nil, fmt.Errorf("scan riwayat item: %w", err)
		}

		idx, exists := sectionIndexByDate[sectionDate]
		if !exists {
			dateCopy := sectionDate
			sections = append(sections, model.RiwayatSectionModel{
				Date:  &dateCopy,
				Items: []model.RiwayatItemModel{},
			})
			idx = len(sections) - 1
			sectionIndexByDate[sectionDate] = idx
		}

		codeCopy := code
		plateCopy := plateNumber
		vehicleCopy := vehicleType
		timeCopy := itemTime
		amountCopy := amount
		isEntry := isEntryInt == 1

		sections[idx].Items = append(sections[idx].Items, model.RiwayatItemModel{
			Code:        &codeCopy,
			PlateNumber: &plateCopy,
			VehicleType: &vehicleCopy,
			Time:        &timeCopy,
			Amount:      &amountCopy,
			IsEntry:     &isEntry,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate riwayat: %w", err)
	}

	return &model.RiwayatModel{
		Sections: sections,
	}, nil
}
