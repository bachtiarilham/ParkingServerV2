package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type rateLimitBucket struct {
	tokens   float64
	last     time.Time
	lastSeen time.Time
}

type RateLimiter struct {
	mu        sync.Mutex
	buckets   map[string]*rateLimitBucket
	rate      float64
	burst     float64
	window    time.Duration
	lastPurge time.Time
}

func NewRateLimiter(requests int, per time.Duration, burst int) *RateLimiter {
	if requests <= 0 {
		requests = 1
	}
	if per <= 0 {
		per = time.Minute
	}
	if burst <= 0 {
		burst = requests
	}
	return &RateLimiter{
		buckets: make(map[string]*rateLimitBucket),
		rate:    float64(requests) / per.Seconds(),
		burst:   float64(burst),
		window:  per * 2,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := rl.keyForRequest(r)
		if !rl.allow(key) {
			http.Error(w, "terlalu banyak request, coba lagi nanti", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) AllowLogin(next http.Handler) http.Handler {
	return rl.wrapForAuthEndpoint(next, "login")
}

func (rl *RateLimiter) AllowRefresh(next http.Handler) http.Handler {
	return rl.wrapForAuthEndpoint(next, "refresh")
}

func (rl *RateLimiter) wrapForAuthEndpoint(next http.Handler, endpoint string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := rl.keyForRequest(r) + ":" + endpoint
		if !rl.allow(key) {
			http.Error(w, fmt.Sprintf("terlalu banyak percobaan %s, coba lagi nanti", endpoint), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) keyForRequest(r *http.Request) string {
	ip := clientIP(r)
	return ip + "|" + r.Method + "|" + cleanPath(r.URL.Path)
}

func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if now.Sub(rl.lastPurge) > rl.window {
		rl.purgeLocked(now)
		rl.lastPurge = now
	}

	bucket, ok := rl.buckets[key]
	if !ok {
		rl.buckets[key] = &rateLimitBucket{
			tokens:   rl.burst - 1,
			last:     now,
			lastSeen: now,
		}
		return true
	}

	elapsed := now.Sub(bucket.last).Seconds()
	bucket.tokens = minFloat64(rl.burst, bucket.tokens+elapsed*rl.rate)
	bucket.last = now
	bucket.lastSeen = now

	if bucket.tokens < 1 {
		return false
	}

	bucket.tokens--
	return true
}

func (rl *RateLimiter) purgeLocked(now time.Time) {
	for key, bucket := range rl.buckets {
		if now.Sub(bucket.lastSeen) > rl.window {
			delete(rl.buckets, key)
		}
	}
}

func clientIP(r *http.Request) string {
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		if len(parts) > 0 {
			if ip := strings.TrimSpace(parts[0]); ip != "" {
				return ip
			}
		}
	}
	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}

func cleanPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return "/"
	}
	return path
}

func minFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
