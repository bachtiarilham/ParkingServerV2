package usecase

import (
	"context"
	"fmt"
	"modulegue/internal/domain/user" // Import domain user
)

// GetProfileInput tidak perlu input selain userID dari context
type GetProfileInput struct {
	UserID int64
}

// GetProfileOutput sesuai dengan UserProfileDto
type GetProfileOutput struct {
	UserID     int64  `json:"user_id"`
	Nik        string `json:"nik"`
	FullName   string `json:"full_name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	RoleID     int64  `json:"role"` // Sesuaikan dengan field di UserProfileDto: val role: Long?
	IsVerified bool   `json:"is_verified"`
}

type GetProfileUseCase struct {
	userRepo user.Repository
}

func NewGetProfileUseCase(userRepo user.Repository) *GetProfileUseCase {
	return &GetProfileUseCase{
		userRepo: userRepo,
	}
}

func (uc *GetProfileUseCase) Execute(ctx context.Context, input GetProfileInput) (GetProfileOutput, error) {
	// Cari user berdasarkan ID dari context
	u, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return GetProfileOutput{}, fmt.Errorf("get profile: %w", err)
	}

	// Map dari domain entity ke output usecase
	output := GetProfileOutput{
		UserID:     u.ID,
		Nik:        u.Nik,
		FullName:   u.FullName,
		Phone:      u.Phone,
		Email:      u.Email,
		Username:   u.Username,
		RoleID:     u.RoleID, // Kirim RoleID sebagai Long (int64)
		IsVerified: u.IsVerified,
	}

	return output, nil
}
