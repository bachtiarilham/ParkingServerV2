// // //mount ketiga delivery router

// // package server

// // import (
// // 	"net/http"

// // 	"github.com/bachtiarilham/ParkServerV2/internal/delivery/mobile/handler"
// // )

// // func NewRouter(
// // 	authHandler *handler.AuthHandler,
// // ) http.Handler {

// // 	mux := http.NewServeMux()

// // 	mux.HandleFunc(
// // 		"POST /auth/register",
// // 		authHandler.Register,
// // 	)

// // 	return mux
// // }

// package server

// import (
// 	"database/sql"
// 	"log"
// 	"net/http"

// 	"modulegue/config"
// 	// "modulegue/internal/delivery/desktop"
// 	"modulegue/internal/delivery/mobile/customer"
// 	"modulegue/internal/delivery/mobile/jukir"

// 	// "modulegue/internal/delivery/web"
// 	"modulegue/internal/repository"
// 	authuc "modulegue/internal/usecase/auth"
// 	homeuc "modulegue/internal/usecase/home"
// 	payuc "modulegue/internal/usecase/payment"
// 	"modulegue/pkg/queue"

// 	// "modulegue/pkg/response"
// 	"modulegue/pkg/worker"
// )

// // func NewRouter(
// // 	cfg *config.Config,
// // 	db *sql.DB,
// // 	queue *queue.Queue,
// // 	logger *log.Logger,
// // ) (http.Handler, *worker.WorkerPool) {
// // 	mux := http.NewServeMux()

// // 	workerPool := worker.NewWorkerPool(queue, cfg.QueueWorkerCount, logger) // sesuaikan constructor

// // 	// Daftarkan route dari ketiga delivery
// // 	// customer.RegisterRoutes(mux, db, queue, logger)
// // 	// jukir.RegisterRoutes(mux, db, queue, logger)
// // 	customer.RegisterRoutes(mux, registerUC, loginUC, queue, logger) // contoh
// // 	jukir.RegisterRoutes(mux, registerUC, loginUC, queue, logger)
// // 	desktop.RegisterRoutes(mux, db, queue, logger)
// // 	web.RegisterRoutes(mux, db, queue, logger)

// // 	// Siapkan worker pool (contoh: baca konfigurasi dari cfg)

// // 	return mux, workerPool
// // }

// func NewRouter(
// 	cfg *config.Config,
// 	db *sql.DB,
// 	queue *queue.Queue,
// 	logger *log.Logger,
// ) (http.Handler, *worker.WorkerPool) {
// 	mux := http.NewServeMux()

// 	workerPool := worker.NewWorkerPool(queue, cfg.QueueWorkerCount, logger)

// 	// Buat semua repositories
// 	userRepo := repository.NewUserRepository(db)
// 	// authRepo := repository.NewAuthRepository(db)
// 	// parkingRepo := repository.NewParkingRepository(db) // contoh tambahan
// 	// transactionRepo := repository.NewTransactionRepository(db) // contoh tambahan
// 	// Buat repositori lainnya sesuai kebutuhan
// 	// tambahkan repositori lain sesuai kebutuhan
// 	// misalnya: parkingRepo, transactionRepo, dll

// 	// Buat semua usecases
// 	jwtSecret := cfg.JWTSecret
// 	accessTTL := cfg.AccessTokenMinutes
// 	refreshTTL := cfg.RefreshTokenHours

// 	registerUC := authuc.NewRegisterUseCase(userRepo, authRepo, 3) // Contoh: 3 untuk role customer
// 	loginUC := authuc.NewLoginUseCase(userRepo, authRepo, jwtSecret, accessTTL, refreshTTL)
// 	homeUC := homeuc.NewGetDashboardUseCase(userRepo /*, parkingRepo, transactionRepo, etc */)
// 	executePaymentUC := payuc.NewExecutePaymentUseCase(/* repositori yang dibutuhkan */)
// 	scanDetailUC := payuc.NewGetScanDetailUseCase(/* repositori yang dibutuhkan */)
// 	submitQrUC := payuc.NewSubmitQrUseCase(/* repositori yang dibutuhkan */)

// 	// Buat usecase lainnya sesuai kebutuhan
// 	// misalnya: parkingUC, transactionUC, dll

// 	// Daftarkan routes dengan usecase
// 	customer.RegisterRoutes(mux, registerUC, loginUC, homeUC, executePaymentUC, scanDetailUC, submitQrUC, queue, logger)
// 	jukir.RegisterRoutes(mux, registerUC, loginUC, queue, logger)
// 	// desktop.RegisterRoutes(mux, registerUC, loginUC, queue, logger)
// 	// web.RegisterRoutes(mux, registerUC, loginUC, queue, logger)

