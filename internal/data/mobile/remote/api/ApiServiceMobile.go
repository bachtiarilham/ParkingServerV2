package api

import (
	"log"
	"net/http"
	"time"

	"modulegue/config"
	mobile_handler "modulegue/internal/delivery/mobile/customer/handler"
	jukir_handler "modulegue/internal/delivery/mobile/handler"
	shared_handler "modulegue/internal/delivery/shared/handler"

	middleware "modulegue/internal/middleware"
	authuc "modulegue/internal/usecase/auth"
	homeuc "modulegue/internal/usecase/home"
	paymentuc "modulegue/internal/usecase/payment"
	useruc "modulegue/internal/usecase/user"

	qruc "modulegue/internal/usecase/qr"

	"modulegue/core/queue"
)

func RegisterRoutes(
	mux *http.ServeMux,
	//auth
	registerUC *authuc.RegisterUseCase,
	loginUC *authuc.LoginUseCase,
	logoutUC *authuc.LogoutUseCase,
	refreshUC *authuc.RefreshTokenUseCase,
	changePasswordUC *authuc.ChangePasswordUseCase,
	//setting
	getUserProfileUC *useruc.GetProfileUseCase,
	//home
	homeUC *homeuc.GetDashboardUseCase,
	//payment
	qrUC *qruc.GenerateQRUseCase,
	initiatePaymentUC *paymentuc.InitiatePaymentUseCase,
	//pkg
	q *queue.Queue,
	logger *log.Logger,
) {
	// Buat handler
	authHandler := shared_handler.NewAuthHandler(registerUC, loginUC, refreshUC, logoutUC, changePasswordUC) // Pastikan constructor AuthHandler sudah menerima semua usecase yang diperlukan
	homeHandler := mobile_handler.NewHomeHandler(homeUC)
	// scanHandler := mobile_handler.NewScanHandler(submitQrUC, scanDetailUC)
	userHandler := shared_handler.NewUserHandler(getUserProfileUC)
	paymentHandler := jukir_handler.NewPaymentHandler(initiatePaymentUC)
	qrHandler := mobile_handler.NewQRHandler(qrUC)
	authLimiter := middleware.NewRateLimiter(10, time.Minute, 5)

	// Daftarkan endpoint sesuai Android
	// route tanpa otentikasi
	mux.HandleFunc("POST /api/v2/linespot/auth/register", authHandler.Register)
	mux.Handle("POST /api/v2/linespot/auth/login", authLimiter.AllowLogin(http.HandlerFunc(authHandler.Login)))
	mux.Handle("POST /api/v2/linespot/auth/refresh", authLimiter.AllowRefresh(http.HandlerFunc(authHandler.Refresh)))

	//inisiasi middleware JWT untuk route yang membutuhkan otentikasi
	protectedHomeHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(homeHandler.GetDashboard))
	protectedProfileHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(userHandler.GetCurrentUser))
	protectedLogoutHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(authHandler.Logout))                 // jika logout juga butuh otentikasi
	protectedChangePasswordHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(authHandler.ChangePassword)) // Ganti cfg.JWTSecret dengan cara kamu mengakses secret
	protectedQrGenerator := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(qrHandler.GenerateQR))
	protectedPaymentHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(paymentHandler.InitiatePayment))

	//route dengan otentikasi (middleware JWT)
	mux.Handle("GET /api/v2/linespot/home", protectedHomeHandler)
	mux.Handle("GET /api/v2/linespot/users/me", protectedProfileHandler)
	mux.Handle("POST /api/v2/linespot/auth/logout", protectedLogoutHandler)
	mux.Handle("POST /api/v2/linespot/auth/change-password", protectedChangePasswordHandler)
	mux.Handle("GET /api/v2/linespot/qr/generate", protectedQrGenerator)
	mux.Handle("POST /api/v2/linespot/pay", protectedPaymentHandler)
}
