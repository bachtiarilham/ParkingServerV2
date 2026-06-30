package dto

type UserDto struct {
	UserId     int64  `json:"user_id"`
	Nik        string `json:"nik"`
	FullName   string `json:"full_name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Role       int64  `json:"role"` // <-- Sesuaikan dengan Kotlin: val role: Long?
	IsVerified bool   `json:"is_verified"`
}
