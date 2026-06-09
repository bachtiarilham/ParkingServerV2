package auth

import (
	"context"
	"errors"
	"fmt"
	domain_auth "modulegue/internal/domain/auth" // Alias untuk menghindari konflik
	"modulegue/internal/domain/user"
	"modulegue/pkg/hash" // Gunakan pkg/hash kamu
	"time"
)

var (
	ErrEmailAlreadyExists    = errors.New("email sudah terdaftar")
	ErrUsernameAlreadyExists = errors.New("username sudah digunakan")
)

type RegisterInput struct {
	// Id	 int64 // Biasanya ID akan di-generate oleh database, jadi ini bisa diabaikan saat input
	FullName string
	Nik      string
	Phone    string
	Email    string
	Username string
	Password string
	// Tidak termasuk Role, karena biasanya default ke 'customer' atau ditentukan oleh admin
}

type RegisterOutput struct {
	UserID  int64
	Message string
}

type RegisterUseCase struct {
	userRepo      user.Repository
	authRepo      domain_auth.Repository // Jika nanti perlu buat session otomatis setelah register
	defaultRoleID int64                  // Role default untuk user baru, misalnya customer
}

func NewRegisterUseCase(
	userRepo user.Repository,
	authRepo domain_auth.Repository,
	defaultRoleID int64, // Contoh: 3 untuk customer
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:      userRepo,
		authRepo:      authRepo, // Bisa jadi tidak digunakan langsung di Register, tapi kita siapkan jika diperlukan
		defaultRoleID: defaultRoleID,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterInput) (RegisterOutput, error) {
	// 1. Cek apakah email sudah terdaftar
	existingUserByEmail, _ := uc.userRepo.FindByEmail(ctx, input.Email)
	if existingUserByEmail != nil {
		return RegisterOutput{}, ErrEmailAlreadyExists
	}

	// 2. Cek apakah username sudah digunakan
	existingUserByUsername, _ := uc.userRepo.FindByUsername(ctx, input.Username) // Kita perlu method ini di UserRepository
	if existingUserByUsername != nil {
		return RegisterOutput{}, ErrUsernameAlreadyExists
	}

	// 3. Hash password
	hashedPassword, err := hash.Hash(input.Password)
	if err != nil {
		return RegisterOutput{}, fmt.Errorf("gagal hash password: %w", err)
	}

	// 4. Buat entity user baru
	newUser := &user.User{
		RoleID:           uc.defaultRoleID,
		FullName:         input.FullName,
		Nik:              input.Nik,
		Phone:            input.Phone,
		Email:            input.Email,
		Username:         input.Username,
		PasswordHash:     hashedPassword,
		EmploymentStatus: "active", // Default
		IsVerified:       false,    // Bisa di-set true jika verifikasi otomatis
		RegisteredAt:     time.Now(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// 5. Simpan user ke database
	err = uc.userRepo.Create(ctx, newUser)
	if err != nil {
		return RegisterOutput{}, fmt.Errorf("gagal menyimpan user: %w", err)
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

	return RegisterOutput{
		UserID:  newUser.ID,
		Message: "registrasi berhasil",
	}, nil
}
