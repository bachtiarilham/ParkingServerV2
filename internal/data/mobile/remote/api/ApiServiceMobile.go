package api

import (
	"log"
	"net/http"
	"time"

	"modulegue/config"

	"modulegue/core/queue"
	middleware "modulegue/internal/middleware"

	authUc "modulegue/internal/domain/mobile/usecase/auth"
	helperUc "modulegue/internal/domain/mobile/usecase/helper"
	homeUc "modulegue/internal/domain/mobile/usecase/home"
	laporanUc "modulegue/internal/domain/mobile/usecase/laporan"
	riwayatUc "modulegue/internal/domain/mobile/usecase/riwayat"
	subscriptionUc "modulegue/internal/domain/mobile/usecase/subscription"
	authHandler "modulegue/internal/handler/mobile/auth"
	helperHandler "modulegue/internal/handler/mobile/helper"
	homeHandler "modulegue/internal/handler/mobile/home"
	laporanHandler "modulegue/internal/handler/mobile/laporan"
	riwayatHandler "modulegue/internal/handler/mobile/riwayat"
	subscriptionHandler "modulegue/internal/handler/mobile/subscription"
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
	//home
	homeUc *homeUc.GetHomeUseCase,
	//laporan
	laporanUc *laporanUc.GetLaporanUseCase,
	//riwayat
	riwayatUc *riwayatUc.GetRiwayatUseCase,
	//subscription
	subscriptionUc *subscriptionUc.SubscriptionUseCase,
	//helper
	GetLokasiUc *helperUc.GetLokasiUseCase,
	//setting
	// getUserProfileUC *useruc.GetProfileUseCase,
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
	//home
	homeHandler := homeHandler.NewHomeHandler(homeUc)
	//riwayat
	riwayatHandler := riwayatHandler.NewRiwayatHandler(riwayatUc)
	//laporan
	laporanHandler := laporanHandler.NewLaporanHandler(laporanUc)
	//subscription
	subscriptionHandler := subscriptionHandler.NewSubscriptionHandler(subscriptionUc)
	//helper
	getLokasiHandler := helperHandler.NewGetLocationHandler(GetLokasiUc)
	// scanHandler := mobile_handler.NewScanHandler(submitQrUC, scanDetailUC)
	// userHandler := shared_handler.NewUserHandler(getUserProfileUC)
	// paymentHandler := jukir_handler.NewPaymentHandler(initiatePaymentUC)
	// qrHandler := mobile_handler.NewQRHandler(qrUC)
	// authHandler := shared_handler.NewAuthHandler(registerUC, loginUC, refreshUC, logoutUC, changePasswordUC) // Pastikan constructor AuthHandler sudah menerima semua usecase yang diperlukan

	//OtentikasiHandler
	//auth
	protectedLogoutHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(logoutHandler.Logout)) // jika logout juga butuh otentikasi
	// protectedChangePasswordHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(authHandler.ChangePassword)) // Ganti cfg.JWTSecret dengan cara kamu mengakses secret
	//home
	protectedHomeHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(homeHandler.GetDashboard))
	//riwayat
	protectedRiwayatHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(riwayatHandler.Execute))
	//laporan
	protectedLaporanHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(laporanHandler.Execute))
	//subscription
	protectedSubscriptionHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(subscriptionHandler.Execute))
	//helper
	protectedGetLokasiHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(getLokasiHandler.Execute))

	// protectedProfileHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(userHandler.GetCurrentUser))
	// protectedQrGenerator := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(qrHandler.GenerateQR))
	// protectedPaymentHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(paymentHandler.InitiatePayment))

	//Endpoint
	//auth
	mux.HandleFunc("POST /api/v2/linespot/auth/register", registerHandler.Register)
	mux.Handle("POST /api/v2/linespot/auth/login", authLimiter.AllowLogin(http.HandlerFunc(loginHandler.Login)))
	mux.Handle("POST /api/v2/linespot/auth/refreshToken", authLimiter.AllowRefresh(http.HandlerFunc(refreshHandler.Execute)))
	mux.Handle("POST /api/v2/linespot/auth/logout", protectedLogoutHandler)
	// mux.Handle("GET /api/v2/linespot/users/me", protectedProfileHandler)
	// mux.Handle("POST /api/v2/linespot/auth/change-password", protectedChangePasswordHandler)

	//home
	mux.Handle("GET /api/v2/linespot/home", protectedHomeHandler)
	//laporan
	mux.Handle("GET /api/v2/linespot/laporan", protectedLaporanHandler)
	//riwayat
	mux.Handle("GET /api/v2/linespot/riwayat", protectedRiwayatHandler)
	//subscription
	mux.Handle("GET /api/v2/linespot/subscribe", protectedSubscriptionHandler)
	//helper
	mux.Handle("GET /api/v2/linespot/get_lokasi", protectedGetLokasiHandler)

	//route dengan otentikasi (middleware JWT)
	// mux.Handle("GET /api/v2/linespot/qr/generate", protectedQrGenerator)
	// mux.Handle("POST /api/v2/linespot/pay", protectedPaymentHandler)
}
