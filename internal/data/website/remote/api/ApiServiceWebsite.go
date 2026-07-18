package api

import (
	"log"
	"net/http"
	"time"

	// "modulegue/config"

	"modulegue/core/queue"
	middleware "modulegue/internal/middleware"

	//auth
	authUc "modulegue/internal/domain/shared/usecase/auth"
	webloginUc "modulegue/internal/domain/web/usecase/login"
	webloginhandler "modulegue/internal/handler/website/auth"
)

func RegisterRoutes(
	mux *http.ServeMux,
	//auth
	loginUc *authUc.LoginUseCase,
	webloginUc *webloginUc.LoginUseCase,
	// registerUc *authUc.RegisterUseCase,
	// logoutUc *authUc.LogoutUseCase,
	// refreshUc *authUc.RefreshTokenUseCase,
	// changePasswordUc *authUc.ChangePasswordUseCase,

	//pkg
	q *queue.Queue,
	logger *log.Logger,
) {
	//core
	authLimiter := middleware.NewRateLimiter(10, time.Minute, 5)

	// Handler
	//auth
	loginHandler := webloginhandler.NewLoginHandler(webloginUc)
	// registerHandler := authHandler.NewRegisterHandler(registerUc)
	// logoutHandler := authHandler.NewLogoutHandler(logoutUc)
	// refreshHandler := authHandler.NewRefreshTokenHandler(refreshUc)
	// changePasswordHandler := authHandler.NewChangePasswordHandler(changePasswordUc)

	//OtentikasiHandler
	//auth
	// protectedLogoutHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(logoutHandler.Logout))                         // jika logout juga butuh otentikasi
	// protectedChangePasswordHandler := middleware.JWTAuth(config.Load().JWTSecret)(http.HandlerFunc(changePasswordHandler.ChangePassword)) // Ganti cfg.JWTSecret dengan cara kamu mengakses secret

	//Endpoint
	//auth
	// mux.HandleFunc("POST /api/v2/linespot/auth/register", registerHandler.Register)
	mux.Handle("POST /api/v2/dashboard/auth/login", authLimiter.AllowLogin(http.HandlerFunc(loginHandler.Login)))
	// mux.Handle("POST /api/v2/linespot/auth/refreshToken", authLimiter.AllowRefresh(http.HandlerFunc(refreshHandler.Execute)))
	// mux.Handle("POST /api/v2/linespot/auth/logout", protectedLogoutHandler)
	// mux.Handle("GET /api/v2/linespot/users/me", protectedProfileHandler)
	// mux.Handle("POST /api/v2/linespot/auth/change-password", protectedChangePasswordHandler)
}
