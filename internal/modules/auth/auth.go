package auth

import "github.com/go-chi/chi/v5"

func Init(r chi.Router, authService Service) {
	handler := NewHandler(authService)
	r.Post("/auth/token", handler.GenerateToken)
}
