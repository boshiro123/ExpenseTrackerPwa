package middleware

import (
	"context"
	"net/http"
	"strings"

	"expense-tracker-pwa/internal/config"
	"expense-tracker-pwa/internal/service"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type contextKey string

var ContextUserIDKey contextKey = "user_id"

type AuthMiddleware struct {
	cfg config.Config
}

func NewAuthMiddleware(cfg config.Config) *AuthMiddleware {
	return &AuthMiddleware{cfg: cfg}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("{\"error\":\"unauthorized\"}"))
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		uid, err := service.ParseUserIDFromToken(token, m.cfg.JWTSecret)
		if err != nil || uid == primitive.NilObjectID {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("{\"error\":\"unauthorized\"}"))
			return
		}
		ctx := context.WithValue(r.Context(), ContextUserIDKey, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
