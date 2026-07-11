package home

import (
	"database/sql"

	"modulegue/core/utils"
	profileModel "modulegue/internal/domain/mobile/model/profile"
)

type customerProfileRow struct {
	UserID                  sql.NullInt64
	FullName                sql.NullString
	Username                sql.NullString
	Email                   sql.NullString
	Phone                   sql.NullString
	PhotoURL                sql.NullString
	RoleCode                sql.NullString
	RoleName                sql.NullString
	Saldo                   sql.NullInt64
	ActiveMembershipID      sql.NullInt64
	MembershipPackageName   sql.NullString
	MembershipExpiredAt     sql.NullTime
	MembershipPackageCode   sql.NullString
	MembershipStatus        sql.NullString
	ActiveParkingSessionID  sql.NullInt64
	TotalParkingCount       sql.NullInt64
	TotalPaymentAmount      sql.NullInt64
	UnreadNotificationCount sql.NullInt64
}

type jukirProfileRow struct {
	UserID                  sql.NullInt64
	Nik                     sql.NullString
	FullName                sql.NullString
	Username                sql.NullString
	Email                   sql.NullString
	Phone                   sql.NullString
	PhotoURL                sql.NullString
	IsVerified              sql.NullBool
	RoleID                  sql.NullInt64
	RoleCode                sql.NullString
	RoleName                sql.NullString
	Saldo                   sql.NullInt64
	LocationID              sql.NullInt64
	LocationCode            sql.NullString
	LocationName            sql.NullString
	Address                 sql.NullString
	CenterLatitude          sql.NullFloat64
	CenterLongitude         sql.NullFloat64
	RadiusMeter             sql.NullInt64
	AreaID                  sql.NullInt64
	AreaName                sql.NullString
	ZoneID                  sql.NullInt64
	ZoneName                sql.NullString
	AssignmentEffectiveFrom sql.NullTime
	AssignmentEffectiveTo   sql.NullTime
	TodayIncome             sql.NullInt64
	TotalIncome             sql.NullInt64
	TodayTransactionCount   sql.NullInt64
	UnreadNotificationCount sql.NullInt64
}

func mapCustomerProfileRow(row customerProfileRow) *profileModel.CustomerModel {
	profile := &profileModel.CustomerModel{}
	profile.UserId = utils.NullInt64Value(row.UserID)
	profile.FullName = utils.NullStringValue(row.FullName)
	profile.Username = utils.NullStringValue(row.Username)
	profile.Email = utils.NullStringValue(row.Email)
	profile.Phone = utils.NullStringValue(row.Phone)
	profile.PhotoUrl = utils.NullStringValue(row.PhotoURL)
	profile.RoleCode = utils.NullStringValue(row.RoleCode)
	profile.RoleName = utils.NullStringValue(row.RoleName)
	profile.Saldo = utils.NullInt64Value(row.Saldo)
	profile.ActiveMembershipId = utils.NullInt64Value(row.ActiveMembershipID)
	profile.MembershipPackageName = utils.NullStringValue(row.MembershipPackageName)
	profile.MembershipExpiredAt = utils.NullTimeValue(row.MembershipExpiredAt)
	profile.MembershipPackageCode = utils.NullStringValue(row.MembershipPackageCode)
	profile.MembershipStatus = utils.NullStringValue(row.MembershipStatus)
	profile.ActiveParkingSession = utils.NullInt64Value(row.ActiveParkingSessionID)
	profile.TotalParkingCount = utils.NullInt64Value(row.TotalParkingCount)
	profile.TotalPaymentAmount = utils.NullInt64Value(row.TotalPaymentAmount)
	profile.UnreadNotificationCount = utils.NullInt64Value(row.UnreadNotificationCount)
	return profile
}

func mapJukirProfileRow(row jukirProfileRow) *profileModel.JukirModel {
	profile := &profileModel.JukirModel{}
	profile.UserId = utils.NullInt64Value(row.UserID)
	profile.Nik = utils.NullStringValue(row.Nik)
	profile.FullName = utils.NullStringValue(row.FullName)
	profile.Username = utils.NullStringValue(row.Username)
	profile.Email = utils.NullStringValue(row.Email)
	profile.Phone = utils.NullStringValue(row.Phone)
	profile.PhotoUrl = utils.NullStringValue(row.PhotoURL)
	profile.IsVerified = utils.NullBoolValue(row.IsVerified)
	profile.RoleId = utils.NullInt64Value(row.RoleID)
	profile.RoleCode = utils.NullStringValue(row.RoleCode)
	profile.RoleName = utils.NullStringValue(row.RoleName)
	profile.Saldo = utils.NullInt64Value(row.Saldo)
	profile.LocationId = utils.NullInt64Value(row.LocationID)
	profile.LocationCode = utils.NullStringValue(row.LocationCode)
	profile.LocationName = utils.NullStringValue(row.LocationName)
	profile.Address = utils.NullStringValue(row.Address)
	profile.CenterLatitude = utils.NullFloat64Value(row.CenterLatitude)
	profile.CenterLongitude = utils.NullFloat64Value(row.CenterLongitude)
	profile.RadiusMeter = utils.NullInt64Value(row.RadiusMeter)
	profile.AreaId = utils.NullInt64Value(row.AreaID)
	profile.AreaName = utils.NullStringValue(row.AreaName)
	profile.ZoneId = utils.NullInt64Value(row.ZoneID)
	profile.ZoneName = utils.NullStringValue(row.ZoneName)
	profile.AssignmentEffectiveFrom = utils.NullTimeValue(row.AssignmentEffectiveFrom)
	profile.AssignmentEffectiveTo = utils.NullTimeValue(row.AssignmentEffectiveTo)
	profile.TodayIncome = utils.NullInt64Value(row.TodayIncome)
	profile.TotalIncome = utils.NullInt64Value(row.TotalIncome)
	profile.TodayTransactionCount = utils.NullInt64Value(row.TodayTransactionCount)
	profile.UnreadNotificationCount = utils.NullInt64Value(row.UnreadNotificationCount)
	return profile
}
