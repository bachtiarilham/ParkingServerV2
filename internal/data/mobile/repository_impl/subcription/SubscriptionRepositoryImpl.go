package subscription

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	model "modulegue/internal/domain/mobile/model/subscription"
	"modulegue/internal/domain/mobile/repository"
)

type SubscriptionRepositoryImpl struct {
	db *sql.DB
}

type activeMembershipRow struct {
	packageName    string
	packageExpired time.Time
}

type packageRow struct {
	packageID       int64
	packageCode     string
	packageName     string
	price           int64
	durationDays    int64
	discountPercent float64
	description     string
	CoverageLokasi  []string
	BenefitPackage  []string
}

type promoRow struct {
	promoName     string
	discountValue float64
}

func NewSubscriptionRepositoryImpl(db *sql.DB) repository.SubscriptionRepository {
	return &SubscriptionRepositoryImpl{db: db}
}

func (r *SubscriptionRepositoryImpl) GetSubscribe(ctx context.Context, userId int64) (*model.SubscriptionResponseModel, error) {
	activeMembership, err := r.getActiveMembership(ctx, userId)
	if err != nil {
		return nil, err
	}

	activeBenefits, err := r.getActiveMembershipBenefits(ctx, userId)
	if err != nil {
		return nil, err
	}

	packages, err := r.getMembershipPackages(ctx)
	if err != nil {
		return nil, err
	}

	packageLocations, err := r.getMembershipPackageLocations(ctx)
	if err != nil {
		return nil, err
	}

	packageBenefits, err := r.getMembershipPackageBenefits(ctx)
	if err != nil {
		return nil, err
	}

	availablePromos, err := r.getAvailablePromos(ctx, userId)
	if err != nil {
		return nil, err
	}

	availablePromoTerms, err := r.getAvailablePromoTerms(ctx, userId)
	if err != nil {
		return nil, err
	}

	for i := range packages {
		pkg := &packages[i]
		pkg.CoverageLokasi = append([]string(nil), packageLocations[pkg.packageID]...)
		pkg.BenefitPackage = append([]string(nil), packageBenefits[pkg.packageID]...)
	}

	listPaket := model.ListPaket{
		Bulanan:   []model.DetailPaket{},
		EnamBulan: []model.DetailPaket{},
		Tahunan:   []model.DetailPaket{},
	}

	for _, pkg := range packages {
		detail := model.DetailPaket{
			NamaPaket:      pkg.packageName,
			Harga:          pkg.price,
			CoverageLokasi: append([]string(nil), pkg.CoverageLokasi...),
			BenefitPackage: append([]string(nil), pkg.BenefitPackage...),
		}

		switch pkg.packageCode {
		case "MONTHLY":
			listPaket.Bulanan = append(listPaket.Bulanan, detail)
		case "SIX_MONTHS":
			listPaket.EnamBulan = append(listPaket.EnamBulan, detail)
		case "YEARLY":
			listPaket.Tahunan = append(listPaket.Tahunan, detail)
		default:
			listPaket.Bulanan = append(listPaket.Bulanan, detail)
		}
	}

	return &model.SubscriptionResponseModel{
		ActivePackageName:    activeMembership.packageName,
		ActivePackageExpired: activeMembership.packageExpired,
		ActivePackageBenefit: append([]string(nil), activeBenefits...),
		ListPaket:            listPaket,
		PromoTersedia: model.PromoTersedia{
			SyaratDanKetentuan: append([]string(nil), availablePromoTerms...),
			EachPromo:          availablePromos,
		},
	}, nil
}

func (r *SubscriptionRepositoryImpl) getActiveMembership(ctx context.Context, userId int64) (activeMembershipRow, error) {
	query := `
SELECT
    mu.id AS userMembershipId,
    mp.id AS packageId,
    mp.package_name AS packageName,
    mp.package_period_code AS packagePeriodCode,
    mu.expired_at AS packageExpired,
    mu.status AS membershipStatus

FROM membership_user mu

JOIN membership_package mp
    ON mp.id = mu.package_id

WHERE mu.user_id = ?
  AND mu.status = 'ACTIVE'
  AND mu.expired_at >= NOW()

ORDER BY mu.expired_at DESC

LIMIT 1;
`

	var (
		userMembershipID int64
		packageID        int64
		packageName      sql.NullString
		packagePeriod    sql.NullString
		packageExpired   sql.NullTime
		membershipStatus sql.NullString
	)

	err := r.db.QueryRowContext(ctx, query, userId).Scan(
		&userMembershipID,
		&packageID,
		&packageName,
		&packagePeriod,
		&packageExpired,
		&membershipStatus,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return activeMembershipRow{
				packageName:    "Belum ada paket aktif",
				packageExpired: time.Time{},
			}, nil
		}
		return activeMembershipRow{}, fmt.Errorf("get active membership: %w", err)
	}

	_ = userMembershipID
	_ = packageID
	_ = packagePeriod
	_ = membershipStatus

	result := activeMembershipRow{
		packageName:    packageName.String,
		packageExpired: time.Time{},
	}
	if packageExpired.Valid {
		result.packageExpired = packageExpired.Time
	}

	return result, nil
}

