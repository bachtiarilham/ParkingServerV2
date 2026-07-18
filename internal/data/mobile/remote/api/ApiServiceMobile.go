package api

import (
	"log"
	"net/http"

	// "time"

	"modulegue/config"

	"modulegue/core/queue"
	middleware "modulegue/internal/middleware"

	// //auth
	// authUc "modulegue/internal/domain/mobile/usecase/auth"
	// authHandler "modulegue/internal/handler/mobile/auth"

	//home
	homeUc "modulegue/internal/domain/mobile/usecase/home"
	homeHandler "modulegue/internal/handler/mobile/home"

	//laporan
	laporanUc "modulegue/internal/domain/mobile/usecase/laporan"
	laporanHandler "modulegue/internal/handler/mobile/laporan"

	//riwayat
	riwayatUc "modulegue/internal/domain/mobile/usecase/riwayat"
	riwayatHandler "modulegue/internal/handler/mobile/riwayat"

	//subscription
	subscriptionUc "modulegue/internal/domain/mobile/usecase/subscription"
	subscriptionHandler "modulegue/internal/handler/mobile/subscription"

	//payment
	paymentUc "modulegue/internal/domain/mobile/usecase/payment_parking"
	paymentHandler "modulegue/internal/handler/mobile/payment_parking"

	//parking
	parkingUc "modulegue/internal/domain/mobile/usecase/parking"
	parkingHandler "modulegue/internal/handler/mobile/parking"

	//helper
	helperUc "modulegue/internal/domain/mobile/usecase/helper"
	helperHandler "modulegue/internal/handler/mobile/helper"

	//topup
	topupUc "modulegue/internal/domain/mobile/usecase/topup"
	getTopUpStatusHandler "modulegue/internal/handler/mobile/topup"
	topupHandler "modulegue/internal/handler/mobile/topup"
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
	// loginUc *authUc.LoginUseCase,
	// registerUc *authUc.RegisterUseCase,
	// logoutUc *authUc.LogoutUseCase,
	// refreshUc *authUc.RefreshTokenUseCase,
	// changePasswordUc *authUc.ChangePasswordUseCase,
	//home
	homeUc *homeUc.GetHomeUseCase,
	//laporan
	laporanUc *laporanUc.GetLaporanUseCase,
	//riwayat
	riwayatUc *riwayatUc.GetRiwayatUseCase,
	//subscription
	subscriptionUc *subscriptionUc.SubscriptionUseCase,
	//payment
	postParkingUc *parkingUc.PostParkingUseCase,
	postPaymentParkingUc *paymentUc.PostPaymentParkingUseCase,
	getPembayaranStatusUc *paymentUc.GetPembayaranStatusUseCase,
	//topup
	topupUc *topupUc.TopUpUseCase,
	topupCallbackUc *topupUc.TopUpCallbackUseCase,
	getStatusTopUpUc *topupUc.GetTopUpStatusUseCase,
	midtransVerifier topupHandler.MidtransSignatureVerifier,
	//helper
	GetLokasiUc *helperUc.GetLokasiUseCase,
	GetTarifUc *helperUc.GetTarifUseCase,
	GetNominalTopUpUc *helperUc.GetNominalTopUpUseCase,
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
	// authLimiter := middleware.NewRateLimiter(10, time.Minute, 5)

	// Handler
	//auth
	// loginHandler := authHandler.NewLoginHandler(loginUc)
	// registerHandler := authHandler.NewRegisterHandler(registerUc)
	// logoutHandler := authHandler.NewLogoutHandler(logoutUc)
	// refreshHandler := authHandler.NewRefreshTokenHandler(refreshUc)
	// changePasswordHandler := authHandler.NewChangePasswordHandler(changePasswordUc)
	//home
	homeHandler := homeHandler.NewHomeHandler(homeUc)
	//riwayat
	riwayatHandler := riwayatHandler.NewRiwayatHandler(riwayatUc)
	//laporan
	laporanHandler := laporanHandler.NewLaporanHandler(laporanUc)
	//subscription
	subscriptionHandler := subscriptionHandler.NewSubscriptionHandler(subscriptionUc)
	//parking
	postParkingHandler := parkingHandler.NewPostParkingHandler(postParkingUc)
	//payment
	postPaymentParkingHandler := paymentHandler.NewPostPaymentParkingHandler(postPaymentParkingUc)
	getPembayaranStatusHandler := paymentHandler.NewGetPembayaranStatusHandler(getPembayaranStatusUc)
	//topup
	topupCreateHandler := topupHandler.NewSubscriptionHandler(topupUc)
	topupCallbackHandler := topupHandler.NewMidtransNotificationHandler(topupCallbackUc, midtransVerifier)
	topupStatusHandler := getTopUpStatusHandler.NewGetTopUpStatusHandler(getStatusTopUpUc)
	//helper
	getLokasiHandler := helperHandler.NewGetLocationHandler(GetLokasiUc)
	getTarifHandler := helperHandler.NewGetTarifHandler(GetTarifUc)
	getNominalTopUpHandler := helperHandler.NewGetNominalTopUpHandler(GetNominalTopUpUc)
	// scanHandler := mobile_handler.NewScanHandler(submitQrUC, scanDetailUC)
	// userHandler := shared_handler.NewUserHandler(getUserProfileUC)
	// paymentHandler := jukir_handler.NewPaymentHandler(initiatePaymentUC)
	// qrHandler := mobile_handler.NewQRHandler(qrUC)
	// authHandler := shared_handler.NewAuthHandler(registerUC, loginUC, refreshUC, logoutUC, changePasswordUC) // Pastikan constructor AuthHandler sudah menerima semua usecase yang diperlukan

	//OtentikasiHandler
	//auth
	// protectedLogoutHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(logoutHandler.Logout))                         // jika logout juga butuh otentikasi
	// protectedChangePasswordHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(changePasswordHandler.ChangePassword)) // Ganti cfg.JWTSecret dengan cara kamu mengakses secret
	//home
	protectedHomeHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(homeHandler.GetDashboard))
	//riwayat
	protectedRiwayatHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(riwayatHandler.Execute))
	//laporan
	protectedLaporanHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(laporanHandler.Execute))
	//subscription
	protectedSubscriptionHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(subscriptionHandler.Execute))
	//payment
	protectedPostParkingHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(postParkingHandler.Execute))
	protectedPostPaymentParkingHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(postPaymentParkingHandler.Execute))
	protectedGetPembayaranStatusHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(getPembayaranStatusHandler.Execute))
	//topup
	protectedGetTopUpStatusHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(topupStatusHandler.Execute))
	//helper
	protectedGetLokasiHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(getLokasiHandler.Execute))
	protectedGetTarifHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(getTarifHandler.Execute))
	protectedGetNominalTopUpHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(getNominalTopUpHandler.Execute))
	// protectedProfileHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(userHandler.GetCurrentUser))
	// protectedQrGenerator := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(qrHandler.GenerateQR))
	// protectedPaymentHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(paymentHandler.InitiatePayment))

	//Endpoint
	//auth
	// mux.HandleFunc("POST /api/v2/linespot/auth/register", registerHandler.Register)
	// mux.Handle("POST /api/v2/linespot/auth/login", authLimiter.AllowLogin(http.HandlerFunc(loginHandler.Login)))
	// mux.Handle("POST /api/v2/linespot/auth/refreshToken", authLimiter.AllowRefresh(http.HandlerFunc(refreshHandler.Execute)))
	// mux.Handle("POST /api/v2/linespot/auth/logout", protectedLogoutHandler)
	// // mux.Handle("GET /api/v2/linespot/users/me", protectedProfileHandler)
	// mux.Handle("POST /api/v2/linespot/auth/change-password", protectedChangePasswordHandler)

	//home
	mux.Handle("GET /api/v2/linespot/customer_home", protectedHomeHandler)
	mux.Handle("GET /api/v2/linespot/jukir_home", protectedHomeHandler)
	//laporan
	mux.Handle("POST /api/v2/linespot/laporan", protectedLaporanHandler)
	//riwayat
	mux.Handle("POST /api/v2/linespot/riwayat", protectedRiwayatHandler)
	//subscription
	mux.Handle("GET /api/v2/linespot/subscribe", protectedSubscriptionHandler)
	//parking
	mux.Handle("POST /api/v2/linespot/parking", protectedPostParkingHandler)
	//payment
	mux.Handle("POST /api/v2/linespot/parking/payment", protectedPostPaymentParkingHandler)
	mux.Handle("GET /api/v2/linespot/parking/{sessionId}/status", protectedGetPembayaranStatusHandler)
	//topup
	mux.Handle("POST /api/v2/linespot/topup/create", middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(topupCreateHandler.Execute)))
	mux.Handle("GET /api/v2/linespot/topup/{topupCode}/status", protectedGetTopUpStatusHandler)
	mux.Handle("POST /api/v2/linespot/topup/midtrans/callback", http.HandlerFunc(topupCallbackHandler.Execute))

	//helper
	mux.Handle("GET /api/v2/linespot/get_lokasi", protectedGetLokasiHandler)
	mux.Handle("GET /api/v2/linespot/get_tarif", protectedGetTarifHandler)
	mux.Handle("GET /api/v2/linespot/get_topup", protectedGetNominalTopUpHandler)

	//route dengan otentikasi (middleware JWT)
	// mux.Handle("GET /api/v2/linespot/qr/generate", protectedQrGenerator)
	// mux.Handle("POST /api/v2/linespot/pay", protectedPaymentHandler)
}