// 	return mux, workerPool
// }

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
	authuc "modulegue/internal/usecase/auth"
	homeuc "modulegue/internal/usecase/home"
	transuc "modulegue/internal/usecase/payment"
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
	authRepo := repository.NewAuthRepository(db)               // <-- Uncomment dan buat jika belum ada
	transactionRepo := repository.NewTransactionRepository(db) // <-- Buat jika belum ada
	homeRepo := repository.NewHomeRepository(db)               // <-- Buat jika home butuh repo khusus
	// locationRepo := repository.NewLocationRepository(db)       // <-- Buat jika belum ada
	// vehicleTypeRepo := repository.NewVehicleTypeRepository(db) // <-- Buat jika belum ada
	// Tambahkan repositori lain jika diperlukan (misalnya walletRepo)

	// --- Buat semua Use Cases ---
	jwtSecret := cfg.JWTSecret
	accessTTL := cfg.AccessTokenMinutes // Asumsi ini dalam durasi time.Minute
	refreshTTL := cfg.RefreshTokenHours // Asumsi ini dalam durasi time.Hour

	// Auth Use Cases
	registerUC := authuc.NewRegisterUseCase(
		userRepo,
		authRepo, // Jika register perlu buat session otomatis
		3,        // Contoh: 3 untuk role customer (harus sesuai dengan id di tabel system_role)
	)
	loginUC := authuc.NewLoginUseCase(
		userRepo,
		authRepo, // Untuk menyimpan session
		jwtSecret,
		accessTTL,
		refreshTTL,
	)

	getUserProfileUC := useruc.NewGetProfileUseCase(userRepo)
	refreshUC := authuc.NewRefreshTokenUseCase(
		authRepo,
		userRepo, // Untuk verifikasi user aktif saat refresh
		jwtSecret,
		accessTTL,
		refreshTTL,
	)

	// Home Use Case
	homeUC := homeuc.NewGetDashboardUseCase(
		homeRepo, // Buat homeRepo jika diperlukan, atau gunakan kombinasi repo di bawah
		// userRepo,           // Untuk profil dan summary (wallet)
		// transactionRepo,    // Untuk warnings (finance)
		// locationRepo, // Jika home menampilkan info lokasi
		// Tambahkan repo lain sesuai kebutuhan usecase
	)

	// Transaction Use Cases
	submitQrUC := transuc.NewSubmitQrUseCase(transactionRepo) // <-- Dari 'transaction' package
	scanDetailUC := transuc.NewGetScanDetailUseCase(
		transactionRepo,
		// locationRepo,     // Untuk nama lokasi
		// vehicleTypeRepo,  // Untuk info kendaraan (jika diperlukan lebih lanjut)
		// userRepo, // Jika diperlukan untuk nama customer
	)
	executePaymentUC := transuc.NewExecutePaymentUseCase(
		transactionRepo,
		// paymentGatewayRepo, // Repo untuk interaksi gateway (jika ada)
		// walletRepo, // Jika pembayaran bisa dari wallet
	)

	// --- Daftarkan Routes ---
	// Customer Routes
	// var refreshUC *authuc.RefreshTokenUseCase
	logoutUC := authuc.NewLogoutUseCase(authRepo)

	changePasswordUC := authuc.NewChangePasswordUseCase(userRepo, authRepo) // <-- Tambahkan authRepo sebagai parameter

	customer.RegisterRoutes(
		mux,
		registerUC,
		loginUC,
		logoutUC, // <-- Tambahkan LogoutUseCase
		refreshUC,
		getUserProfileUC,
		changePasswordUC, // <-- Tambahkan ChangePasswordUseCase
		homeUC,
		executePaymentUC, // <-- Gunakan dari 'transaction'
		scanDetailUC,     // <-- Gunakan dari 'transaction'
		submitQrUC,       // <-- Gunakan dari 'transaction'
		queue,
		logger,
	)

	// Jukir Routes (contoh, sesuaikan dengan kebutuhan jukir)
	jukir.RegisterRoutes(
		mux,
		registerUC, // Jika jukir bisa register
		loginUC,    // Jika jukir bisa login
		// Tambahkan UC lain yang spesifik untuk jukir
		queue,
		logger,
	)

	// Tambahkan routes untuk desktop, web jika ada
	// desktop.RegisterRoutes(...)
	// web.RegisterRoutes(...)

	return mux, workerPool
}
