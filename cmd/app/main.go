package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"expense-tracker-pwa/internal/config"
	"expense-tracker-pwa/internal/controller"
	"expense-tracker-pwa/internal/logger"
	appmw "expense-tracker-pwa/internal/middleware"
	"expense-tracker-pwa/internal/service"

	chi "github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load()
	logger.Init()
	log := logger.Log

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := config.ConnectMongo(ctx, cfg)

	authService := service.NewAuthService(db, cfg)
	expenseService := service.NewExpenseService(db)
	authController := controller.NewAuthController(authService)
	expenseController := controller.NewExpenseController(expenseService)

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	m := appmw.NewAuthMiddleware(cfg)

	r.Route("/api", func(r chi.Router) {
		r.Use(appmw.JSONContentType)
		r.Post("/register", authController.Register)
		r.Post("/login", authController.Login)
		r.Group(func(r chi.Router) {
			r.Use(m.RequireAuth)
			r.Get("/expenses", expenseController.List)
			r.Post("/expenses", expenseController.Create)
		})
	})

	r.Handle("/manifest.json", http.FileServer(http.Dir("frontend")))
	r.Handle("/service-worker.js", http.FileServer(http.Dir("frontend")))
	r.Handle("/icons/*", http.StripPrefix("/icons/", http.FileServer(http.Dir("frontend/icons"))))
	r.Handle("/*", http.FileServer(http.Dir("frontend")))

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r}
	go func() {
		log.Info("server_start", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server_error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	_ = srv.Shutdown(ctxShutdown)
	log.Info("server_stopped")
}
