package model

import "time"

type SessionModel struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
}

func (s SessionModel) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
