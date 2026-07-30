package paymentgate

import "time"

type MemberCheckResult struct {
	IsMember    bool       `json:"is_member"`
	PackageID   *int64     `json:"package_id"`
	PackageName *string    `json:"package_name"`
	ExpiredAt   *time.Time `json:"expired_at"`
}
