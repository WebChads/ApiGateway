package router

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/WebChads/ApiGateway/internal/pkg/middleware/circuit"
	"github.com/WebChads/ApiGateway/internal/pkg/middleware/service"
	"github.com/go-chi/chi"
	"github.com/sony/gobreaker"
)

func TournamentServiceRouter() http.Handler {
	productServiceURL, _ := url.Parse("http://tournament-service:8082")
	proxy := httputil.NewSingleHostReverseProxy(productServiceURL)

	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "tournament-service",
		MaxRequests: 5,
		Interval:    10 * time.Second,
		Timeout:     15 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5
		},
	})

	proxy.Director = func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", productServiceURL.Host)
		req.URL.Scheme = productServiceURL.Scheme
		req.URL.Host = productServiceURL.Host
		req.URL.Path = "/api/v1" + strings.TrimPrefix(req.URL.Path, "/tournament")
	}

	r := chi.NewRouter()
	r.Use(service.ServiceMiddleware("tournament-service"))
	r.Handle("/*", circuit.CircuitBreakerMiddleware(cb, proxy))
	return r
}
