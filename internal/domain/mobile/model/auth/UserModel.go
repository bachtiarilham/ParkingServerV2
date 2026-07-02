package auth

import (
	domain "modulegue/internal/domain/mobile/model/helper"
	"time"
)

type UserModel struct {
	UserId       int64               `json:"user_id"`
	Nik          string              `json:"nik"`
	FullName     string              `json:"full_name"`
	Phone        string              `json:"phone"`
	Email        string              `json:"email"`
	Username     string              `json:"username"`
	Password     string              `json:"password"`
	PasswordHash string              `json:"password_hash"`
	RoleId       int64               `json:"role_id"`
	IsVerified   bool                `json:"is_verified"`
	Lokasi       string              `json:"lokasi"`
	Zona         string              `json:"zona"`
	Tarif        []domain.TarifModel `json:"tarif"`
	RegisteredAt time.Time           `json:"registered_at"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}
