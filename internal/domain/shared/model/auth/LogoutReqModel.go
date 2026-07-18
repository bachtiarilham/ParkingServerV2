package auth

type LogoutReqModel struct {
	// Kita bisa menerima access token untuk extract user_id dan/atau refresh token
	// Tapi karena kita menggunakan middleware otentikasi, user_id bisa diambil dari context
	// Dan refresh token bisa diambil dari body request (jika ingin revoke refresh token)
	UserId       int64
	AccessToken  string
	RefreshToken string // Opsional: jika ingin revoke refresh token juga
}
