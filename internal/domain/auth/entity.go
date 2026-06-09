//Token, Session

package auth

import "time"

type Credential struct {
	ID           int64
	FullName     string
	Role         string
	Email        string
	PasswordHash string
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

type Session struct {
	ID           int64
	UserID       int64
	RefreshToken string
	ExpiresAt    time.Time
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

func (s Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
