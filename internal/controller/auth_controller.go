package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"expense-tracker-pwa/internal/service"
)

type AuthController struct {
	service *service.AuthService
}

func NewAuthController(s *service.AuthService) *AuthController {
	return &AuthController{service: s}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var in service.AuthCredentials
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Email == "" || in.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid_input"})
		return
	}
	ctx, cancel := timeoutCtx(r)
	defer cancel()
	if err := c.service.Register(ctx, in.Email, in.Password); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "user_exists"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var in service.AuthCredentials
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Email == "" || in.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid_input"})
		return
	}
	ctx, cancel := timeoutCtx(r)
	defer cancel()
	token, err := c.service.Login(ctx, in.Email, in.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid_credentials"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(service.AuthToken{Token: token})
}

func timeoutCtx(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), 5*time.Second)
}
