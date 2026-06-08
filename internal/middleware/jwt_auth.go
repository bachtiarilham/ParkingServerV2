//Bearer token (mobile + desktop)

package middleware

import (
	"context"
	"net/http"
	"strings"

	"modulegue/pkg/jwt"
	"modulegue/pkg/response"
)

type contextKey string

const ContextKeyUserID contextKey = "user_id"

// JWTAuth memvalidasi Bearer token dan menyimpan userID ke context
func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "token tidak ditemukan")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				response.Error(w, http.StatusUnauthorized, "format token tidak valid")
				return
			}

			claims, err := jwt.ParseClaimsHS256(parts[1], secret)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "token tidak valid atau sudah kadaluarsa")
				return
			}

			if claims.Type != "access" {
				response.Error(w, http.StatusUnauthorized, "tipe token tidak valid")
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromContext mengambil userID dari context hasil middleware
func UserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(ContextKeyUserID).(int64)
	return id, ok
}
