//mount semua route mobile

package jukir

import (
	"log"
	"net/http"

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
	// authHandler := handler.NewAuthHandler(registerUC, loginUC)
	// mux.HandleFunc("POST /api/v2/linespot/auth/register", authHandler.Register)
	// mux.HandleFunc("POST /api/v2/linespot/auth/login", authHandler.Login)
	// mux.HandleFunc("POST /auth/logout", authHandler.Logout)
	// mux.HandleFunc("POST /auth/refresh", authHandler.Refresh)
	// mux.HandleFunc("POST /auth/forgot-password", authHandler.ForgotPassword)
	// mux.HandleFunc("POST /auth/reset-password", authHandler.ResetPassword)
	// mux.HandleFunc("POST /auth/verify-email", authHandler.VerifyEmail)
	// mux.HandleFunc("POST /auth/resend-verification", authHandler.ResendVerification)
	// mux.HandleFunc("POST /api/v2/linespot/home", authHandler.Login)
	// mux.HandleFunc("POST /api/v2/linespot/scan", authHandler.Scan)
	// mux.HandleFunc("POST /api/v2/linespot/pay", authHandler.Pay)
	// ... lainnya
}
