package home

import (
	"context"
	"database/sql"
	"fmt"

	model "modulegue/internal/domain/mobile/model/home"
	profileModel "modulegue/internal/domain/mobile/model/profile"
	"modulegue/internal/domain/mobile/repository"
)

type HomeRepositoryImpl struct {
	db *sql.DB
}

func NewHomeRepositoryImpl(db *sql.DB) repository.HomeRepository {
	return &HomeRepositoryImpl{db: db}
}

func (r *HomeRepositoryImpl) GetJukirHome(ctx context.Context, reqModel model.GetHomeReqModel) (*model.JukirHomeModel, error) {
	var profile *profileModel.JukirModel
	var err error

	profile, err = r.getJukirProfile(ctx, reqModel.UserID)
	if err != nil {
		return nil, err
	}

	contents, err := r.getContents(ctx, reqModel.UserID)
	if err != nil {
		return nil, err
	}

	warnings, err := r.getUnreadNotifCounts(ctx, reqModel.UserID)
	if err != nil {
		return nil, err
	}

	return &model.JukirHomeModel{
		Profile:          profile,
		Contents:         &contents,
		UnreadNotifCount: warnings,
	}, nil
}

func (r *HomeRepositoryImpl) GetCustomerHome(ctx context.Context, reqModel model.GetHomeReqModel) (*model.CustomerHomeModel, error) {
	var profile *profileModel.CustomerModel
	var err error

	profile, err = r.getCustomerProfile(ctx, reqModel.UserID)
	if err != nil {
		return nil, err
	}

	contents, err := r.getContents(ctx, reqModel.UserID)
	if err != nil {
		return nil, err
	}

	warnings, err := r.getUnreadNotifCounts(ctx, reqModel.UserID)
	if err != nil {
		return nil, err
	}

	return &model.CustomerHomeModel{
		Profile:          profile,
		Contents:         &contents,
		UnreadNotifCount: warnings,
	}, nil
}

func (r *HomeRepositoryImpl) getCustomerProfile(ctx context.Context, userID int64) (*profileModel.CustomerModel, error) {
	query := `
	SELECT
		ui.id AS user_id,
		ui.full_name,
		ui.username,
		ui.email,
		ui.phone_number,
		ui.photo_url,

		mr.role_code,
		mr.role_name,

		COALESCE(suh.saldo, 0) AS saldo,

		sch.active_membership_id,
		sch.membership_package_name,
		sch.membership_expired_at,

		mp.package_code AS membership_package_code,
		mp.package_name AS membership_package_real_name,
		mu.status AS membership_status,

		sch.active_parking_session_id,
		sch.total_parking_count,
		COALESCE(sch.total_payment_amount, 0) AS total_payment_amount,

		COALESCE(unread_notif.unread_count, 0) AS unread_notification_count

	FROM user_identity ui

	JOIN master_role mr
		ON mr.id = ui.role_id

	LEFT JOIN summary_user_home suh
		ON suh.user_id = ui.id

	LEFT JOIN summary_customer_home sch
		ON sch.user_id = ui.id

	LEFT JOIN membership_user mu
		ON mu.id = sch.active_membership_id

	LEFT JOIN membership_package mp
		ON mp.id = mu.package_id

	LEFT JOIN (
		SELECT
			user_id,
			COUNT(*) AS unread_count
		FROM notification_user
		WHERE is_read = 0
		GROUP BY user_id
	) unread_notif
		ON unread_notif.user_id = ui.id

	WHERE ui.id = ?
	AND mr.role_code = 'CUSTOMER'
	AND ui.status = 'ACTIVE'

	LIMIT 1;
	`

	var row customerProfileRow

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&row.UserID,
		&row.FullName,
		&row.Username,
		&row.Email,
		&row.Phone,
		&row.PhotoURL,
		&row.RoleCode,
		&row.RoleName,
		&row.Saldo,
		&row.ActiveMembershipID,
		&row.MembershipPackageName,
		&row.MembershipExpiredAt,
		&row.MembershipPackageCode,
		&row.MembershipStatus,
		&row.ActiveParkingSessionID,
		&row.TotalParkingCount,
		&row.TotalPaymentAmount,
		&row.UnreadNotificationCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("get customer profile: %w", err)
	}

	return mapCustomerProfileRow(row), nil

}

