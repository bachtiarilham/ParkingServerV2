package server

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"modulegue/config"
	"modulegue/core/queue"
	"modulegue/core/worker"
	mobileapi "modulegue/internal/data/mobile/remote/api"
	authrepo "modulegue/internal/data/mobile/repository_impl/auth"
	helperrepo "modulegue/internal/data/mobile/repository_impl/helper"
	homerepo "modulegue/internal/data/mobile/repository_impl/home"
	laporanrepo "modulegue/internal/data/mobile/repository_impl/laporan"
	paymentrepo "modulegue/internal/data/mobile/repository_impl/payment"
	riwayatrepo "modulegue/internal/data/mobile/repository_impl/riwayat"
	subscriptionrepo "modulegue/internal/data/mobile/repository_impl/subcription"
	webdelivery "modulegue/internal/delivery/web"
	authmodel "modulegue/internal/domain/mobile/model/auth"
	authuc "modulegue/internal/domain/mobile/usecase/auth"
	helperuc "modulegue/internal/domain/mobile/usecase/helper"
	homeuc "modulegue/internal/domain/mobile/usecase/home"
	laporanuc "modulegue/internal/domain/mobile/usecase/laporan"
	paymentuc "modulegue/internal/domain/mobile/usecase/payment"
	riwayatuc "modulegue/internal/domain/mobile/usecase/riwayat"
	subscriptionuc "modulegue/internal/domain/mobile/usecase/subscription"
)

func NewRouter(
	cfg *config.Config,
	db *sql.DB,
	queue *queue.Queue,
	logger *log.Logger,
) (http.Handler, *worker.WorkerPool) {
	mux := http.NewServeMux()
	workerPool := worker.NewWorkerPool(queue, cfg.QueueWorkerCount, logger)

	accessTTL := time.Duration(cfg.AccessTokenMinutes) * time.Minute
	refreshTTL := time.Duration(cfg.RefreshTokenHours) * time.Hour

	authRepository := authrepo.NewAuthRepositoryImpl(db)
	sessionRepository := authrepo.NewSessionRepositoryImpl(db)
	homeRepository := homerepo.NewHomeRepositoryImpl(db)
	laporanRepository := laporanrepo.NewLaporanRepositoryImpl(db)
	riwayatRepository := riwayatrepo.NewRiwayatRepositoryImpl(db)
	subscriptionRepository := subscriptionrepo.NewSubscriptionRepositoryImpl(db)
	paymentRepository := paymentrepo.NewPaymentRepositoryImpl(db)
	helperRepository := helperrepo.NewHelperRepositoryImpl(db)

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
	postParkingUC := paymentuc.NewPostParkingUseCase(paymentRepository)
	postPaymentParkingUC := paymentuc.NewPostPaymentParkingUseCase(paymentRepository)
	getPembayaranStatusUC := paymentuc.NewGetPembayaranStatusUseCase(paymentRepository)
	getLokasiUC := helperuc.NewGetLokasiUseCase(helperRepository)
	getTarifUC := helperuc.NewGetTarifUseCase(helperRepository)

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
		getLokasiUC,
		getTarifUC,
		queue,
		logger,
	)

	webdelivery.RegisterRoutes(mux, cfg, db, queue, logger)

	return mux, workerPool
}
