package service

import "net/http"

func ServiceMiddleware(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Add("X-API-Gateway", "true")
			r.Header.Add("X-Target-Service", serviceName)
			next.ServeHTTP(w, r)
		})
	}
}
