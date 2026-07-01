package errorstring

import (
	"errors"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidCredentials   = errors.New("email atau password salah")
	ErrOldPasswordMismatch  = errors.New("password lama tidak cocok")
	ErrNewPasswordSameAsOld = errors.New("password baru tidak boleh sama dengan password lama")
	ErrInvalidUserID        = errors.New("user tidak valid")
	ErrInvalidRefreshToken  = errors.New("refresh token tidak valid")
	ErrExpiredRefreshToken  = errors.New("refresh token sudah kadaluarsa")
)
