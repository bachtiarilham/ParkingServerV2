package profile

import (
	"time"
)

type CustomerDto struct {
	UserId                  int64     `json:"user_id"`
	Nik                     string    `json:"nik"`
	FullName                string    `json:"full_name"`
	Username                string    `json:"username"`
	Email                   string    `json:"email"`
	Phone                   string    `json:"phone"`
	PhotoUrl                string    `json:"photo_url"`
	IsVerified              bool      `json:"is_verified"`
	RoleId                  int64     `json:"role_id"`
	RoleCode                string    `json:"role_code"`
	RoleName                string    `json:"role_name"`
	Saldo                   int64     `json:"saldo"`
	ActiveMembershipId      int64     `json:"active_membership_id"`
	MembershipPackageName   string    `json:"membership_package_name"`
	MembershipExpiredAt     time.Time `json:"membership_expired_at"`
	MembershipPackageCode   string    `json:"membership_package_code"`
	MembershipStatus        string    `json:"membership_status"`
	ActiveParkingSession    int64     `json:"active_parking_session"`
	TotalParkingCount       int64     `json:"total_parking_count"`
	TotalPaymentAmount      int64     `json:"total_payment_amount"`
	UnreadNotificationCount int64     `json:"unread_notification_count"`
}
