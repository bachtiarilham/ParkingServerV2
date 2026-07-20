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
	webhelperuc "modulegue/internal/domain/web/usecase/helper"
	webloginUc "modulegue/internal/domain/web/usecase/login"
	webmonitoringuc "modulegue/internal/domain/web/usecase/monitoring"
	webpetugasuc "modulegue/internal/domain/web/usecase/petugas"
	websettinguc "modulegue/internal/domain/web/usecase/setting"
	webtopupuc "modulegue/internal/domain/web/usecase/topup"
	webloginhandler "modulegue/internal/handler/website/auth"
	helperhandler "modulegue/internal/handler/website/helper"
	monitoringhandler "modulegue/internal/handler/website/monitoring"
	petugashandler "modulegue/internal/handler/website/petugas"
	settinghandler "modulegue/internal/handler/website/setting"
	topuphandler "modulegue/internal/handler/website/topup"
)

func RegisterRoutes(
	mux *http.ServeMux,
	//auth
	loginUc *authUc.LoginUseCase,
	webloginUc *webloginUc.LoginUseCase,
	getLokasiUc *webhelperuc.GetLokasiUseCase,
	getTarifUc *webhelperuc.GetTarifUseCase,
	monitoringUc *webmonitoringuc.MonitoringUseCase,
	petugasUc *webpetugasuc.PetugasUseCase,
	addParlokUc *websettinguc.AddParlokUseCase,
	registerJukirUc *websettinguc.RegisterJukirUseCase,
	saveScheduleUc *websettinguc.SaveScheduleUseCase,
	saveTarifUc *websettinguc.SaveTarifUseCase,
	showSelectedJukirUc *websettinguc.ShowSelectedJukirUseCase,
	updateProfileUc *websettinguc.UpdateProfileUseCase,
	topupUc *webtopupuc.TopUpUseCase,
	// registerUc *authUc.RegisterUseCase,
	// logoutUc *authUc.LogoutUseCase,
	// refreshUc *authUc.RefreshTokenUseCase,
	// changePasswordUc *authUc.ChangePasswordUseCase,

	//pkg
	q *queue.Queue,
	logger *log.Logger,
	jwtSecret string,
) {
	//core
	authLimiter := middleware.NewRateLimiter(10, time.Minute, 5)

	// Handler
	//auth
	loginHandler := webloginhandler.NewLoginHandler(webloginUc)
	getLokasiHandler := helperhandler.NewGetLokasiHandler(getLokasiUc)
	getTarifHandler := helperhandler.NewGetTarifHandler(getTarifUc)
	monitoringHandler := monitoringhandler.NewGetLokasiHandler(monitoringUc)
	petugasHandler := petugashandler.NewPetugasHandler(petugasUc)
	settingHandler := settinghandler.NewSettingHandler(
		addParlokUc,
		registerJukirUc,
		saveScheduleUc,
		saveTarifUc,
		showSelectedJukirUc,
		updateProfileUc,
	)
	topupHandler := topuphandler.NewTopUpHandler(topupUc)
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
	mux.Handle("POST /api/web/login", authLimiter.AllowLogin(http.HandlerFunc(loginHandler.Login)))
	// mux.Handle("POST /api/v2/linespot/auth/refreshToken", authLimiter.AllowRefresh(http.HandlerFunc(refreshHandler.Execute)))
	// mux.Handle("POST /api/v2/linespot/auth/logout", protectedLogoutHandler)
	// mux.Handle("GET /api/v2/linespot/users/me", protectedProfileHandler)
	// mux.Handle("POST /api/v2/linespot/auth/change-password", protectedChangePasswordHandler)

	protected := middleware.JWTAuth(jwtSecret)
	mux.Handle("POST /api/web/oldtarif", protected(http.HandlerFunc(getTarifHandler.Execute)))
	mux.Handle("POST /api/web/getlokasilist", protected(http.HandlerFunc(getLokasiHandler.Execute)))
	mux.Handle("POST /api/web/monitoring", protected(http.HandlerFunc(monitoringHandler.Execute)))
	mux.Handle("POST /api/web/petugas", protected(http.HandlerFunc(petugasHandler.GetPetugas)))
	mux.Handle("POST /api/web/addparlok", protected(http.HandlerFunc(settingHandler.AddParlok)))
	mux.Handle("POST /api/web/register/jukir", protected(http.HandlerFunc(settingHandler.RegisterJukir)))
	mux.Handle("POST /api/web/saveschedule", protected(http.HandlerFunc(settingHandler.SaveSchedule)))
	mux.Handle("POST /api/web/savetarif", protected(http.HandlerFunc(settingHandler.SaveTarif)))
	mux.Handle("POST /api/web/showselectedjukir", protected(http.HandlerFunc(settingHandler.ShowSelectedJukir)))
	mux.Handle("POST /api/web/updateprofil", protected(http.HandlerFunc(settingHandler.UpdateProfil)))
	mux.Handle("POST /api/web/topup", protected(http.HandlerFunc(topupHandler.Execute)))
}
