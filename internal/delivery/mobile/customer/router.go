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
	useruc "modulegue/internal/usecase/user"
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
	logoutUC *authuc.LogoutUseCase,
	refreshUC *authuc.RefreshTokenUseCase,
	getUserProfileUC *useruc.GetProfileUseCase,
	changePasswordUC *authuc.ChangePasswordUseCase,
	homeUC *homeuc.GetDashboardUseCase,
	executePaymentUC *payuc.ExecutePaymentUseCase,
	scanDetailUC *payuc.GetScanDetailUseCase,
	submitQrUC *payuc.SubmitQrUseCase,
	q *queue.Queue,
	logger *log.Logger,
) {
	// Buat handler
	authHandler := handler.NewAuthHandler(registerUC, loginUC, refreshUC, logoutUC, changePasswordUC) // Pastikan constructor AuthHandler sudah menerima semua usecase yang diperlukan
	homeHandler := handler.NewHomeHandler(homeUC)
	scanHandler := handler.NewScanHandler(submitQrUC, scanDetailUC)
	paymentHandler := handler.NewPaymentHandler(executePaymentUC)
	userHandler := handler.NewUserHandler(getUserProfileUC)

	// Daftarkan endpoint sesuai Android
	// route tanpa otentikasi
	mux.HandleFunc("POST /api/v2/linespot/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v2/linespot/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v2/linespot/auth/refresh", authHandler.Refresh)

	//inisiasi middleware JWT untuk route yang membutuhkan otentikasi
	protectedHomeHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(homeHandler.GetDashboard))
	protectedProfileHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(userHandler.GetCurrentUser))
	protectedLogoutHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(authHandler.Logout))                 // jika logout juga butuh otentikasi
	protectedChangePasswordHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(authHandler.ChangePassword)) // Ganti cfg.JWTSecret dengan cara kamu mengakses secret

	//route dengan otentikasi (middleware JWT)
	mux.Handle("GET /api/v2/linespot/home", protectedHomeHandler)           // ← GET, bukan POST
	mux.HandleFunc("POST /api/v2/linespot/scan", scanHandler.SubmitAndScan) // ← satu endpoint untuk submit + scan
	mux.HandleFunc("POST /api/v2/linespot/pay", paymentHandler.ExecutePayment)
	mux.Handle("GET /api/v2/linespot/users/me", protectedProfileHandler)
	mux.Handle("POST /api/v2/linespot/auth/logout", protectedLogoutHandler)
	mux.Handle("POST /api/v2/linespot/auth/change-password", protectedChangePasswordHandler) // Gunakan path yang kamu inginkan
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
