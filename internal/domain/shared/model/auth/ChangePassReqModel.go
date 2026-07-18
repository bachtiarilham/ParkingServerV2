package auth

type ChangePassReqModel struct {
	UserId      int64
	OldPassword string
	NewPassword string
}
