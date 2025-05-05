package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/WebChads/ApiGateway/internal/pkg/middleware/auth"
	router "github.com/WebChads/ApiGateway/internal/routers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func main() {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(auth.AuthMiddleware)
	r.Use(httprate.Limit(
		100,
		time.Minute,
		httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
			return r.RemoteAddr, nil
		}),
	))

	// Service routes
	r.Mount("/api/v1/auth", router.AuthServiceRouter())
	r.Mount("/api/v1/accounts", router.AccountServiceRouter())
	r.Mount("/api/v1/tournaments", router.TournamentServiceRouter())

	// Health check endpoint
	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	slog.Info("Starting API Gateway on :8080")
	http.ListenAndServe(":8080", r)
}
