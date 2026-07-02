package auth

import (
	dto "modulegue/internal/data/mobile/remote/dto/helper"
)

type UserDto struct {
	UserId     int64          `json:"user_id"`
	Nik        string         `json:"nik"`
	FullName   string         `json:"full_name"`
	Phone      string         `json:"phone"`
	Email      string         `json:"email"`
	Username   string         `json:"username"`
	Password   string         `json:"password"`
	RoleId     int64          `json:"role_id"`
	IsVerified bool           `json:"is_verified"`
	Lokasi     string         `json:"lokasi"`
	Zona       string         `json:"zona"`
	Tarif      []dto.TarifDto `json:"tarif"`
}
