package auth

type LoginResponseModel struct {
	UserModel     UserModel     `json:"userDto"`
	TokenSetModel TokenSetModel `json:"token"`
}