func (r *SubscriptionRepositoryImpl) getActiveMembershipBenefits(ctx context.Context, userId int64) ([]string, error) {
	query := `
SELECT
    mb.id AS benefitId,
    mb.benefit_name AS benefitName,
    mb.benefit_value AS benefitValue,
    mb.sort_order AS sortOrder

FROM membership_user mu

JOIN membership_package mp
    ON mp.id = mu.package_id

JOIN membership_benefit mb
    ON mb.package_id = mp.id

WHERE mu.user_id = ?
  AND mu.status = 'ACTIVE'
  AND mu.expired_at >= NOW()

ORDER BY
    mb.sort_order ASC,
    mb.id ASC;
`

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("get active membership benefits: %w", err)
	}
	defer rows.Close()

	benefits := []string{}
	for rows.Next() {
		var (
			benefitID    int64
			benefitName  sql.NullString
			benefitValue sql.NullString
			sortOrder    int64
		)

		if err := rows.Scan(&benefitID, &benefitName, &benefitValue, &sortOrder); err != nil {
			return nil, fmt.Errorf("scan active membership benefit: %w", err)
		}
		_ = benefitID
		_ = sortOrder

		benefits = append(benefits, buildLabelValue(benefitName.String, benefitValue.String))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate active membership benefits: %w", err)
	}

	return benefits, nil
}

func (r *SubscriptionRepositoryImpl) getMembershipPackages(ctx context.Context) ([]packageRow, error) {
	query := `
SELECT
    mp.id AS packageId,
    mp.package_period_code AS packagePeriodCode,
    mp.package_name AS packageName,
    mp.price AS price,
    mp.duration_days AS durationDays,
    mp.discount_percent AS discountPercent,
    mp.description AS description

FROM membership_package mp

WHERE mp.is_active = 1

ORDER BY
    CASE mp.package_period_code
        WHEN 'MONTHLY' THEN 1
        WHEN 'SIX_MONTHS' THEN 2
        WHEN 'YEARLY' THEN 3
        ELSE 99
    END,
    mp.price ASC,
    mp.id ASC;
`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get membership packages: %w", err)
	}
	defer rows.Close()

	result := []packageRow{}
	for rows.Next() {
		var (
			row             packageRow
			discountPercent sql.NullFloat64
			durationDays    int64
			packageName     sql.NullString
			packagePeriod   sql.NullString
			description     sql.NullString
		)

		if err := rows.Scan(
			&row.packageID,
			&packagePeriod,
			&packageName,
			&row.price,
			&durationDays,
			&discountPercent,
			&description,
		); err != nil {
			return nil, fmt.Errorf("scan membership package: %w", err)
		}

		row.packageCode = packagePeriod.String
		row.packageName = packageName.String
		row.durationDays = durationDays
		row.discountPercent = discountPercent.Float64
		row.description = description.String
		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate membership packages: %w", err)
	}

	return result, nil
}

func (r *SubscriptionRepositoryImpl) getMembershipPackageLocations(ctx context.Context) (map[int64][]string, error) {
	query := `
SELECT
    mp.id AS packageId,

    lp.id AS locationId,
    lp.location_name AS locationName,
    lp.address AS locationAddress

FROM membership_package mp

JOIN membership_package_location mpl
    ON mpl.package_id = mp.id

JOIN location_parking lp
    ON lp.id = mpl.location_id
   AND lp.is_active = 1

WHERE mp.is_active = 1

ORDER BY
    mp.id ASC,
    lp.location_name ASC;
`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get membership package locations: %w", err)
	}
	defer rows.Close()

	result := map[int64][]string{}
	for rows.Next() {
		var (
			packageID       int64
			locationID      int64
			locationName    sql.NullString
			locationAddress sql.NullString
		)

		if err := rows.Scan(&packageID, &locationID, &locationName, &locationAddress); err != nil {
			return nil, fmt.Errorf("scan membership package location: %w", err)
		}
		_ = locationID

		result[packageID] = append(result[packageID], buildLabelValue(locationName.String, locationAddress.String))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate membership package locations: %w", err)
	}

	return result, nil
}

func (r *SubscriptionRepositoryImpl) getMembershipPackageBenefits(ctx context.Context) (map[int64][]string, error) {
	query := `
SELECT
    mp.id AS packageId,

    mb.id AS benefitId,
    mb.benefit_name AS benefitName,
    mb.benefit_value AS benefitValue,
    mb.sort_order AS sortOrder

FROM membership_package mp

JOIN membership_benefit mb
    ON mb.package_id = mp.id

WHERE mp.is_active = 1

ORDER BY
    mp.id ASC,
    mb.sort_order ASC,
    mb.id ASC;
`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get membership package benefits: %w", err)
	}
	defer rows.Close()

	result := map[int64][]string{}
	for rows.Next() {
		var (
			packageID    int64
			benefitID    int64
			benefitName  sql.NullString
			benefitValue sql.NullString
			sortOrder    int64
		)

		if err := rows.Scan(&packageID, &benefitID, &benefitName, &benefitValue, &sortOrder); err != nil {
			return nil, fmt.Errorf("scan membership package benefit: %w", err)
		}
		_ = benefitID
		_ = sortOrder

		result[packageID] = append(result[packageID], buildLabelValue(benefitName.String, benefitValue.String))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate membership package benefits: %w", err)
	}

	return result, nil
}

