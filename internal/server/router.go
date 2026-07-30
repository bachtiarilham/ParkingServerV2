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
	searchrepo "modulegue/internal/data/mobile/repository_impl/filter_pencarian"
	helperrepo "modulegue/internal/data/mobile/repository_impl/helper"
	homerepo "modulegue/internal/data/mobile/repository_impl/home"
	invoicerepo "modulegue/internal/data/mobile/repository_impl/invoice"
	parkingrepo "modulegue/internal/data/mobile/repository_impl/parking"
	paymentrepo "modulegue/internal/data/mobile/repository_impl/payment_parking"
	subscriptionrepo "modulegue/internal/data/mobile/repository_impl/subcription"
	topuprepo "modulegue/internal/data/mobile/repository_impl/topup"
	authrepo "modulegue/internal/data/shared/repository_impl/auth"
	webhelperrepo "modulegue/internal/data/website/repository_impl/helper"
	webloginrepo "modulegue/internal/data/website/repository_impl/home"
	webmonitoringrepo "modulegue/internal/data/website/repository_impl/monitoring"
	webpetugasrepo "modulegue/internal/data/website/repository_impl/petugas"
	websettingrepo "modulegue/internal/data/website/repository_impl/setting"
	webtopuprepo "modulegue/internal/data/website/repository_impl/topup"

	webdelivery "modulegue/internal/delivery/web"
	//usecase
	searchuc "modulegue/internal/domain/mobile/usecase/filter_pencarian"
	helperuc "modulegue/internal/domain/mobile/usecase/helper"
	homeuc "modulegue/internal/domain/mobile/usecase/home"
	invoiceuc "modulegue/internal/domain/mobile/usecase/invoice"
	parkinguc "modulegue/internal/domain/mobile/usecase/parking"
	paymentuc "modulegue/internal/domain/mobile/usecase/payment_parking"
	subscriptionuc "modulegue/internal/domain/mobile/usecase/subscription"
	topupuc "modulegue/internal/domain/mobile/usecase/topup"
	authmodel "modulegue/internal/domain/shared/model/auth"
	authuc "modulegue/internal/domain/shared/usecase/auth"
	webhelperuc "modulegue/internal/domain/web/usecase/helper"
	webloginuc "modulegue/internal/domain/web/usecase/login"
	webmonitoringuc "modulegue/internal/domain/web/usecase/monitoring"
	webpetugasuc "modulegue/internal/domain/web/usecase/petugas"
	websettinguc "modulegue/internal/domain/web/usecase/setting"
	webtopupuc "modulegue/internal/domain/web/usecase/topup"

	//payment_gate
	paymentgaterepo "modulegue/internal/data/mobile/repository_impl/payment_gate"
	paymentgateuc "modulegue/internal/domain/mobile/usecase/payment_gate"
	midtrans "modulegue/internal/service/payment_gateway/midtrans"

	//sync
	syncrepo "modulegue/internal/data/mobile/repository_impl/sync"

	//settings
	settingsrepo "modulegue/internal/data/mobile/repository_impl/settings"
	settingsuc "modulegue/internal/domain/mobile/usecase/settings"
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
	searchrepo := searchrepo.NewFilterPencarianRepoImpl(db)
	invoicerepo := invoicerepo.NewInvoiceRepositoryImpl(db)
	subscriptionRepository := subscriptionrepo.NewSubscriptionRepositoryImpl(db)
	parkingRepository := parkingrepo.NewParkingRepositoryImpl(db)
	paymentRepository := paymentrepo.NewPaymentRepositoryImpl(db)
	statusPaymentRepository := paymentrepo.NewStatusPaymentRepositoryImpl(db)
	helperRepository := helperrepo.NewHelperRepositoryImpl(db)
	topUpRepository := topuprepo.NewTopUpRepositoryImpl(db)
	paymentCallbackRepository := paymentgaterepo.NewPaymentCallbackRepositoryImpl(db)
	syncRepository := syncrepo.NewSyncRepoImpl(db)
	paymentGateRepository := paymentgaterepo.NewPaymentGateRepositoryImpl(db)
	settingsRepository := settingsrepo.NewSettingsRepositoryImpl(db)

	isProduction := cfg.AppEnv == "production"
	midtransClient := midtrans.NewMidtransClient(cfg.MidtransServerKey, isProduction)

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
	laporanuc := searchuc.NewGetLaporanUseCase(searchrepo)
	riwayatmemberuc := searchuc.NewGetRiwayatMembershipUseCase(searchrepo)
	riwayatparkiruc := searchuc.NewGetRiwayatParkirUseCase(searchrepo)
	riwayattrxuc := searchuc.NewGetRiwayatTransaksiUseCase(searchrepo)
	invoiceuc := invoiceuc.NewGetInvoiceUseCase(invoicerepo)
	subscriptionUC := subscriptionuc.NewSubscriptionUseCase(subscriptionRepository)
	postParkingUC := parkinguc.NewPostParkingUseCase(parkingRepository)
	topupUC := topupuc.NewTopUpUseCase(topUpRepository, cfg.AdminFeeTopUp)
	topupCallbackUC := topupuc.NewTopUpCallbackUseCase(topUpRepository)
	postPaymentParkingUC := paymentuc.NewPostPaymentParkingUseCase(paymentRepository)
	getPembayaranStatusUC := paymentuc.NewGetPembayaranStatusUseCase(statusPaymentRepository)
	getLokasiUC := helperuc.NewGetLokasiUseCase(helperRepository)
	getTarifUC := helperuc.NewGetTarifUseCase(helperRepository)
	GetNominalPaymentUc := helperuc.NewNominalTopUpUseCase(helperRepository)
	GetPaymentMethodUc := helperuc.NewGetPaymentMethodUseCase(helperRepository)

	getStatusTopUpUc := topupuc.NewGetTopUpStatusUseCase(topUpRepository)
	paymentCallbackUC := paymentgateuc.NewPaymentCallbackUseCase(paymentCallbackRepository, syncRepository, midtransClient)
	getPaymentStatusUC := paymentgateuc.NewGetPaymentStatusUseCase(paymentCallbackRepository)
	payTransferUC := paymentgateuc.NewPayTransferUseCase(paymentGateRepository)
	payTopUpUC := paymentgateuc.NewPayTopUpUseCase(paymentGateRepository, midtransClient)
	payMembershipUC := paymentgateuc.NewPayMembershipUseCase(paymentGateRepository, midtransClient)
	payParkingUC := paymentgateuc.NewPayParkingUseCase(paymentGateRepository, midtransClient)
	payCashParkingUc := paymentgateuc.NewPayCashParkirUseCase(paymentGateRepository, syncRepository, paymentCallbackRepository)
	changeProfileUC := settingsuc.NewChangeProfileUseCase(settingsRepository)

	//uc website
	loginWebRepo := webloginrepo.NewHomeRepositoryImpl(db)
	webHelperRepo := webhelperrepo.NewHelperRepositoryImpl(db)
	webMonitoringRepo := webmonitoringrepo.NewMonitoringRepositoryImpl(db)
	petugasWebRepo := webpetugasrepo.NewPetugasRepositoryImpl(db)
	settingWebRepo := websettingrepo.NewSettingRepositoryImpl(db)
	topupWebRepo := webtopuprepo.NewTopUpRepositoryImpl(db)

	webLoginUC := webloginuc.NewLoginUseCase(
		loginWebRepo,
		sessionRepository,
		authmodel.SessionModel{},
		cfg.JWTSecret,
		accessTTL,
		refreshTTL,
	)
	webGetLokasiUC := webhelperuc.NewGetLokasiUseCase(webHelperRepo)
	webGetTarifUC := webhelperuc.NewGetTarifUseCase(webHelperRepo)
	webMonitoringUC := webmonitoringuc.NewMonitoringUseCase(webMonitoringRepo)
	webPetugasUC := webpetugasuc.NewPetugasUseCase(petugasWebRepo)
	webAddParlokUC := websettinguc.NewAddParlokUseCase(settingWebRepo)
	webRegisterJukirUC := websettinguc.NewRegisterJukirUseCase(settingWebRepo)
	webSaveScheduleUC := websettinguc.NewSaveScheduleUseCase(settingWebRepo)
	webSaveTarifUC := websettinguc.NewSaveTarifUseCase(settingWebRepo)
	webShowSelectedJukirUC := websettinguc.NewShowSelectedJukirUseCase(settingWebRepo)
	webUpdateProfileUC := websettinguc.NewUpdateProfileUseCase(settingWebRepo)
	webTopUpUC := webtopupuc.NewTopUpUseCase(topupWebRepo)

	mobileapi.RegisterRoutes(
		mux,
		homeUC,
		laporanuc,
		riwayatmemberuc,
		riwayatparkiruc,
		riwayattrxuc,
		invoiceuc,
		paymentCallbackUC,
		getPaymentStatusUC,
		payTransferUC,
		payTopUpUC,
		payMembershipUC,
		payParkingUC,
		payCashParkingUc,
		subscriptionUC,
		postParkingUC,
		postPaymentParkingUC,
		getPembayaranStatusUC,
		topupUC,
		topupCallbackUC,
		getStatusTopUpUc,
		midtransClient,
		getLokasiUC,
		getTarifUC,
		GetNominalPaymentUc,
		GetPaymentMethodUc,
		changeProfileUC,
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
		webGetLokasiUC,
		webGetTarifUC,
		webMonitoringUC,
		webPetugasUC,
		webAddParlokUC,
		webRegisterJukirUC,
		webSaveScheduleUC,
		webSaveTarifUC,
		webShowSelectedJukirUC,
		webUpdateProfileUC,
		webTopUpUC,
		queue,
		logger,
		cfg.JWTSecret,
	)

	webdelivery.RegisterRoutes(mux, cfg, db, queue, logger)

	return mux, workerPool
}
