package auth

import (
	"context"
	"errors"
	"fmt"
	"modulegue/core/hash"

	// "modulegue/internal/data/mobile/remote/dto"
	model "modulegue/internal/domain/mobile/model/auth"
	"modulegue/internal/domain/mobile/repository"
	"time"
)

var (
	ErrInvalidInput          = errors.New("tolong isi field yang wajib")
	ErrEmailAlreadyExists    = errors.New("Email sudah terdaftar")
	ErrPhoneAlreadyExists    = errors.New("Phone sudah terdaftar")
	ErrUsernameAlreadyExists = errors.New("Username sudah terdaftar")
)

type RegisterUseCase struct {
	AuthRepo  repository.AuthRepository
	reqModel  model.RegisterRequestModel
	respModel model.RegisterResponseModel
}

func NewRegisterUseCase(
	AuthRepo repository.AuthRepository,
	reqModel model.RegisterRequestModel,
	respModel model.RegisterResponseModel,
) *RegisterUseCase {
	return &RegisterUseCase{
		AuthRepo:  AuthRepo,
		reqModel:  reqModel,
		respModel: respModel,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, reqModel model.RegisterRequestModel) (*model.RegisterResponseModel, error) {
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

	hashedPassword, err := hash.Hash(reqModel.Password)
	if err != nil {
		return nil, err
	}

	newUser := &model.UserModel{
		RoleId:       2,
		FullName:     reqModel.FullName,
		Nik:          reqModel.NIK,
		Phone:        reqModel.Phone,
		Email:        reqModel.Email,
		Username:     reqModel.Username,
		PasswordHash: hashedPassword,
		IsVerified:   false, // Bisa di-set true jika verifikasi otomatis
		RegisteredAt: time.Now(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	newUser, err = uc.AuthRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan user: %w", err)
	}

	// 6. (Opsional) Buat wallet default untuk user baru
	// Kita anggap wallet akan dibuat oleh usecase lain atau oleh trigger di database.
	// Jika ingin dibuat disini, kamu perlu repository wallet dan logic-nya.
	// Misalnya:
	// err = uc.walletRepo.CreateDefaultWallet(ctx, newUser.ID)
	// if err != nil {
	//     // Log error, tapi jangan hentikan register jika wallet opsional
	//     log.Printf("Gagal buat wallet default untuk user %d: %v", newUser.ID, err)
	// }

	return &model.RegisterResponseModel{
		Message: "registrasi berhasil",
		User: model.UserModel{
			UserId:     newUser.UserId,
			Nik:        newUser.Nik,
			FullName:   newUser.FullName,
			Phone:      newUser.Phone,
			Email:      newUser.Email,
			Username:   newUser.Username,
			Password:   "",
			RoleId:     newUser.RoleId,
			IsVerified: newUser.IsVerified,
			Lokasi:     newUser.Lokasi,
			Zona:       newUser.Zona,
			Tarif:      nil,
		},
	}, nil
}
