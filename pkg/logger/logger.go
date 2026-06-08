package logger

import (
	"log"
	"time"
)

type RequestLogEntry struct {
	Method     string
	Path       string
	StatusCode int
	Duration   time.Duration
}

func WriteRequest(l *log.Logger, entry RequestLogEntry) {
	if l == nil {
		return
	}

	l.Printf(
		"request method=%s path=%s status=%d duration_ms=%d",
		entry.Method,
		entry.Path,
		entry.StatusCode,
		entry.Duration.Milliseconds(),
	)
}
