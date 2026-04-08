package middleware

import (
	"backend-challenge/internal/modules/auth"
	"backend-challenge/pkg/response"
	"net/http"
	"strings"
)

func AuthMiddleware(authService auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "missing token")
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			if !authService.ValidateToken(token) {
				response.Error(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
