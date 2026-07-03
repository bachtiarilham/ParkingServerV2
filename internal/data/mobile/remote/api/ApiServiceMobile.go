package api

import (
	"log"
	"net/http"
	"time"

	"modulegue/config"

	"modulegue/core/queue"
	middleware "modulegue/internal/middleware"

	authUc "modulegue/internal/domain/mobile/usecase/auth"
	authHandler "modulegue/internal/handler/mobile/auth"

	homeUc "modulegue/internal/domain/mobile/usecase/home"
	homeHandler "modulegue/internal/handler/mobile/home"
	// mobile_handler "modulegue/internal/delivery/mobile/customer/handler"
	// jukir_handler "modulegue/internal/delivery/mobile/handler"
	// shared_handler "modulegue/internal/delivery/shared/handler"
	// paymentuc "modulegue/internal/usecase/payment"
	// useruc "modulegue/internal/usecase/user"
	// qruc "modulegue/internal/usecase/qr"
)

func RegisterRoutes(
	mux *http.ServeMux,
	//auth
	loginUc *authUc.LoginUseCase,
	registerUc *authUc.RegisterUseCase,
	logoutUc *authUc.LogoutUseCase,
	refreshUc *authUc.RefreshTokenUseCase,
	changePasswordUc *authUc.ChangePasswordUseCase,
	//setting
	// getUserProfileUC *useruc.GetProfileUseCase,
	//home
	homeUc *homeUc.GetHomeUseCase,
	//payment
	// qrUC *qruc.GenerateQRUseCase,
	// initiatePaymentUC *paymentuc.InitiatePaymentUseCase,
	//pkg
	q *queue.Queue,
	logger *log.Logger,
) {
	//core
	authLimiter := middleware.NewRateLimiter(10, time.Minute, 5)

	// Handler
	//auth
	loginHandler := authHandler.NewLoginHandler(loginUc)
	registerHandler := authHandler.NewRegisterHandler(registerUc)
	logoutHandler := authHandler.NewLogoutHandler(logoutUc)
	refreshHandler := authHandler.NewRefreshTokenHandler(refreshUc)
	homeHandler := homeHandler.NewHomeHandler(homeUc)
	// scanHandler := mobile_handler.NewScanHandler(submitQrUC, scanDetailUC)
	// userHandler := shared_handler.NewUserHandler(getUserProfileUC)
	// paymentHandler := jukir_handler.NewPaymentHandler(initiatePaymentUC)
	// qrHandler := mobile_handler.NewQRHandler(qrUC)
	// authHandler := shared_handler.NewAuthHandler(registerUC, loginUC, refreshUC, logoutUC, changePasswordUC) // Pastikan constructor AuthHandler sudah menerima semua usecase yang diperlukan

	// route tanpa otentikasi
	mux.HandleFunc("POST /api/v2/linespot/auth/register", registerHandler.Register)
	mux.Handle("POST /api/v2/linespot/auth/login", authLimiter.AllowLogin(http.HandlerFunc(loginHandler.Login)))
	mux.Handle("POST /api/v2/linespot/auth/refreshToken", authLimiter.AllowRefresh(http.HandlerFunc(refreshHandler.Execute)))

	//inisiasi middleware JWT untuk route yang membutuhkan otentikasi
	protectedHomeHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(homeHandler.GetDashboard))
	// protectedProfileHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(userHandler.GetCurrentUser))
	protectedLogoutHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(logoutHandler.Logout)) // jika logout juga butuh otentikasi
	// protectedChangePasswordHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(authHandler.ChangePassword)) // Ganti cfg.JWTSecret dengan cara kamu mengakses secret
	// protectedQrGenerator := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(qrHandler.GenerateQR))
	// protectedPaymentHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(paymentHandler.InitiatePayment))

	//route dengan otentikasi (middleware JWT)
	mux.Handle("GET /api/v2/linespot/home", protectedHomeHandler)
	// mux.Handle("GET /api/v2/linespot/users/me", protectedProfileHandler)
	mux.Handle("POST /api/v2/linespot/auth/logout", protectedLogoutHandler)
	// mux.Handle("POST /api/v2/linespot/auth/change-password", protectedChangePasswordHandler)
	// mux.Handle("GET /api/v2/linespot/qr/generate", protectedQrGenerator)
	// mux.Handle("POST /api/v2/linespot/pay", protectedPaymentHandler)
}
