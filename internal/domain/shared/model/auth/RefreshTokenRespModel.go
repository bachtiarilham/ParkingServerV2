package auth

type RefreshTokenRespModel struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64
}
