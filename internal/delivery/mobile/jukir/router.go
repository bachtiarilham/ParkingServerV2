//mount semua route mobile

package jukir

import (
	"log"
	"modulegue/config"
	"net/http"

	jukir_handler "modulegue/internal/delivery/mobile/jukir/handler"
	paymentuc "modulegue/internal/usecase/payment"

	middleware "modulegue/internal/middleware"
	authuc "modulegue/internal/usecase/auth"
	homeuc "modulegue/internal/usecase/home_jukir"
	useruc "modulegue/internal/usecase/user"

	"modulegue/pkg/queue"
)

func RegisterRoutes(
	mux *http.ServeMux,
	registerUC *authuc.RegisterUseCase,
	loginUC *authuc.LoginUseCase,
	logoutUC *authuc.LogoutUseCase,
	refreshUC *authuc.RefreshTokenUseCase,
	changePasswordUC *authuc.ChangePasswordUseCase,
	getUserProfileUC *useruc.GetProfileUseCase,
	homeUC *homeuc.GetDashboardUseCase,
	initiatePaymentUC *paymentuc.InitiatePaymentUseCase,
	q *queue.Queue,
	logger *log.Logger,
) {

	paymentHandler := jukir_handler.NewPaymentHandler(initiatePaymentUC)

	protectedPaymentHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(paymentHandler.InitiatePayment))

	mux.Handle("POST /api/v2/linespot/pay", protectedPaymentHandler)
}
