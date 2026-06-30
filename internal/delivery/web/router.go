package web

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"modulegue/config"
	webhandler "modulegue/internal/delivery/web/handler"
	"modulegue/internal/middleware"
	sharedrepo "modulegue/internal/repository"
	webrepo "modulegue/internal/repository/web"
	webuc "modulegue/internal/usecase/web"
	"modulegue/pkg/queue"
)

func RegisterRoutes(mux *http.ServeMux, cfg *config.Config, db *sql.DB, q *queue.Queue, logger *log.Logger) {
	_ = q
	_ = logger

	adminRepo := webrepo.NewMySQL(db)
	authRepo := sharedrepo.NewAuthRepository(db)

	authUC := webuc.NewAuthUseCase(adminRepo, authRepo, cfg.JWTSecret, cfg.AccessTokenMinutes, cfg.RefreshTokenHours)
	dashboardUC := webuc.NewDashboardUseCase(adminRepo)

	authHandler := webhandler.NewAuthHandler(authUC)
	dashboardHandler := webhandler.NewDashboardHandler(dashboardUC)
	authLimiter := middleware.NewRateLimiter(10, time.Minute, 5)

	mux.Handle("POST /api/v2/dashboard-parkir/auth/login", authLimiter.AllowLogin(http.HandlerFunc(authHandler.Login)))
	mux.Handle("POST /api/v2/dashboard-parkir/auth/refresh", authLimiter.AllowRefresh(http.HandlerFunc(authHandler.Refresh)))

	protected := middleware.JWTAuth(cfg.JWTSecret)
	mux.Handle("GET /api/v2/dashboard-parkir/dashboard/overview", protected(http.HandlerFunc(dashboardHandler.GetDashboardOverview)))
	mux.Handle("GET /api/v2/dashboard-parkir/monitoring/overview", protected(http.HandlerFunc(dashboardHandler.GetMonitoringOverview)))
	mux.Handle("GET /api/v2/dashboard-parkir/officers/overview", protected(http.HandlerFunc(dashboardHandler.GetOfficerOverview)))
	mux.Handle("GET /api/v2/dashboard-parkir/transactions/overview", protected(http.HandlerFunc(dashboardHandler.GetTransactionsOverview)))
	mux.Handle("GET /api/v2/dashboard-parkir/settings/overview", protected(http.HandlerFunc(dashboardHandler.GetSettingsOverview)))
	mux.Handle("PUT /api/v2/dashboard-parkir/settings/overview", protected(http.HandlerFunc(dashboardHandler.UpdateSettingsOverview)))
	mux.Handle("PATCH /api/v2/dashboard-parkir/monitoring/locations/{id}", protected(http.HandlerFunc(dashboardHandler.UpdateLocationSettings)))
	mux.Handle("PUT /api/v2/dashboard-parkir/monitoring/locations/{id}/shift-templates", protected(http.HandlerFunc(dashboardHandler.SaveLocationShiftTemplates)))
	mux.Handle("PATCH /api/v2/dashboard-parkir/officers/{id}/status", protected(http.HandlerFunc(dashboardHandler.UpdateOfficerStatus)))
	mux.Handle("POST /api/v2/dashboard-parkir/officers/mutations", protected(http.HandlerFunc(dashboardHandler.ApplyOfficerMutation)))
	mux.Handle("POST /api/v2/dashboard-parkir/transactions/disputes", protected(http.HandlerFunc(dashboardHandler.CreateDisputeCase)))
	mux.Handle("PATCH /api/v2/dashboard-parkir/transactions/disputes/{id}", protected(http.HandlerFunc(dashboardHandler.UpdateDisputeCaseStatus)))
	mux.Handle("POST /api/v2/dashboard-parkir/transactions/refunds", protected(http.HandlerFunc(dashboardHandler.CreateRefundTransaction)))
	mux.Handle("PATCH /api/v2/dashboard-parkir/transactions/refunds/{id}", protected(http.HandlerFunc(dashboardHandler.UpdateRefundTransactionStatus)))
	mux.Handle("POST /api/v2/dashboard-parkir/transactions/closings", protected(http.HandlerFunc(dashboardHandler.CreateClosingBatch)))
	mux.Handle("PATCH /api/v2/dashboard-parkir/transactions/closings/{id}", protected(http.HandlerFunc(dashboardHandler.UpdateClosingBatchStatus)))
}
