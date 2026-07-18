package auth

type LoginRespModel struct {
	UserId       int64
	RoleId       int64
	PasswordHash string
}
