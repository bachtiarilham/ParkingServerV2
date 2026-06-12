package dto

// // --- Register ---

type RegisterRequest struct {
	FullName    string `json:"full_name"`
	NIK         string `json:"nik"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	UserID  int64  `json:"user_id"` // Karena ID bisa besar, kirim sebagai string
}

// --- Login ---

type LoginRequest struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AuthUser AuthUser `json:"auth_user"`
	TokenSet TokenSet `json:"token"`
}

type AuthUser struct {
	UserId     int64  `json:"user_id"`
	FullName   string `json:"full_name"`
	Nik        string `json:"nik"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Role       int64  `json:"role_id"`
	IsVerified bool   `json:"is_verified"`
}

type TokenSet struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresInSeconds int64  `json:"expires_at"`
}

type RefreshTokenRequestDto struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponseDto struct {
	AuthUser `json:",inline"` // Embed AuthUser
	TokenSet `json:",inline"` // Embed TokenSet
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ChangePasswordResponse struct {
	Message string `json:"message"`
}
