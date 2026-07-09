package auth

type RefreshTokenResponseDto struct {
	TokenSetDto `json:"TokenSetDto"` // Embed TokenSet
}
