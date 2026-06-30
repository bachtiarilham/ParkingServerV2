package dto

type LoginRequest struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthUser struct {
	UserID     int64  `json:"user_id"`
	FullName   string `json:"full_name"`
	Phone      string `json:"phone,omitempty"`
	Email      string `json:"email,omitempty"`
	Username   string `json:"username"`
	Role       string `json:"role"`
	IsVerified bool   `json:"is_verified"`
}

type TokenSet struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresInSeconds int64  `json:"expires_in_seconds"`
}

type AuthEnvelope struct {
	User   AuthUser `json:"user"`
	Tokens TokenSet `json:"tokens"`
}
