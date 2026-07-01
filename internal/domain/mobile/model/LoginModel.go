package model

import "time"

type LoginModel struct {
	ID           int64     `json:"id"`
	FullName     string    `json:"full_name"`
	RoleId       int64     `json:"role_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
}
