package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	model "modulegue/internal/domain/mobile/model/auth"
	"modulegue/internal/domain/mobile/repository"
)

var (
	ErrInvalidInput          = errors.New("tolong isi field yang wajib")
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

func (uc *RegisterUseCase) Execute(ctx context.Context, reqModel model.RegisterRequestModel) (*model.RegisterResponseModel, error) {
	reqModel.FullName = strings.TrimSpace(reqModel.FullName)
	reqModel.NIK = strings.TrimSpace(reqModel.NIK)
	reqModel.Phone = strings.TrimSpace(reqModel.Phone)
	reqModel.Email = strings.ToLower(strings.TrimSpace(reqModel.Email))
	reqModel.Username = strings.TrimSpace(reqModel.Username)
	reqModel.Password = strings.TrimSpace(reqModel.Password)

	if reqModel.Username == "" || reqModel.Email == "" || reqModel.Password == "" {
		return nil, ErrInvalidInput
	}

	exists, err := uc.AuthRepo.ExistsByEmailOrUsernameOrPhone(ctx, reqModel.Email, reqModel.Username, reqModel.Phone)
	if err != nil {
		return nil, err
	}
	if exists.EmailExists {
		return nil, ErrEmailAlreadyExists
	}
	if exists.UsernameExists {
		return nil, ErrUsernameAlreadyExists
	}
	if exists.PhoneExists {
		return nil, ErrPhoneAlreadyExists
	}

	result, err := uc.AuthRepo.RegisterUser(ctx, reqModel)
	if err != nil {
		return nil, fmt.Errorf("register user: %w", err)
	}

	if result == nil {
		return &model.RegisterResponseModel{
			Message: "registrasi berhasil",
		}, nil
	}

	return result, nil
}
