package home

import (
	"context"
	"database/sql"
	"fmt"

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

	customerSummary, err := r.getCustomerSummary(ctx, reqModel.UserID)
	if err != nil {
		return nil, err
	}

	jukirSummary, err := r.getJukirSummary(ctx, reqModel.UserID)
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
		Profile:         profile,
		CustomerSummary: customerSummary,
		JukirSummary:    jukirSummary,
		Events:          events,
		News:            news,
		Warnings:        warnings,
	}, nil
}

func (r *HomeRepositoryImpl) getProfile(ctx context.Context, userID int64) (*model.ProfileModel, error) {
	query := `SELECT id, full_name FROM system_user WHERE id = ? LIMIT 1`

	var profile model.ProfileModel
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&profile.ID, &profile.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("get profile: %w", err)
	}

	return &profile, nil
}

func (r *HomeRepositoryImpl) getCustomerSummary(ctx context.Context, userID int64) (*model.CustomerSummaryModel, error) {
	query := `
		SELECT COALESCE(current_balance_amount, 0)
		FROM user_wallet
		WHERE user_id = ? AND wallet_status = 'active'
		LIMIT 1
	`

	var balance int64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&balance)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("get customer summary: %w", err)
	}
	if err == sql.ErrNoRows {
		balance = 0
	}

	return &model.CustomerSummaryModel{
		Saldo:       balance,
		ExpiredDate: "",
	}, nil
}

func (r *HomeRepositoryImpl) getJukirSummary(ctx context.Context, userID int64) (*model.JukirSummaryModel, error) {
	query := `
		SELECT
			COALESCE(SUM(fpt.jukir_share), 0) AS pendapatan,
			COALESCE(MAX(pl.location_name), '') AS lokasi,
			COALESCE(MAX(pa.area_name), '') AS area,
			COALESCE(MAX(pz.zone_name), '') AS zona
		FROM financial_parking_transaction fpt
		LEFT JOIN parking_location pl ON pl.id = fpt.location_id
		LEFT JOIN parking_area pa ON pa.id = pl.area_id
		LEFT JOIN parking_zone pz ON pz.id = pl.zone_id
		WHERE fpt.jukir_user_id = ?
	`

	var summary model.JukirSummaryModel
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&summary.Pendapatan,
		&summary.Lokasi,
		&summary.Area,
		&summary.Zona,
	)
	if err != nil {
		return nil, fmt.Errorf("get jukir summary: %w", err)
	}

	return &summary, nil
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
