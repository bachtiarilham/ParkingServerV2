package profile

import (
	dto "modulegue/internal/data/mobile/remote/dto/profile"
	model "modulegue/internal/domain/mobile/model/profile"
)

func ToCustomerDto(src *model.CustomerModel) *dto.CustomerDto {
	if src == nil {
		return nil
	}
	return &dto.CustomerDto{
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
		ActiveMembershipId:      src.ActiveMembershipId,
		MembershipPackageName:   src.MembershipPackageName,
		MembershipExpiredAt:     src.MembershipExpiredAt,
		MembershipPackageCode:   src.MembershipPackageCode,
		MembershipStatus:        src.MembershipStatus,
		ActiveParkingSessionId:  src.ActiveParkingSession,
		TotalParkingCount:       src.TotalParkingCount,
		TotalPaymentAmount:      src.TotalPaymentAmount,
		UnreadNotificationCount: src.UnreadNotificationCount,
	}
}
