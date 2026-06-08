package middleware

import (
	"log"
	"net/http"
	"time"

	"modulegue/pkg/logger"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func WithRequestLogger(next http.Handler, l *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		logger.WriteRequest(l, logger.RequestLogEntry{
			Method:     r.Method,
			Path:       r.URL.Path,
			StatusCode: rec.status,
			Duration:   time.Since(start),
		})
	})
}
