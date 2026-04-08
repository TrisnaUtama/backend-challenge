package books

import (
	middleware "backend-challenge/internal/middlewares"
	"backend-challenge/internal/modules/auth"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Init(r chi.Router, db *pgxpool.Pool, authService auth.Service) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.Route("/books", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.FindAll)
		r.Get("/{id}", handler.FindByID)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(authService))
			r.Get("/protected", handler.FindAll)
		})
	})
}
