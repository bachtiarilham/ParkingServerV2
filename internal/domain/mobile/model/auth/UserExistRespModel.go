package auth

type UserExistRespModel struct {
	NikExists      bool
	EmailExists    bool
	UsernameExists bool
	PhoneExists    bool
}
