package api

import (
	"log"
	"modulegue/config"
	"modulegue/core/queue"
	middleware "modulegue/internal/middleware"
	"net/http"

	//home
	homeUc "modulegue/internal/domain/mobile/usecase/home"
	homeHandler "modulegue/internal/handler/mobile/home"

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

	//search
	searchuc "modulegue/internal/domain/mobile/usecase/filter_pencarian"
	searchhandler "modulegue/internal/handler/mobile/filter_pencarian"

	//invoice
	invoiceuc "modulegue/internal/domain/mobile/usecase/invoice"
	invoicehandler "modulegue/internal/handler/mobile/invoice"

	//payment_gate
	paymentgateuc "modulegue/internal/domain/mobile/usecase/payment_gate"
	paymentgatehandler "modulegue/internal/handler/mobile/payment_gate"

	//settings
	settingsuc "modulegue/internal/domain/mobile/usecase/settings"
	settingshandler "modulegue/internal/handler/mobile/settings"
)

func RegisterRoutes(
	mux *http.ServeMux,
	//home
	homeUc *homeUc.GetHomeUseCase,
	//filter_pencarian
	laporanuc *searchuc.GetLaporanUseCase,
	riwayatmemberuc *searchuc.GetRiwayatMembershipUseCase,
	riwayatparkiruc *searchuc.GetRiwayatParkirUseCase,
	riwayattrxuc *searchuc.GetRiwayatTransaksiUseCase,
	//invoice
	invoiceuc *invoiceuc.GetInvoiceUseCase,
	//payment_gate
	paymentCallbackUc *paymentgateuc.PaymentCallbackUseCase,
	getPaymentStatusUc *paymentgateuc.GetPaymentStatusUseCase,
	payTransferUc *paymentgateuc.PayTransferUseCase,
	payTopUpUc *paymentgateuc.PayTopUpUseCase,
	payMembershipUc *paymentgateuc.PayMembershipUseCase,
	payParkingUc *paymentgateuc.PayParkingUseCase,
	payCashParkingUc *paymentgateuc.PayCashParkirUseCase,
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
	GetNominalPaymentUc *helperUc.GetNominalTopUpUseCase,
	GetPaymentMethodUc *helperUc.GetPaymentMethodUseCase,

	// setting
	changeProfileUc *settingsuc.ChangeProfileUseCase,
	//payment
	// qrUC *qruc.GenerateQRUseCase,
	// initiatePaymentUC *paymentuc.InitiatePaymentUseCase,
	//pkg
	q *queue.Queue,
	logger *log.Logger,
) {
	//home
	homeHandler := homeHandler.NewHomeHandler(homeUc)
	//filterpencarian
	searchHandler := searchhandler.NewFilterPencarianHandler(
		riwayatparkiruc,
		riwayattrxuc,
		riwayatmemberuc,
		laporanuc,
	)
	//invoice
	invoiceHandler := invoicehandler.NewGetInvoiceHandler(invoiceuc)
	//payment_gate
	paymentCallbackHandler := paymentgatehandler.NewPaymentCallbackHandler(paymentCallbackUc, getPaymentStatusUc)
	paymentGateHandler := paymentgatehandler.NewPaymentGateHandler(payTransferUc, payTopUpUc, payMembershipUc, payParkingUc, payCashParkingUc)
	//settings
	changeProfileHandler := settingshandler.NewChangeProfileHandler(changeProfileUc)
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
	getNominalPaymentHandler := helperHandler.NewGetNominalTopUpHandler(GetNominalPaymentUc)
	getPaymentMethodHandler := helperHandler.NewGetPaymentMethodHandler(GetPaymentMethodUc)

	//home
	protectedHomeHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(homeHandler.GetDashboard))
	//search
	protectedSearchHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(searchHandler.Execute))
	//invoice
	protectedInvoiceHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(invoiceHandler.Execute))
	//subscription
	protectedSubscriptionHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(subscriptionHandler.Execute))
	//parking
	protectedPostParkingHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(postParkingHandler.Execute))
	protectedPostPaymentParkingHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(postPaymentParkingHandler.Execute))
	protectedGetPembayaranStatusHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(getPembayaranStatusHandler.Execute))
	//topup
	protectedGetTopUpStatusHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(topupStatusHandler.Execute))
	//helper
	protectedGetLokasiHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(getLokasiHandler.Execute))
	protectedGetTarifHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(getTarifHandler.Execute))
	protectedGetNominalPaymentHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(getNominalPaymentHandler.Execute))
	protectedGetPaymentMethodHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(getPaymentMethodHandler.Execute))
	//images
	fileServer := http.FileServer(http.Dir("D:\\parking_data\\images"))
	//invoicepdf
	fileServerInvoice := http.FileServer(http.Dir("D:\\parking_data\\invoice"))

	//home
	mux.Handle("GET /api/v2/linespot/customer_home", protectedHomeHandler)
	mux.Handle("GET /api/v2/linespot/jukir_home", protectedHomeHandler)
	//search
	mux.Handle("POST /api/v2/linespot/search/riwayat_membership", protectedSearchHandler)
	mux.Handle("POST /api/v2/linespot/search/riwayat_parkir", protectedSearchHandler)
	mux.Handle("POST /api/v2/linespot/search/riwayat_transaksi", protectedSearchHandler)
	mux.Handle("POST /api/v2/linespot/search/laporan", protectedSearchHandler)
	//invoice
	mux.Handle("GET /api/v2/linespot/invoice/{invoice_number}", protectedInvoiceHandler)
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
	mux.Handle("GET /api/v2/linespot/get_nominal_payment", protectedGetNominalPaymentHandler)
	mux.Handle("GET /api/v2/linespot/get_payment_method", protectedGetPaymentMethodHandler)

	//payment_gate callbacks & polling
	mux.Handle("POST /api/v2/payment/midtrans/callback", http.HandlerFunc(paymentCallbackHandler.HandleCallback))
	mux.Handle("GET /api/v2/linespot/payment/status/{transaction_code}", middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(paymentCallbackHandler.GetStatus)))
	mux.Handle("POST /api/v2/linespot/payment/checkout", middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(paymentGateHandler.Execute)))

	//imgaes
	mux.Handle("GET /images/", http.StripPrefix("/images/", fileServer))

	//invoicepdf
	mux.Handle("GET /invoice/", http.StripPrefix("/invoice/", fileServerInvoice))

	//settings
	mux.Handle("PUT /api/v2/linespot/change_profile", middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(changeProfileHandler.Execute)))

}
