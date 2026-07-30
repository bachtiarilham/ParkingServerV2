package subscription

import (
	"context"
	"database/sql"
	"fmt"

	model "modulegue/internal/domain/mobile/model/subscription"
	"modulegue/internal/domain/mobile/repository"
)

type SubscriptionRepositoryImpl struct {
	db *sql.DB
}

func NewSubscriptionRepositoryImpl(db *sql.DB) repository.SubscriptionRepository {
	return &SubscriptionRepositoryImpl{db: db}
}

func (r *SubscriptionRepositoryImpl) GetSubscribe(ctx context.Context, userId int64) (*model.SubscribeResponseModel, error) {
	resp := &model.SubscribeResponseModel{
		Benefits:  []model.BenefitsModel{},
		ListPaket: []model.DetailPaketModel{},
		Faq:       []model.FaqModel{},
	}
	// A. Fetch Active Paket & Benefits
	activeQuery := `
	SELECT 
		mp.package_name AS active_package_name,
		mu.expired_at AS active_package_expired,
		mb.benefit_name AS name,
		mb.benefit_value AS description
	FROM membership_user mu
	JOIN membership_package mp ON mu.package_id = mp.id
	LEFT JOIN membership_benefit mb ON mb.package_id = mp.id
	WHERE mu.user_id = ? 
	AND mu.status = 'ACTIVE' 
	AND mu.expired_at > NOW();
	`
	activeRows, err := r.db.QueryContext(ctx, activeQuery, userId)
	if err != nil {
		return nil, fmt.Errorf("query active paket: %w", err)
	}
	defer activeRows.Close()

	seenBenefits := make(map[string]bool)

	for activeRows.Next() {
		var (
			packageName  string
			expiredAt    string
			benefitName  sql.NullString
			benefitValue sql.NullString
		)
		if err := activeRows.Scan(&packageName, &expiredAt, &benefitName, &benefitValue); err != nil {
			return nil, fmt.Errorf("scan active paket row: %w", err)
		}
		if resp.ActivePaket == nil {
			resp.ActivePaket = &model.ActivePaketModel{
				ActivePackageName:    &packageName,
				ActivePackageExpired: &expiredAt,
			}
		}
		if benefitName.Valid && benefitName.String != "" {
			bName := benefitName.String
			if !seenBenefits[bName] {
				seenBenefits[bName] = true
				resp.Benefits = append(resp.Benefits, model.BenefitsModel{
					Name:        bName,
					Description: benefitValue.String,
				})
			}
		}
	}

	// B. Fetch Statistik
	statistikQuery := `
	SELECT 
		last_month_period,
		last_month_duration_minutes,
		last_month_amount_paid,
		last_month_amount_saved
	FROM user_identity
	WHERE id = ?;
	`
	var (
		period          string
		durationMinutes int
		amountPaid      int64
		amountSaved     int64
	)
	err = r.db.QueryRowContext(ctx, statistikQuery, userId).Scan(&period, &durationMinutes, &amountPaid, &amountSaved)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("query statistik: %w", err)
		}
		resp.Statistik = nil
	} else {
		// Temporarily store raw values inside the Statistik model fields.
		// Usecase will parse these and overwrite with formatted values.
		resp.Statistik = &model.StatistikModel{
			TotalJamParkirBulanLalu:       durationMinutes,
			TotalBiayaParkirBulanLaluText: fmt.Sprintf("%d", amountPaid),
			TotalPersentaseHematText:      fmt.Sprintf("%d", amountSaved),
		}
	}

	// C. Fetch List Paket
	listPaketQuery := `
	SELECT 
		mp.package_name AS name,
		mp.price AS price,
		mp.package_period_code AS period_label,
		mp.description AS info_label,
		mp.discount_percent AS badge_label,
		mb.benefit_name AS benefit
	FROM membership_package mp
	LEFT JOIN membership_benefit mb ON mb.package_id = mp.id
	WHERE mp.is_active = 1
	ORDER BY mp.price ASC;
	`
	paketRows, err := r.db.QueryContext(ctx, listPaketQuery)
	if err != nil {
		return nil, fmt.Errorf("query list paket: %w", err)
	}
	defer paketRows.Close()

	type pkgKey struct {
		Name     string
		Price    int64
		Period   string
		Info     string
		DiscPcnt int
	}
	var pkgOrder []pkgKey
	pkgMap := make(map[pkgKey][]string)

	for paketRows.Next() {
		var (
			name        string
			price       int64
			periodLabel string
			infoLabel   string
			discPercent int
			benefit     sql.NullString
		)
		if err := paketRows.Scan(&name, &price, &periodLabel, &infoLabel, &discPercent, &benefit); err != nil {
			return nil, fmt.Errorf("scan list paket row: %w", err)
		}
		key := pkgKey{
			Name:     name,
			Price:    price,
			Period:   periodLabel,
			Info:     infoLabel,
			DiscPcnt: discPercent,
		}
		if _, exists := pkgMap[key]; !exists {
			pkgOrder = append(pkgOrder, key)
		}
		if benefit.Valid && benefit.String != "" {
			pkgMap[key] = append(pkgMap[key], benefit.String)
		}
	}
	for _, key := range pkgOrder {
		var badgeStr *string
		if key.DiscPcnt > 0 {
			val := fmt.Sprintf("%d", key.DiscPcnt)
			badgeStr = &val
		}
		resp.ListPaket = append(resp.ListPaket, model.DetailPaketModel{
			Name:        key.Name,
			Price:       key.Price,
			PeriodLabel: key.Period, // Temporarily holds package_period_code
			InfoLabel:   key.Info,
			BadgeLabel:  badgeStr, // Temporarily holds discount_percent as string
			Benefits:    pkgMap[key],
		})
	}

	resp.Faq = append(resp.Faq, model.GetDefaultFaqs()...)
	return resp, nil
}
