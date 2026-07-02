package auth

type RegisterResponseModel struct {
	Message string    `json:"message"`
	User    UserModel `json:"user"`
}
