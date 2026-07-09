package auth

type RegisterRequestDto struct {
	FullName string `json:"full_name"`
	NIK      string `json:"nik"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
