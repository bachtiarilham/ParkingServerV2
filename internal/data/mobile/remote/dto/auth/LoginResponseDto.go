package auth

type LoginResponseDto struct {
	UserDto     UserDto     `json:"userDto"`
	TokenSetDto TokenSetDto `json:"token"`
}
