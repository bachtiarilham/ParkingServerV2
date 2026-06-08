package auth

import (
	"context"
	"errors"
	"fmt"

	"modulegue/internal/domain/user"
	"modulegue/pkg/hash"
)

var ErrEmailAlreadyExists = errors.New("email sudah terdaftar")

type RegisterRequest struct {
	Name     string
	Email    string
	Password string
}

type RegisterUseCase struct {
	userRepo user.Repository
}

func NewRegisterUseCase(userRepo user.Repository) *RegisterUseCase {
	return &RegisterUseCase{userRepo: userRepo}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, req RegisterRequest) error {
	// Cek email sudah terdaftar
	existing, _ := uc.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return ErrEmailAlreadyExists
	}

	// Hash password
	hashed, err := hash.Hash(req.Password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	return uc.userRepo.Create(ctx, &user.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashed,
	})
}
