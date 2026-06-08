//mount semua route mobile

package jukir

import (
	"database/sql"
	"log"
	"net/http"

	"modulegue/internal/delivery/mobile/jukir/handler"
	"modulegue/pkg/queue"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB, q *queue.Queue, logger *log.Logger) {
	authHandler := handler.NewAuthHandler(db, q, logger) // buat handler
	// Daftarkan endpoint
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	// endpoint mobile lainnya...
}
