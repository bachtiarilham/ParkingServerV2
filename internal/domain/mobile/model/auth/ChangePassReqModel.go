package auth

type ChangePassReqModel struct {
	Username    string
	Email       string
	Phone       string
	OldPassword string
	NewPassword string
}
