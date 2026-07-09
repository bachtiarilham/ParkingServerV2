package auth

import (
	"context"
	"errors"
	"strings"

	model "modulegue/internal/domain/mobile/model/auth"
	"modulegue/internal/domain/mobile/repository"
)

var (
	ErrInvalidInput          = errors.New("tolong isi field yang wajib")
	ErrNikAlreadyExists      = errors.New("NIK sudah terdaftar")
	ErrEmailAlreadyExists    = errors.New("Email sudah terdaftar")
	ErrPhoneAlreadyExists    = errors.New("Phone sudah terdaftar")
	ErrUsernameAlreadyExists = errors.New("Username sudah terdaftar")
)

type RegisterUseCase struct {
	AuthRepo repository.AuthRepository
}

func NewRegisterUseCase(
	AuthRepo repository.AuthRepository,
) *RegisterUseCase {
	return &RegisterUseCase{
		AuthRepo: AuthRepo,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, reqModel model.RegisterRequestModel) error {
	reqModel.FullName = strings.TrimSpace(reqModel.FullName)
	reqModel.NIK = strings.TrimSpace(reqModel.NIK)
	reqModel.Phone = strings.TrimSpace(reqModel.Phone)
	reqModel.Email = strings.ToLower(strings.TrimSpace(reqModel.Email))
	reqModel.Username = strings.TrimSpace(reqModel.Username)
	reqModel.Password = strings.TrimSpace(reqModel.Password)

	if reqModel.Username == "" || reqModel.Email == "" || reqModel.Password == "" {
		return ErrInvalidInput
	}

	exists, err := uc.AuthRepo.ExistsByIdentity(ctx, reqModel.NIK, reqModel.Email, reqModel.Username, reqModel.Phone)
	if err != nil {
		return err
	}
	if exists.NikExists {
		return ErrNikAlreadyExists
	}
	if exists.EmailExists {
		return ErrEmailAlreadyExists
	}
	if exists.UsernameExists {
		return ErrUsernameAlreadyExists
	}
	if exists.PhoneExists {
		return ErrPhoneAlreadyExists
	}

	if err := uc.AuthRepo.RegisterUser(ctx, reqModel); err != nil {
		return err
	}

	return nil
}
