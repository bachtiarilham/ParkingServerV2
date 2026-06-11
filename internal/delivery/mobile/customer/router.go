//mount semua route mobile

package customer

import (
	// "database/sql"
	"log"
	"net/http"

	"modulegue/config"
	"modulegue/internal/delivery/mobile/customer/handler"
	middleware "modulegue/internal/middleware"
	authuc "modulegue/internal/usecase/auth"
	homeuc "modulegue/internal/usecase/home"
	payuc "modulegue/internal/usecase/payment"
	"modulegue/pkg/queue"
)

// func RegisterRoutes(mux *http.ServeMux, db *sql.DB, q *queue.Queue, logger *log.Logger) {
// 	authHandler := handler.NewAuthHandler(db, q) // buat handler
// 	// Daftarkan endpoint
// 	mux.HandleFunc("POST /auth/register", authHandler.Register)
// 	// endpoint mobile lainnya...
// }

func RegisterRoutes(
	mux *http.ServeMux,
	registerUC *authuc.RegisterUseCase,
	loginUC *authuc.LoginUseCase,
	refreshUC *authuc.RefreshTokenUseCase,
	homeUC *homeuc.GetDashboardUseCase,
	executePaymentUC *payuc.ExecutePaymentUseCase,
	scanDetailUC *payuc.GetScanDetailUseCase,
	submitQrUC *payuc.SubmitQrUseCase,
	q *queue.Queue,
	logger *log.Logger,
) {
	// Buat handler
	authHandler := handler.NewAuthHandler(registerUC, loginUC, refreshUC)
	homeHandler := handler.NewHomeHandler(homeUC)
	scanHandler := handler.NewScanHandler(submitQrUC, scanDetailUC)
	paymentHandler := handler.NewPaymentHandler(executePaymentUC)

	// Daftarkan endpoint sesuai Android
	// route tanpa otentikasi
	mux.HandleFunc("POST /api/v2/linespot/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v2/linespot/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v2/linespot/auth/refresh", authHandler.Refresh)
	// route dengan otentikasi (middleware JWT)
	protectedHomeHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(homeHandler.GetDashboard))
	mux.Handle("GET /api/v2/linespot/home", protectedHomeHandler)           // ← GET, bukan POST
	mux.HandleFunc("POST /api/v2/linespot/scan", scanHandler.SubmitAndScan) // ← satu endpoint untuk submit + scan
	mux.HandleFunc("POST /api/v2/linespot/pay", paymentHandler.ExecutePayment)
	// mux.HandleFunc("GET /api/v2/linespot/users/me", userHandler.GetCurrentUser) // tambahkan jika ada
}

// func RegisterRoutes(
// 	mux *http.ServeMux,
// 	registerUC *authuc.RegisterUseCase,
// 	loginUC *authuc.LoginUseCase,
// 	homeUC *homeuc.GetDashboardUseCase,
// 	excpayUC *payuc.ExecutePaymentUseCase,
// 	scanDetailUC *payuc.GetScanDetailUseCase,
// 	submitQrUC *payuc.SubmitQrUseCase,
// 	q *queue.Queue,
// 	logger *log.Logger,
// ) {
// 	// authHandler := handler.NewAuthHandler(registerUC, loginUC)
// 	// mux.HandleFunc("POST /auth/register", authHandler.Register)
// 	// mux.HandleFunc("POST /auth/login", authHandler.Login)
// 	// // ... lainnya

// 	//register handler
// 	authHandler := handler.NewAuthHandler(registerUC, loginUC)
// 	homeHandler := handler.NewHomeHandler(homeUC)
// 	payHandler := handler.NewPayHandler(payUC)

// 	//register endpoint
// 	mux.HandleFunc("POST /api/v2/linespot/auth/register", authHandler.Register)
// 	mux.HandleFunc("POST /api/v2/linespot/auth/login", authHandler.Login)
// 	// mux.HandleFunc("POST /auth/logout", authHandler.Logout)
// 	// mux.HandleFunc("POST /auth/refresh", authHandler.Refresh)
// 	// mux.HandleFunc("POST /auth/forgot-password", authHandler.ForgotPassword)
// 	// mux.HandleFunc("POST /auth/reset-password", authHandler.ResetPassword)
// 	// mux.HandleFunc("POST /auth/verify-email", authHandler.VerifyEmail)
// 	// mux.HandleFunc("POST /auth/resend-verification", authHandler.ResendVerification)
// 	mux.HandleFunc("POST /api/v2/linespot/home", homeHandler.GetDashboard)
// 	mux.HandleFunc("POST /api/v2/linespot/scan", authHandler.Scan)
// 	mux.HandleFunc("POST /api/v2/linespot/pay", payHandler.Pay)
// }
