package home

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	model "modulegue/internal/domain/mobile/model/home"
	"modulegue/internal/domain/mobile/repository"
)

type HomeRepositoryImpl struct {
	db *sql.DB
}

func NewHomeRepositoryImpl(db *sql.DB) repository.HomeRepository {
	return &HomeRepositoryImpl{db: db}
}

func (r *HomeRepositoryImpl) GetHome(ctx context.Context, reqModel model.GetHomeReqModel) (*model.HomeModel, error) {
	profile, err := r.getProfile(ctx, reqModel.UserID)
	if err != nil {
		return nil, err
	}

	events, news, err := r.getRecentEventsAndNews(ctx, 10, 0)
	if err != nil {
		return nil, err
	}

	warnings, err := r.getWarnings(ctx, reqModel.UserID)
	if err != nil {
		return nil, err
	}

	return &model.HomeModel{
		Profile:  profile,
		Events:   events,
		News:     news,
		Warnings: warnings,
	}, nil
}

func (r *HomeRepositoryImpl) getProfile(ctx context.Context, userID int64) (*model.ProfileModel, error) {
	query := `
		SELECT
			su.id,
			su.full_name AS name,
			COALESCE(uw.current_balance_amount, 0) AS saldo,
			oa.effective_to AS expired_date,
			COALESCE(SUM(fpt.jukir_share), 0) AS pendapatan,
			pl.location_name AS lokasi,
			pa.area_name AS area,
			pz.zone_name AS zona
		FROM system_user su
		LEFT JOIN user_wallet uw
			ON uw.user_id = su.id
			AND uw.wallet_type = 'emoney'
		LEFT JOIN officer_assignment_current oa
			ON oa.officer_user_id = su.id
		LEFT JOIN parking_location pl
			ON pl.id = oa.location_id
		LEFT JOIN parking_area pa
			ON pa.id = oa.area_id
		LEFT JOIN parking_zone pz
			ON pz.id = oa.zone_id
		LEFT JOIN financial_parking_transaction fpt
			ON fpt.jukir_user_id = su.id
		WHERE su.id = ? 
			GROUP BY
			su.id, su.full_name, uw.current_balance_amount,
			oa.effective_to, pl.location_name, pa.area_name, pz.zone_name;
	`

	var (
		profile     model.ProfileModel
		name        sql.NullString
		saldo       sql.NullFloat64
		expiredDate sql.NullTime
		pendapatan  sql.NullFloat64
		lokasi      sql.NullString
		area        sql.NullString
		zona        sql.NullString
	)

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&profile.ID,
		&name,
		&saldo,
		&expiredDate,
		&pendapatan,
		&lokasi,
		&area,
		&zona,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("get profile: %w", err)
	}

	profile.Name = name.String
	profile.Saldo = int64(saldo.Float64)
	profile.Pendapatan = int64(pendapatan.Float64)
	profile.Lokasi = lokasi.String
	profile.Area = area.String
	profile.Zona = zona.String

	if expiredDate.Valid {
		profile.ExpiredDate = expiredDate.Time.Format(time.RFC3339)
	}

	return &profile, nil
}

func (r *HomeRepositoryImpl) getRecentEventsAndNews(ctx context.Context, limit, offset int) ([]model.EventsModel, []model.NewsModel, error) {
	query := `
		SELECT id, title, description, publish_date, image_url, content_type
		FROM customer_news_and_events
		WHERE is_active = 1
		ORDER BY publish_date DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("get events and news: %w", err)
	}
	defer rows.Close()

	events := []model.EventsModel{}
	news := []model.NewsModel{}

	for rows.Next() {
		var (
			id          int64
			title       string
			description string
			publishDate sql.NullTime
			imageURL    sql.NullString
			contentType string
		)

		if err := rows.Scan(&id, &title, &description, &publishDate, &imageURL, &contentType); err != nil {
			return nil, nil, fmt.Errorf("scan events/news: %w", err)
		}

		if contentType == "news" {
			news = append(news, model.NewsModel{
				ID:          id,
				Title:       title,
				Description: description,
				Date:        publishDate.Time,
				ImageURL:    imageURL.String,
				ContentType: contentType,
			})
			continue
		}

		events = append(events, model.EventsModel{
			ID:          id,
			Title:       title,
			Description: description,
			Date:        publishDate.Time,
			ImageURL:    imageURL.String,
			ContentType: contentType,
		})
	}

	return events, news, nil
}

func (r *HomeRepositoryImpl) getWarnings(ctx context.Context, userID int64) (*model.WarningsModel, error) {
	warnings := &model.WarningsModel{}

	var nik sql.NullString
	if err := r.db.QueryRowContext(ctx, `SELECT nik FROM system_user WHERE id = ? LIMIT 1`, userID).Scan(&nik); err == nil {
		if nik.String == "" {
			warnings.Profile = "Profil belum lengkap"
		}
	}

	var openAlertCount int64
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM admin_alert_event
		WHERE alert_status = 'open' AND officer_user_id = ?
	`, userID).Scan(&openAlertCount); err == nil && openAlertCount > 0 {
		warnings.Parking = fmt.Sprintf("Ada %d alert yang belum ditangani", openAlertCount)
	}

	var unpaidTxCount int64
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM financial_parking_transaction
		WHERE customer_user_id = ? AND transaction_status = 'unpaid'
	`, userID).Scan(&unpaidTxCount); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("get unpaid transactions warning: %w", err)
	}

	var openDisputeCount int64
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM financial_dispute_case
		WHERE opened_by_user_id = ? AND case_status = 'open'
	`, userID).Scan(&openDisputeCount); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("get dispute warning: %w", err)
	}

	if unpaidTxCount > 0 || openDisputeCount > 0 {
		warnings.Finance = fmt.Sprintf("Ada %d transaksi belum lunas dan %d sengketa aktif", unpaidTxCount, openDisputeCount)
	}

	return warnings, nil
}
