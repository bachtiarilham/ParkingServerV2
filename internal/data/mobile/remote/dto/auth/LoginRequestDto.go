package auth

type LoginRequestDto struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}
