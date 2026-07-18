package auth

type LogoutRequestDto struct {
	RefreshToken string `json:"refresh_token"`
}
