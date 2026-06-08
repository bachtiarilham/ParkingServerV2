package web

import (
	"database/sql"
	"log"
	"net/http"

	// "modulegue/internal/delivery/web/handler"
	"modulegue/pkg/queue"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB, q *queue.Queue, logger *log.Logger) {
	// authHandler := handler.NewAuthHandler(db, q, logger) // buat handler
	// Daftarkan endpoint
	// mux.HandleFunc("POST /auth/register", authHandler.Register)
	// endpoint mobile lainnya...
}
