package main

import (
	"fmt"
	"log"
	"net/http"

	"backend-challenge/internal/configs"
	postgres "backend-challenge/internal/database"
	"backend-challenge/internal/server"
	"backend-challenge/pkg/logger"
)

func main() {
	cfg, err := configs.NewSetting()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger.Init(cfg.App.Env)
	defer logger.Log.Sync()

	db, err := postgres.NewPostgresConn(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	srv := server.New(cfg, db)

	serverAddr := fmt.Sprintf(":%d", cfg.App.Port)
	log.Printf("Server is running on http://localhost%s", serverAddr)
	log.Printf("Docs available on http://localhost%s/docs", serverAddr)
	log.Printf("Environment: %s", cfg.App.Env)

	if err = http.ListenAndServe(serverAddr, srv.Handler()); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
