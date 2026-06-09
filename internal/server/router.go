// //mount ketiga delivery router

// package server

// import (
// 	"net/http"

// 	"github.com/bachtiarilham/ParkServerV2/internal/delivery/mobile/handler"
// )

// func NewRouter(
// 	authHandler *handler.AuthHandler,
// ) http.Handler {

// 	mux := http.NewServeMux()

// 	mux.HandleFunc(
// 		"POST /auth/register",
// 		authHandler.Register,
// 	)

// 	return mux
// }

package server

import (
	"database/sql"
	"log"
	"net/http"

	"modulegue/config"
	// "modulegue/internal/delivery/desktop"
	"modulegue/internal/delivery/mobile/customer"
	"modulegue/internal/delivery/mobile/jukir"

	// "modulegue/internal/delivery/web"
	"modulegue/internal/repository"
	authuc "modulegue/internal/usecase/auth"
	"modulegue/pkg/queue"

	// "modulegue/pkg/response"
	"modulegue/pkg/worker"
)

// func NewRouter(
// 	cfg *config.Config,
// 	db *sql.DB,
// 	queue *queue.Queue,
// 	logger *log.Logger,
// ) (http.Handler, *worker.WorkerPool) {
// 	mux := http.NewServeMux()

// 	workerPool := worker.NewWorkerPool(queue, cfg.QueueWorkerCount, logger) // sesuaikan constructor

// 	// Daftarkan route dari ketiga delivery
// 	// customer.RegisterRoutes(mux, db, queue, logger)
// 	// jukir.RegisterRoutes(mux, db, queue, logger)
// 	customer.RegisterRoutes(mux, registerUC, loginUC, queue, logger) // contoh
// 	jukir.RegisterRoutes(mux, registerUC, loginUC, queue, logger)
// 	desktop.RegisterRoutes(mux, db, queue, logger)
// 	web.RegisterRoutes(mux, db, queue, logger)

// 	// Siapkan worker pool (contoh: baca konfigurasi dari cfg)

// 	return mux, workerPool
// }

func NewRouter(
	cfg *config.Config,
	db *sql.DB,
	queue *queue.Queue,
	logger *log.Logger,
) (http.Handler, *worker.WorkerPool) {
	mux := http.NewServeMux()

	workerPool := worker.NewWorkerPool(queue, cfg.QueueWorkerCount, logger)

	// Buat semua repositories
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)
	// tambahkan repositori lain sesuai kebutuhan
	// misalnya: parkingRepo, transactionRepo, dll

	// Buat semua usecases
	jwtSecret := cfg.JWTSecret
	accessTTL := cfg.AccessTokenMinutes
	refreshTTL := cfg.RefreshTokenHours

	registerUC := authuc.NewRegisterUseCase(userRepo, authRepo, 3) // Contoh: 3 untuk role customer
	loginUC := authuc.NewLoginUseCase(userRepo, authRepo, jwtSecret, accessTTL, refreshTTL)

	// Buat usecase lainnya sesuai kebutuhan
	// misalnya: parkingUC, transactionUC, dll

	// Daftarkan routes dengan usecase
	customer.RegisterRoutes(mux, registerUC, loginUC, queue, logger)
	jukir.RegisterRoutes(mux, registerUC, loginUC, queue, logger)
	// desktop.RegisterRoutes(mux, registerUC, loginUC, queue, logger)
	// web.RegisterRoutes(mux, registerUC, loginUC, queue, logger)

	return mux, workerPool
}
