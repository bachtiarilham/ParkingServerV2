package auth

type RefreshTokenRequestDto struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponseDto struct {
	UserDto     `json:",inline"` // Embed AuthUser
	TokenSetDto `json:",inline"` // Embed TokenSet
}
