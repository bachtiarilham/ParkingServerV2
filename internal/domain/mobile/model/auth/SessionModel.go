package auth

import "time"

type SessionModel struct {
	UserID       int64  `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
	DeviceId     string
	DeviceName   string
	FcmToken     string
	IpAddress    string
	UserAgent    string
	ExpiresAt    time.Time `json:"expires_at"`
	RevokedAt    time.Time
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (s SessionModel) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
