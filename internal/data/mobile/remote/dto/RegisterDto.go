package dto

type RegisterRequestDto struct {
	FullName    string `json:"full_name"`
	NIK         string `json:"nik"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

type RegisterResponseDto struct {
	Message string `json:"message"`
	UserID  int64  `json:"user_id"` // Karena ID bisa besar, kirim sebagai string
}
