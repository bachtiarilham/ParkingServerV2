package auth

type UserExistRespModel struct {
	EmailExists    bool
	UsernameExists bool
	PhoneExists    bool
}
