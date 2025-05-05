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

func AccountServiceRouter() http.Handler {
	userServiceURL, _ := url.Parse("http://account-service:8081")
	proxy := httputil.NewSingleHostReverseProxy(userServiceURL)

	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "account-service",
		MaxRequests: 5,
		Interval:    10 * time.Second,
		Timeout:     15 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5
		},
	})

	proxy.Director = func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", userServiceURL.Host)
		req.URL.Scheme = userServiceURL.Scheme
		req.URL.Host = userServiceURL.Host
		req.URL.Path = "/api/v1" + strings.TrimPrefix(req.URL.Path, "/account")
	}

	r := chi.NewRouter()
	r.Use(service.ServiceMiddleware("account-service"))
	r.Handle("/*", circuit.CircuitBreakerMiddleware(cb, proxy))
	return r
}
