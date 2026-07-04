package subscription

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	model "modulegue/internal/domain/mobile/model/subscription"
	"modulegue/internal/domain/mobile/repository"
)

type SubscriptionRepositoryImpl struct {
	db *sql.DB
}

func NewSubscriptionRepositoryImpl(db *sql.DB) repository.SubscriptionRepository {
	return &SubscriptionRepositoryImpl{db: db}
}

func (r *SubscriptionRepositoryImpl) GetSubscribe(ctx context.Context, userId int64) (*model.SubscribeModel, error) {
	statusCard, err := r.getStatusCard(ctx, userId)
	if err != nil {
		return nil, err
	}

	packageCards, err := r.getPackageCards(ctx)
	if err != nil {
		return nil, err
	}

	return &model.SubscribeModel{
		StatusCard:  statusCard,
		PackageCard: packageCards,
		Promo:       []model.PromoModel{},
	}, nil
}

func (r *SubscriptionRepositoryImpl) getStatusCard(ctx context.Context, userId int64) (*model.StatusCardModel, error) {
	query := `
		SELECT
			mp.name,
			cm.end_date
		FROM customer_memberships cm
		JOIN membership_plans mp ON mp.id = cm.plan_id
		WHERE cm.customer_user_id = ?
		  AND cm.status = 'active'
		ORDER BY cm.end_date DESC
		LIMIT 1
	`

	var (
		paketAktif sql.NullString
		endDate    sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query, userId).Scan(&paketAktif, &endDate)
	if err != nil {
		if err == sql.ErrNoRows {
			defaultPackage := "Belum ada paket aktif"
			defaultExpiry := "-"
			defaultBenefit := "Aktifkan membership untuk mendapatkan benefit parkir"
			return &model.StatusCardModel{
				PaketAktif: &defaultPackage,
				Kadaluarsa: &defaultExpiry,
				Benefit:    &defaultBenefit,
			}, nil
		}
		return nil, fmt.Errorf("get subscription status card: %w", err)
	}

	paket := paketAktif.String
	kadaluarsa := endDate.Time.Format("2006-01-02")
	benefit := "Membership aktif dapat digunakan selama periode berlangganan"

	return &model.StatusCardModel{
		PaketAktif: &paket,
		Kadaluarsa: &kadaluarsa,
		Benefit:    &benefit,
	}, nil
}

func (r *SubscriptionRepositoryImpl) getPackageCards(ctx context.Context) ([]model.PackageCardModel, error) {
	query := `
		SELECT
			id,
			name,
			price,
			duration_days
		FROM membership_plans
		WHERE is_active = 1
		ORDER BY price ASC, duration_days ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get subscription package cards: %w", err)
	}
	defer rows.Close()

	result := []model.PackageCardModel{}
	for rows.Next() {
		var (
			id           int64
			name         string
			price        int64
			durationDays int
		)

		if err := rows.Scan(&id, &name, &price, &durationDays); err != nil {
			return nil, fmt.Errorf("scan subscription package card: %w", err)
		}

		masaBerlaku := strconv.Itoa(durationDays) + " hari"
		jumlahDiskon := int64(0)
		deskripsi := "Paket membership parkir"
		benefits := []string{
			"Akses membership aktif",
			"Durasi " + masaBerlaku,
		}

		result = append(result, model.PackageCardModel{
			NamaPaket:    &name,
			Harga:        &price,
			MasaBerlaku:  &masaBerlaku,
			JumlahDiskon: &jumlahDiskon,
			Deskripsi:    &deskripsi,
			Benefit:      benefits,
		})
	}

	return result, nil
}
