//mount semua route mobile

package jukir

import (
	"log"
	"net/http"

	"modulegue/internal/delivery/mobile/jukir/handler"
	authuc "modulegue/internal/usecase/auth"
	"modulegue/pkg/queue"
)

// func RegisterRoutes(mux *http.ServeMux, db *sql.DB, q *queue.Queue, logger *log.Logger) {
// 	authHandler := handler.NewAuthHandler(db, q, logger) // buat handler
// 	// Daftarkan endpoint
// 	mux.HandleFunc("POST /auth/register", authHandler.Register)
// 	// endpoint mobile lainnya...
// }

func RegisterRoutes(
	mux *http.ServeMux,
	registerUC *authuc.RegisterUseCase,
	loginUC *authuc.LoginUseCase,
	q *queue.Queue,
	logger *log.Logger,
) {
	authHandler := handler.NewAuthHandler(registerUC, loginUC)
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	// ... lainnya
}
