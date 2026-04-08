package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"backend-challenge/internal/configs"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	Pool *pgxpool.Pool
}

func NewPostgresConn(cfg *configs.Setting) (*PostgresDB, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.Database.ConnStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %v", err)
	}

	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 30 * time.Minute
	poolConfig.MaxConnIdleTime = 5 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %v", err)
	}

	log.Println("Database connected successfully to:", cfg.Database.Name)
	return &PostgresDB{Pool: pool}, nil
}

func (db *PostgresDB) HealthCheck(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

func (db *PostgresDB) GetPool() *pgxpool.Pool {
	return db.Pool
}

func (db *PostgresDB) Close() {
	if db.Pool != nil {
		log.Println("🔌 Closing database connection pool...")
		db.Pool.Close()
	}
}
