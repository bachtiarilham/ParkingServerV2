package dto

type ChangePasswordRequestDto struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ChangePasswordResponseDto struct {
	Message string `json:"message"`
}
