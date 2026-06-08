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
	"modulegue/internal/delivery/desktop"
	"modulegue/internal/delivery/mobile/customer"
	"modulegue/internal/delivery/mobile/jukir"
	"modulegue/internal/delivery/web"
	"modulegue/pkg/queue"

	// "modulegue/pkg/response"
	"modulegue/pkg/worker"
)

func NewRouter(
	cfg *config.Config,
	db *sql.DB,
	queue *queue.Queue,
	logger *log.Logger,
) (http.Handler, *worker.WorkerPool) {
	// _ = cfg
	mux := http.NewServeMux()

	// adminService := adminstore.New(db)
	// adminHTTPHandler := adminendpoints.NewHTTPHandler(adminService, queue, cfg)
	// mobileHTTPHandler := mobile.NewHTTPHandler(db, queue, cfg)

	workerPool := worker.NewWorkerPool(queue, cfg.QueueWorkerCount, logger) // sesuaikan constructor

	// workerPool := worker.NewWorkerPool(queue, cfg.QueueWorkerCount, logger)
	// adminHTTPHandler.RegisterQueueProcessors(workerPool)
	// mobileHTTPHandler.RegisterQueueProcessors(workerPool)

	// healthHandler := func(service string) http.HandlerFunc {
	// 	return func(w http.ResponseWriter, r *http.Request) {

	// 		healthData := map[string]any{
	// 			"success": true,
	// 			"service": service,
	// 			"status":  "UP",
	// 		}

	// 		response.Success(w, http.StatusOK, "Service is healthy", healthData)

	// 		// response.JSON(w, http.StatusOK, response.APIResponse{
	// 		// 	"success": true,
	// 		// 	"service": service,
	// 		// 	"status":  "UP",
	// 		// })
	// 	}
	// }

	// unhealthy := func(service string) http.HandlerFunc {
	// 	return func(w http.ResponseWriter, r *http.Request) {
	// 		// healthData := map[string]any{
	// 		// 	"success": false,
	// 		// 	"service": service,
	// 		// 	"status":  "DOWN",
	// 		// 	}

	// 		response.Error(w, http.StatusBadRequest, "Service is unhealthy")

	// 		// response.JSON(w, http.StatusNotImplemented, response.APIResponse{
	// 		// 	"success": false,
	// 		// 	"service": service,
	// 		// 	"message": "namespace disiapkan, endpoint bisnis belum diimplementasikan pada rombak v2 ini",
	// 		// })
	// 	}
	// }
	// mux := http.NewServeMux()

	// Daftarkan route dari ketiga delivery
	customer.RegisterRoutes(mux, db, queue, logger)
	jukir.RegisterRoutes(mux, db, queue, logger)
	desktop.RegisterRoutes(mux, db, queue, logger)
	web.RegisterRoutes(mux, db, queue, logger)

	// Siapkan worker pool (contoh: baca konfigurasi dari cfg)

	return mux, workerPool
}
