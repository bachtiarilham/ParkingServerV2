package auth

type RegisterResponseDto struct {
	Message string  `json:"message"`
	User    UserDto `json:"user"`
}