func (r *HomeRepositoryImpl) getJukirProfile(ctx context.Context, userID int64) (*profileModel.JukirModel, error) {
	query := `
		SELECT
			ui.id AS user_id,
			ui.nik,
			ui.full_name,
			ui.username,
			ui.email,
			ui.phone_number,
			ui.photo_url,
			ui.is_verified,
			ui.role_id,

			mr.role_code,
			mr.role_name,

			COALESCE(suh.saldo, 0) AS saldo,

			aoc.location_id,
			lp.location_code,
			COALESCE(soh.location_name, lp.location_name) AS location_name,
			lp.address,
			lp.center_latitude,
			lp.center_longitude,
			lp.radius_meter,

			aoc.area_id,
			COALESCE(soh.area_name, la.area_name) AS area_name,

			aoc.zone_id,
			COALESCE(soh.zone_name, lz.zone_name) AS zone_name,

			aoc.effective_from AS assignment_effective_from,
			aoc.effective_to AS assignment_effective_to,

			COALESCE(soh.today_income, 0) AS today_income,
			COALESCE(soh.total_income, 0) AS total_income,
			COALESCE(soh.today_transaction_count, 0) AS today_transaction_count,

			COALESCE(unread_notif.unread_count, 0) AS unread_notification_count

		FROM user_identity ui

		JOIN master_role mr
			ON mr.id = ui.role_id

		LEFT JOIN summary_user_home suh
			ON suh.user_id = ui.id

		LEFT JOIN assignment_officer_current aoc
			ON aoc.officer_user_id = ui.id

		LEFT JOIN location_parking lp
			ON lp.id = aoc.location_id

		LEFT JOIN location_area la
			ON la.id = aoc.area_id

		LEFT JOIN location_zone lz
			ON lz.id = aoc.zone_id

		LEFT JOIN summary_officer_home soh
			ON soh.user_id = ui.id

		LEFT JOIN (
			SELECT
				user_id,
				COUNT(*) AS unread_count
			FROM notification_user
			WHERE is_read = 0
			GROUP BY user_id
		) unread_notif
			ON unread_notif.user_id = ui.id

		WHERE ui.id = ?
		AND mr.role_code = 'OFFICER'
		AND ui.status = 'ACTIVE'

	LIMIT 1;
	`

	var row jukirProfileRow

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&row.UserID,
		&row.Nik,
		&row.FullName,
		&row.Username,
		&row.Email,
		&row.Phone,
		&row.PhotoURL,
		&row.IsVerified,
		&row.RoleID,
		&row.RoleCode,
		&row.RoleName,
		&row.Saldo,
		&row.LocationID,
		&row.LocationCode,
		&row.LocationName,
		&row.Address,
		&row.CenterLatitude,
		&row.CenterLongitude,
		&row.RadiusMeter,
		&row.AreaID,
		&row.AreaName,
		&row.ZoneID,
		&row.ZoneName,
		&row.AssignmentEffectiveFrom,
		&row.AssignmentEffectiveTo,
		&row.TodayIncome,
		&row.TotalIncome,
		&row.TodayTransactionCount,
		&row.UnreadNotificationCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("jukir not found")
		}
		return nil, fmt.Errorf("get jukir profile: %w", err)
	}

	return mapJukirProfileRow(row), nil

}

func (r *HomeRepositoryImpl) getContents(ctx context.Context, userId int64) ([]model.ContentsModel, error) {
	query := `
		SELECT
			ci.id AS contentId,

			mct.id AS contentTypeId,
			mct.content_type_code AS contentTypeCode,
			mct.content_type_name AS contentTypeName,

			ci.title,
			ci.summary,
			ci.body,
			ci.thumbnail_url AS thumbnailUrl,
			ci.banner_url AS bannerUrl,

			ci.event_location AS eventLocation,
			ci.event_start_at AS eventStartAt,
			ci.event_end_at AS eventEndAt,

			ci.publish_at AS publishAt,
			ci.expired_at AS expiredAt,
			ci.priority

		FROM user_identity ui

		JOIN content_target_role ctr
			ON ctr.role_id = ui.role_id

		JOIN content_item ci
			ON ci.id = ctr.content_id

		JOIN master_content_type mct
			ON mct.id = ci.content_type_id

		WHERE ui.id = ?
		AND ci.status = 'PUBLISHED'
		AND mct.is_active = 1
		AND (ci.publish_at IS NULL OR ci.publish_at <= NOW())
		AND (ci.expired_at IS NULL OR ci.expired_at >= NOW())

		ORDER BY
			ci.priority DESC,
			ci.publish_at DESC,
			ci.created_at DESC
	
		LIMIT 6
	`

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("get events and news: %w", err)
	}
	defer rows.Close()

	contents := make([]model.ContentsModel, 0)

	for rows.Next() {
		var row contentsRow
		if err := rows.Scan(
			&row.ContentId,
			&row.ContentTypeId,
			&row.ContentTypeCode,
			&row.ContentTypeName,
			&row.Title,
			&row.Summary,
			&row.Body,
			&row.ThumbnailUrl,
			&row.BannerUrl,
			&row.EventLocation,
			&row.EventStartAt,
			&row.EventEndAt,
			&row.PublishAt,
			&row.ExpiredAt,
			&row.Priority,
		); err != nil {
			return nil, fmt.Errorf("scan contents: %w", err)
		}

		contents = append(contents, mapContentsRowToModel(row))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate contents: %w", err)
	}

	return contents, nil
}

func (r *HomeRepositoryImpl) getUnreadNotifCounts(ctx context.Context, userID int64) (int64, error) {
	var notifCount int64

	var unreadNotifCount sql.NullInt64
	if err := r.db.QueryRowContext(
		ctx,
		`
		SELECT
			COUNT(*) AS unreadNotificationCount
		FROM notification_user
		WHERE user_id = ?
		AND is_read = 0;
		`,
		userID,
	).Scan(&unreadNotifCount); err == nil {
		notifCount = 0
	}
	return notifCount, nil
}
