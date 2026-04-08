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

    // Pasang middleware Chi untuk menangani trailing slashes secara global di Server.go
    // r.Use(middleware.StripSlashes) 

    r.Route("/books", func(r chi.Router) {
        // Route tanpa parameter di atas
        r.Post("/", handler.Create)
        r.Get("/", handler.FindAll)

        // Route dengan parameter di bawah
        r.Get("/{id}", handler.FindByID)
        r.Put("/{id}", handler.Update)
        r.Delete("/{id}", handler.Delete)

        r.Group(func(r chi.Router) {
            r.Use(middleware.AuthMiddleware(authService))
            r.Get("/protected", handler.FindAll)
        })
    })
}
