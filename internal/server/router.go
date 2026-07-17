package server

import (
	"database/sql"
	"log"
	"net/http"

	"modulegue/config"
	"modulegue/core/queue"
	"modulegue/core/worker"
	mobileapi "modulegue/internal/data/mobile/remote/api"
	authrepo "modulegue/internal/data/mobile/repository_impl/auth"
	helperrepo "modulegue/internal/data/mobile/repository_impl/helper"
	homerepo "modulegue/internal/data/mobile/repository_impl/home"
	laporanrepo "modulegue/internal/data/mobile/repository_impl/laporan"
	parkingrepo "modulegue/internal/data/mobile/repository_impl/parking"
	paymentrepo "modulegue/internal/data/mobile/repository_impl/payment_parking"
	riwayatrepo "modulegue/internal/data/mobile/repository_impl/riwayat"
	subscriptionrepo "modulegue/internal/data/mobile/repository_impl/subcription"
	topuprepo "modulegue/internal/data/mobile/repository_impl/topup"
	webdelivery "modulegue/internal/delivery/web"
	authmodel "modulegue/internal/domain/mobile/model/auth"
	authuc "modulegue/internal/domain/mobile/usecase/auth"
	helperuc "modulegue/internal/domain/mobile/usecase/helper"
	homeuc "modulegue/internal/domain/mobile/usecase/home"
	laporanuc "modulegue/internal/domain/mobile/usecase/laporan"
	parkinguc "modulegue/internal/domain/mobile/usecase/parking"
	paymentuc "modulegue/internal/domain/mobile/usecase/payment_parking"
	riwayatuc "modulegue/internal/domain/mobile/usecase/riwayat"
	subscriptionuc "modulegue/internal/domain/mobile/usecase/subscription"
	topupuc "modulegue/internal/domain/mobile/usecase/topup"
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

	mobileapi.RegisterRoutes(
		mux,
		loginUC,
		registerUC,
		logoutUC,
		refreshUC,
		changePasswordUC,
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

	webdelivery.RegisterRoutes(mux, cfg, db, queue, logger)

	return mux, workerPool
}
