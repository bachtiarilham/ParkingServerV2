// // //mount ketiga delivery router

// server/router.go
package server

import (
	"database/sql"
	"log"
	"net/http"

	"modulegue/config"
	"modulegue/internal/delivery/mobile/customer"
	"modulegue/internal/delivery/mobile/jukir"
	"modulegue/internal/repository"
	pay_svc "modulegue/internal/service/payment_gateway"
	authuc "modulegue/internal/usecase/auth"
	homecustomeruc "modulegue/internal/usecase/home_customer"
	homejukiruc "modulegue/internal/usecase/home_jukir"
	payuc "modulegue/internal/usecase/payment"
	qruc "modulegue/internal/usecase/qr"
	useruc "modulegue/internal/usecase/user"
	"modulegue/pkg/queue"
	"modulegue/pkg/worker"
)

func NewRouter(
	cfg *config.Config,
	db *sql.DB,
	queue *queue.Queue,
	logger *log.Logger,
) (http.Handler, *worker.WorkerPool) {
	mux := http.NewServeMux()

	workerPool := worker.NewWorkerPool(queue, cfg.QueueWorkerCount, logger)

	// --- Buat semua Repositories ---
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db) // <-- Uncomment dan buat jika belum ada
	homeCustomerRepo := repository.NewHomeRepository(db)
	jukirCustomerRepo := repository.NewHomeJukirRepository(db) // <-- Buat jika home butuh repo khusus
	paymentRepo := repository.NewPaymentRepository(db)         // <-- Buat jika belum ada
	jukirWalletRepo := repository.NewWalletRepository(db)
	membershipRepo := repository.NewMembershipRepository(db)

	// locationRepo := repository.NewLocationRepository(db)       // <-- Buat jika belum ada
	// vehicleTypeRepo := repository.NewVehicleTypeRepository(db) // <-- Buat jika belum ada
	// Tambahkan repositori lain jika diperlukan (misalnya walletRepo)
	// paymentGatewaySvc, err := pay_svc.NewMidtransService(pay_svc.MidtransConfig{
	// 	ServerKey:   cfg.MidtransServerKey,
	// 	ClientKey:   cfg.MidtransClientKey,
	// 	Environment: midtrans.Sandbox, // Ganti ke Production saat deploy
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to initialize payment gateway: %v", err)
	// }
	paymentGatewaySvc := pay_svc.NewMockMidtransService()

	jwtSecret := cfg.JWTSecret
	accessTTL := cfg.AccessTokenMinutes // Asumsi ini dalam durasi time.Minute
	refreshTTL := cfg.RefreshTokenHours // Asumsi ini dalam durasi time.Hour
	initiatePaymentUC := payuc.NewInitiatePaymentUseCase(
		paymentRepo, // Gunakan repo payment yang baru
		paymentGatewaySvc,
		jukirWalletRepo,
		membershipRepo,
		0.20, // Gov %
		0.40, // Company %
		0.40, // Jukir %
	)

	//Use Cases
	//auth
	registerUC := authuc.NewRegisterUseCase(userRepo, authRepo, 3)
	loginUC := authuc.NewLoginUseCase(userRepo, authRepo, jwtSecret, accessTTL, refreshTTL)
	refreshUC := authuc.NewRefreshTokenUseCase(authRepo, userRepo, jwtSecret, accessTTL, refreshTTL)
	logoutUC := authuc.NewLogoutUseCase(authRepo)
	changePasswordUC := authuc.NewChangePasswordUseCase(userRepo, authRepo) // <-- Tambahkan authRepo sebagai parameter
	//profile
	getUserProfileUC := useruc.NewGetProfileUseCase(userRepo)
	//home
	homeCustomerUC := homecustomeruc.NewGetDashboardUseCase(homeCustomerRepo)
	homeJukirUC := homejukiruc.NewGetDashboardUseCase(jukirCustomerRepo)
	//payment
	qrUC := qruc.NewGenerateQRUseCase()
	// --- Daftarkan Routes ---
	// Customer Routes
	// var refreshUC *authuc.RefreshTokenUseCase

	customer.RegisterRoutes(
		mux,
		//auth
		registerUC,
		loginUC,
		logoutUC,
		refreshUC,
		changePasswordUC,
		//profile
		getUserProfileUC,
		//home
		homeCustomerUC,
		//payment
		qrUC,
		//pkg
		queue,
		logger,
	)

	// Jukir Routes (contoh, sesuaikan dengan kebutuhan jukir)
	jukir.RegisterRoutes(
		mux,
		//auth
		registerUC,
		loginUC,
		logoutUC,
		refreshUC,
		changePasswordUC,
		//profile
		getUserProfileUC,
		//home
		homeJukirUC,
		//payment
		initiatePaymentUC,
		//pkg
		queue,
		logger,
	)

	// Tambahkan routes untuk desktop, web jika ada
	// desktop.RegisterRoutes(...)
	// web.RegisterRoutes(...)

	return mux, workerPool
}
