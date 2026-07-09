package auth

type ChangePasswordRequestDto struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
