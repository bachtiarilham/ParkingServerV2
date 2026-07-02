package auth

type LoginRequestModel struct {
	Email    string `json:"identity"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
