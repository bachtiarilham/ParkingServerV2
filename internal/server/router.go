package server

import (
	"database/sql"
	"log"
	"net/http"

	"modulegue/config"
	"modulegue/core/queue"
	"modulegue/core/worker"

	//api
	mobileapi "modulegue/internal/data/mobile/remote/api"
	sharedapi "modulegue/internal/data/shared/remote/api"
	webapi "modulegue/internal/data/website/remote/api"

	//repo
	helperrepo "modulegue/internal/data/mobile/repository_impl/helper"
	homerepo "modulegue/internal/data/mobile/repository_impl/home"
	laporanrepo "modulegue/internal/data/mobile/repository_impl/laporan"
	parkingrepo "modulegue/internal/data/mobile/repository_impl/parking"
	paymentrepo "modulegue/internal/data/mobile/repository_impl/payment_parking"
	riwayatrepo "modulegue/internal/data/mobile/repository_impl/riwayat"
	subscriptionrepo "modulegue/internal/data/mobile/repository_impl/subcription"
	topuprepo "modulegue/internal/data/mobile/repository_impl/topup"
	authrepo "modulegue/internal/data/shared/repository_impl/auth"
	webloginrepo "modulegue/internal/data/website/repository_impl/home"

	webdelivery "modulegue/internal/delivery/web"
	//usecase
	helperuc "modulegue/internal/domain/mobile/usecase/helper"
	homeuc "modulegue/internal/domain/mobile/usecase/home"
	laporanuc "modulegue/internal/domain/mobile/usecase/laporan"
	parkinguc "modulegue/internal/domain/mobile/usecase/parking"
	paymentuc "modulegue/internal/domain/mobile/usecase/payment_parking"
	riwayatuc "modulegue/internal/domain/mobile/usecase/riwayat"
	subscriptionuc "modulegue/internal/domain/mobile/usecase/subscription"
	topupuc "modulegue/internal/domain/mobile/usecase/topup"
	authmodel "modulegue/internal/domain/shared/model/auth"
	authuc "modulegue/internal/domain/shared/usecase/auth"
	webloginuc "modulegue/internal/domain/web/usecase/login"
	"modulegue/internal/service/payment_gateway"
)

func NewRouter(
	cfg *config.Config,
	db *sql.DB,
	queue *queue.Queue,
	logger *log.Logger,
) (http.Handler, *worker.WorkerPool) {
	mux := http.NewServeMux()
	workerPool := worker.NewWorkerPool(queue, cfg.QueueWorkerCount, logger)

	accessTTL := cfg.AccessTokenMinutes
	refreshTTL := cfg.RefreshTokenHours

	authRepository := authrepo.NewAuthRepositoryImpl(db)
	sessionRepository := authrepo.NewSessionRepositoryImpl(db)
	homeRepository := homerepo.NewHomeRepositoryImpl(db)
	laporanRepository := laporanrepo.NewLaporanRepositoryImpl(db)
	riwayatRepository := riwayatrepo.NewRiwayatRepositoryImpl(db)
	subscriptionRepository := subscriptionrepo.NewSubscriptionRepositoryImpl(db)
	parkingRepository := parkingrepo.NewParkingRepositoryImpl(db)
	paymentRepository := paymentrepo.NewPaymentRepositoryImpl(db)
	statusPaymentRepository := paymentrepo.NewStatusPaymentRepositoryImpl(db)
	helperRepository := helperrepo.NewHelperRepositoryImpl(db)
	topUpRepository := topuprepo.NewTopUpRepositoryImpl(db)

	loginWebRepo := webloginrepo.NewHomeRepositoryImpl(db)

	var midtransService *payment_gateway.MidtransService
	if cfg.MidtransServerKey != "" && cfg.MidtransClientKey != "" {
		service, err := payment_gateway.NewMidtransService(payment_gateway.MidtransConfig{
			ServerKey:   cfg.MidtransServerKey,
			ClientKey:   cfg.MidtransClientKey,
			Environment: cfg.AppEnv,
		})
		if err != nil {
			logger.Printf("midtrans service disabled: %v", err)
		} else {
			midtransService = service
		}
	}

	//uc shared
	registerUC := authuc.NewRegisterUseCase(authRepository)
	loginUC := authuc.NewLoginUseCase(
		authRepository,
		sessionRepository,
		authmodel.SessionModel{},
		cfg.JWTSecret,
		accessTTL,
		refreshTTL,
	)
	logoutUC := authuc.NewLogoutUseCase(authRepository, sessionRepository)
	refreshUC := authuc.NewRefreshTokenUseCase(
		authRepository,
		sessionRepository,
		cfg.JWTSecret,
		accessTTL,
		refreshTTL,
	)
	changePasswordUC := authuc.NewChangePasswordUseCase(authRepository, sessionRepository)
	//uc mobile
	homeUC := homeuc.NewGetHomeUseCase(homeRepository)
	laporanUC := laporanuc.NewGetLaporanUseCase(laporanRepository)
	riwayatUC := riwayatuc.NewGetRiwayatUseCase(riwayatRepository)
	subscriptionUC := subscriptionuc.NewSubscriptionUseCase(subscriptionRepository)
	postParkingUC := parkinguc.NewPostParkingUseCase(parkingRepository)
	topupUC := topupuc.NewTopUpUseCase(topUpRepository, cfg.AdminFeeTopUp)
	topupCallbackUC := topupuc.NewTopUpCallbackUseCase(topUpRepository)
	postPaymentParkingUC := paymentuc.NewPostPaymentParkingUseCase(paymentRepository)
	getPembayaranStatusUC := paymentuc.NewGetPembayaranStatusUseCase(statusPaymentRepository)
	getLokasiUC := helperuc.NewGetLokasiUseCase(helperRepository)
	getTarifUC := helperuc.NewGetTarifUseCase(helperRepository)
	GetNominalTopUpUc := helperuc.NewNominalTopUpUseCase(helperRepository)
	getStatusTopUpUc := topupuc.NewGetTopUpStatusUseCase(topUpRepository)

	//uc website
	webLoginUC := webloginuc.NewLoginUseCase(
		loginWebRepo,
		sessionRepository,
		authmodel.SessionModel{},
		cfg.JWTSecret,
		accessTTL,
		refreshTTL,
	)

	mobileapi.RegisterRoutes(
		mux,

		homeUC,
		laporanUC,
		riwayatUC,
		subscriptionUC,
		postParkingUC,
		postPaymentParkingUC,
		getPembayaranStatusUC,
		topupUC,
		topupCallbackUC,
		getStatusTopUpUc,
		midtransService,
		getLokasiUC,
		getTarifUC,
		GetNominalTopUpUc,
		queue,
		logger,
	)

	sharedapi.RegisterRoutes(
		mux,
		loginUC,
		registerUC,
		logoutUC,
		refreshUC,
		changePasswordUC,
		queue,
		logger,
	)

	webapi.RegisterRoutes(
		mux,
		loginUC,
		webLoginUC,
		queue,
		logger,
	)

	webdelivery.RegisterRoutes(mux, cfg, db, queue, logger)

	return mux, workerPool
}
