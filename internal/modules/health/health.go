package health

import (
	"github.com/go-chi/chi/v5"
)

func Init(r chi.Router) {
	service := NewService()
	handler := NewHandler(service)

	r.Get("/ping", handler.Ping)
	r.Post("/echo", handler.Echo)
}