func (r *SubscriptionRepositoryImpl) getAvailablePromos(ctx context.Context, userId int64) ([]model.DetailPromo, error) {
	query := `
SELECT
    pc.id AS promoId,
    pc.promo_code AS promoCode,
    pc.promo_name AS promoName,
    pc.description AS description,

    pc.discount_type AS discountType,
    pc.discount_value AS discountValue,

    pc.minimum_purchase_amount AS minimumPurchaseAmount,
    pc.maximum_discount_amount AS maximumDiscountAmount,

    pc.quota_total AS quotaTotal,
    pc.quota_used AS quotaUsed,

    pc.valid_from AS validFrom,
    pc.valid_to AS validTo

FROM user_identity ui

JOIN promo_target_role ptr
    ON ptr.role_id = ui.role_id

JOIN promo_campaign pc
    ON pc.id = ptr.promo_id

WHERE ui.id = ?
  AND ui.status = 'ACTIVE'
  AND pc.is_active = 1
  AND pc.status = 'ACTIVE'
  AND (pc.valid_from IS NULL OR pc.valid_from <= NOW())
  AND (pc.valid_to IS NULL OR pc.valid_to >= NOW())
  AND (
        pc.quota_total IS NULL
        OR pc.quota_used < pc.quota_total
  )

ORDER BY
    pc.discount_value DESC,
    pc.valid_to ASC,
    pc.id ASC;
`

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("get available promos: %w", err)
	}
	defer rows.Close()

	result := []model.DetailPromo{}
	for rows.Next() {
		var (
			promoID               int64
			promoCode             sql.NullString
			promoName             sql.NullString
			description           sql.NullString
			discountType          sql.NullString
			discountValue         sql.NullFloat64
			minimumPurchaseAmount sql.NullFloat64
			maximumDiscountAmount sql.NullFloat64
			quotaTotal            sql.NullInt64
			quotaUsed             sql.NullInt64
			validFrom             sql.NullTime
			validTo               sql.NullTime
		)

		if err := rows.Scan(
			&promoID,
			&promoCode,
			&promoName,
			&description,
			&discountType,
			&discountValue,
			&minimumPurchaseAmount,
			&maximumDiscountAmount,
			&quotaTotal,
			&quotaUsed,
			&validFrom,
			&validTo,
		); err != nil {
			return nil, fmt.Errorf("scan available promo: %w", err)
		}
		_ = promoID
		_ = promoCode
		_ = description
		_ = discountType
		_ = minimumPurchaseAmount
		_ = maximumDiscountAmount
		_ = quotaTotal
		_ = quotaUsed
		_ = validFrom
		_ = validTo

		result = append(result, model.DetailPromo{
			NamaPromo:   promoName.String,
			BesarDiskon: discountValue.Float64,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate available promos: %w", err)
	}

	return result, nil
}

func (r *SubscriptionRepositoryImpl) getAvailablePromoTerms(ctx context.Context, userId int64) ([]string, error) {
	query := `
SELECT DISTINCT
    ptc.term_text AS termText,
    ptc.sort_order AS sortOrder

FROM user_identity ui

JOIN promo_target_role ptr
    ON ptr.role_id = ui.role_id

JOIN promo_campaign pc
    ON pc.id = ptr.promo_id

JOIN promo_term_condition ptc
    ON ptc.promo_id = pc.id

WHERE ui.id = ?
  AND ui.status = 'ACTIVE'
  AND pc.is_active = 1
  AND pc.status = 'ACTIVE'
  AND (pc.valid_from IS NULL OR pc.valid_from <= NOW())
  AND (pc.valid_to IS NULL OR pc.valid_to >= NOW())
  AND (
        pc.quota_total IS NULL
        OR pc.quota_used < pc.quota_total
  )

ORDER BY
    ptc.sort_order ASC,
    ptc.term_text ASC;
`

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("get available promo terms: %w", err)
	}
	defer rows.Close()

	result := []string{}
	for rows.Next() {
		var (
			termText  sql.NullString
			sortOrder int64
		)

		if err := rows.Scan(&termText, &sortOrder); err != nil {
			return nil, fmt.Errorf("scan available promo term: %w", err)
		}
		_ = sortOrder

		result = append(result, termText.String)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate available promo terms: %w", err)
	}

	return result, nil
}

func buildLabelValue(label string, value string) string {
	switch {
	case label == "" && value == "":
		return ""
	case value == "":
		return label
	case label == "":
		return value
	default:
		return label + ": " + value
	}
}
