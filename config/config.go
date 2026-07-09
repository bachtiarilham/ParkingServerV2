//struct + loader dari .env

package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	AppEnv             string
	AppPort            string
	DBHost             string
	DBPort             string
	DBName             string
	DBUser             string
	DBPassword         string
	JWTSecret          string
	QueueSQLitePath    string
	QueueWorkerCount   int
	QueueLeaseSeconds  int
	QueueRetrySeconds  int
	QueueWaitTimeoutMS int
	AccessTokenMinutes time.Duration
	RefreshTokenHours  time.Duration
	OTPExpiryMinutes   int
}

func Load() Config {
	return Config{
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),
		// DBHost:             getEnv("DB_HOST", "127.0.0.1"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "3306"),
		DBName:             getEnv("DB_NAME", "parkir_v3"),
		DBUser:             getEnv("DB_USER", "root"),
		DBPassword:         getEnvAny([]string{"DB_PASSWORD", "DB_PASS"}, "root"),
		JWTSecret:          getEnv("JWT_SECRET", "saya-janji-ngoding-pakai-doa-dan-kesabaran-ekstra"),
		QueueSQLitePath:    getEnv("QUEUE_SQLITE_PATH", "./var/park-server-queue.db"),
		QueueWorkerCount:   getEnvAsInt("QUEUE_WORKER_COUNT", 4),
		QueueLeaseSeconds:  getEnvAsInt("QUEUE_LEASE_SECONDS", 30),
		QueueRetrySeconds:  getEnvAsInt("QUEUE_RETRY_SECONDS", 5),
		QueueWaitTimeoutMS: getEnvAsInt("QUEUE_WAIT_TIMEOUT_MS", 1500),
		AccessTokenMinutes: time.Duration(getEnvAsInt("ACCESS_TOKEN_MINUTES", 60)) * time.Minute,
		RefreshTokenHours:  time.Duration(getEnvAsInt("REFRESH_TOKEN_HOURS", 24*30)) * time.Hour,
		OTPExpiryMinutes:   getEnvAsInt("OTP_EXPIRY_MINUTES", 10),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAny(keys []string, fallback string) string {
	for _, key := range keys {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	var parsed int
	if _, err := fmt.Sscanf(value, "%d", &parsed); err != nil {
		return fallback
	}
	return parsed
}
