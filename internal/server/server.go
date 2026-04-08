package server

import (
	"fmt"
	"net/http"

	"backend-challenge/internal/configs"
	"backend-challenge/internal/database"
	"backend-challenge/internal/modules/auth"
	"backend-challenge/internal/modules/books"
	"backend-challenge/internal/modules/health"
	"backend-challenge/pkg/docs"
	"backend-challenge/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

type Server struct {
	router *chi.Mux
	cfg    *configs.Setting
}

func New(cfg *configs.Setting, db *database.PostgresDB) *Server {
	s := &Server{
		router: chi.NewRouter(),
		cfg:    cfg,
	}
	s.setupMiddleware()
	s.setupRoutes(db)
	return s
}

func (s *Server) setupMiddleware() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   s.cfg.Cors.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))
}

func (s *Server) setupRoutes(db *database.PostgresDB) {
	authService := auth.NewService()
	spec, err := docs.MergeSpecs(s.cfg, "pkg/docs/base.json", "pkg/docs/modules")
	if err != nil {
		logger.Fatal("failed to merge openapi specs", zap.Error(err))
	}

	s.router.Get("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(spec)
	})

	s.router.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!doctype html>
		<html>
		  <head>
			<title>AI Interview API Documentation</title>
			<meta charset="utf-8" />
			<meta name="viewport" content="width=device-width, initial-scale=1" />
			<style> body { margin: 0; } </style>
		  </head>
		  <body>
			<script
			  id="api-reference"
			  data-url="/openapi.json"
			  data-configuration='{
				"theme": "purple",
				"layout": "modern",
				"darkMode": true,
				"searchHotKey": "k"
			  }'></script>
			<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
		  </body>
		</html>`

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, html)
	})

	s.router.Route("/", func(r chi.Router) {
		health.Init(r)
		books.Init(r, db.GetPool(), authService)
		auth.Init(r, authService)
	})
}

func (s *Server) Handler() http.Handler {
	return s.router
}
