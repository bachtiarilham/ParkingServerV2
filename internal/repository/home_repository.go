package repository

import (
	"context"
	"database/sql"
	"fmt"
	"modulegue/internal/domain/home"
)

type HomeRepository struct {
	db *sql.DB
}

func NewHomeRepository(db *sql.DB) home.Repository {
	return &HomeRepository{db: db}
}

func (r *HomeRepository) GetProfile(ctx context.Context, userID int64) (*home.Profile, error) {
	query := `SELECT id, full_name FROM system_user WHERE id = ? LIMIT 1`
	var p home.Profile
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&p.ID, &p.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("get profile: %w", err)
	}
	return &p, nil
}

func (r *HomeRepository) GetSummary(ctx context.Context, userID int64) (*home.Summary, error) {
	query := `SELECT current_balance_amount FROM user_wallet WHERE user_id = ? AND wallet_status = 'active' LIMIT 1`
	var balance int64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&balance)
	if err != nil && err != sql.ErrNoRows {
		// Jika wallet tidak ditemukan (ErrNoRows), asumsikan saldo 0
		return nil, fmt.Errorf("get wallet balance: %w", err)
	}
	if err == sql.ErrNoRows {
		balance = 0 // Atau return error jika wallet wajib ada
	}

	// Asumsi tidak ada expired_date untuk saat ini
	summary := &home.Summary{
		Saldo: balance,
		// ExpiredDate: nil, // Tidak digunakan
	}
	return summary, nil
}

func (r *HomeRepository) GetRecentEventsAndNews(ctx context.Context, limit int, offset int) ([]home.EventOrNews, error) {
	query := `
		SELECT id, title, description, publish_date, image_url, content_type
		FROM customer_news_and_events
		WHERE is_active = 1
		ORDER BY publish_date DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get events/news: %w", err)
	}
	defer rows.Close()

	var events []home.EventOrNews
	for rows.Next() {
		var e home.EventOrNews
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.Date, &e.ImageURL, &e.ContentType)
		if err != nil {
			return nil, fmt.Errorf("scan event/news: %w", err)
		}
		events = append(events, e)
	}
	return events, nil
}

func (r *HomeRepository) GetWarnings(ctx context.Context, userID int64) (*home.Warnings, error) {
	warnings := &home.Warnings{}

	// Cek profil (contoh: cek apakah NIK kosong)
	var nik string
	err := r.db.QueryRowContext(ctx, `SELECT nik FROM system_user WHERE id = ? LIMIT 1`, userID).Scan(&nik)
	if err != nil {
		// Log jika error, tapi lanjutkan
	}
	if nik == "" {
		warnings.Profile = "Profil belum lengkap"
	}

	// Cek parking (contoh: cek alert open yang terkait user jika ada)
	// Kita asumsikan alert bisa terkait dengan user secara tidak langsung (misalnya via officer atau lokasi yang diakses user)
	// Query ini hanya contoh, bisa disesuaikan dengan logika bisnis sebenarnya
	var openAlertCount int64
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM admin_alert_event
		WHERE alert_status = 'open' AND officer_user_id = ?
	`, userID).Scan(&openAlertCount)
	if err != nil && err != sql.ErrNoRows {
		// Log error
	}
	if openAlertCount > 0 {
		warnings.Parking = fmt.Sprintf("Ada %d alert yang belum ditangani", openAlertCount)
	}

	// Cek finance (contoh: cek transaksi unpaid atau dispute open)
	var unpaidTxCount, openDisputeCount int64
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM financial_parking_transaction
		WHERE customer_user_id = ? AND transaction_status = 'unpaid'
	`, userID).Scan(&unpaidTxCount)
	if err != nil && err != sql.ErrNoRows {
		// Log error
	}

	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM financial_dispute_case
		WHERE opened_by_user_id = ? AND case_status = 'open'
	`, userID).Scan(&openDisputeCount)
	if err != nil && err != sql.ErrNoRows {
		// Log error
	}

	if unpaidTxCount > 0 || openDisputeCount > 0 {
		warnings.Finance = fmt.Sprintf("Ada %d transaksi belum lunas dan %d sengketa aktif", unpaidTxCount, openDisputeCount)
	}

	return warnings, nil
}
