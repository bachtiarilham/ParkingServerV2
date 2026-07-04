package auth

type RefreshTokenRequestDto struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponseDto struct {
	TokenSetDto `json:"TokenSetDto"` // Embed TokenSet
}
