package profile

import (
	dto "modulegue/internal/data/mobile/remote/dto/profile"
	model "modulegue/internal/domain/mobile/model/profile"
)

func ToJukirDto(src *model.JukirModel) *dto.JukirDto {
	if src == nil {
		return nil
	}
	return &dto.JukirDto{
		UserId:                  src.UserId,
		Nik:                     src.Nik,
		FullName:                src.FullName,
		Username:                src.Username,
		Email:                   src.Email,
		Phone:                   src.Phone,
		PhotoUrl:                src.PhotoUrl,
		IsVerified:              src.IsVerified,
		RoleId:                  src.RoleId,
		RoleCode:                src.RoleCode,
		RoleName:                src.RoleName,
		Saldo:                   src.Saldo,
		LocationId:              src.LocationId,
		LocationCode:            src.LocationCode,
		LocationName:            src.LocationName,
		Address:                 src.Address,
		MinLatitude:             src.MinLatitude,
		MaxLatitude:             src.MaxLatitude,
		MinLongitude:            src.MinLongitude,
		MaxLongitude:            src.MaxLongitude,
		CenterLatitude:          src.CenterLatitude,
		CenterLongitude:         src.CenterLongitude,
		RadiusMeter:             src.RadiusMeter,
		AreaId:                  src.AreaId,
		AreaName:                src.AreaName,
		ZoneId:                  src.ZoneId,
		ZoneName:                src.ZoneName,
		AssignmentEffectiveFrom: src.AssignmentEffectiveFrom,
		AssignmentEffectiveTo:   src.AssignmentEffectiveTo,
		TodayIncome:             src.TodayIncome,
		TotalIncome:             src.TotalIncome,
		TodayTransactionCount:   src.TodayTransactionCount,
		UnreadNotificationCount: src.UnreadNotificationCount,
	}
}
