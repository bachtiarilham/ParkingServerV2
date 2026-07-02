package auth

type LoginRequestDto struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type LoginResponseDto struct {
	UserDto     UserDto     `json:"userDto"`
	TokenSetDto TokenSetDto `json:"token"`
}
