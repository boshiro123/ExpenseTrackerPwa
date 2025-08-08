package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"expense-tracker-pwa/internal/middleware"
	"expense-tracker-pwa/internal/service"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExpenseController struct {
	service *service.ExpenseService
}

func NewExpenseController(s *service.ExpenseService) *ExpenseController {
	return &ExpenseController{service: s}
}

func (c *ExpenseController) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserIDKey).(primitive.ObjectID)
	ctx, cancel := timeoutCtx(r)
	defer cancel()
	items, err := c.service.List(ctx, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "server_error"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(items)
}

func (c *ExpenseController) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserIDKey).(primitive.ObjectID)
	var in service.CreateExpenseInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Amount <= 0 || in.Category == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid_input"})
		return
	}
	if in.Date.IsZero() {
		in.Date = time.Now()
	}
	ctx, cancel := timeoutCtx(r)
	defer cancel()
	created, err := c.service.Create(ctx, userID, in)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "server_error"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(created)
}
