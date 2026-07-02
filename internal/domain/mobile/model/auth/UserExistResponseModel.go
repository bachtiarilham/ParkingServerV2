package auth

type UserExistResult struct {
	EmailExists    bool
	UsernameExists bool
	PhoneExists    bool
}
