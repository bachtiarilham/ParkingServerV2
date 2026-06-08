package auth

type TokenManager interface {
	Generate(userID int64) (string, error)
	Validate(token string) (int64, error)
}
