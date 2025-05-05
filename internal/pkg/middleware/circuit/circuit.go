package circuit

import (
	"net/http"

	"github.com/sony/gobreaker"
)

func CircuitBreakerMiddleware(cb *gobreaker.CircuitBreaker, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := cb.Execute(func() (any, error) {
			next.ServeHTTP(w, r)
			return nil, nil
		})

		if err != nil {
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		}
	})
}
