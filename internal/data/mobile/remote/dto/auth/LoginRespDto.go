package auth

type LoginRespDto struct {
	TokenSetDto *TokenSetDto `json:"token_set"`
	RoleId      int64        `json:"role_id"`
}
