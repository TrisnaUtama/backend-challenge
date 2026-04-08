package database

import (
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func (db *PostgresDB) RunMigrations(migrationDir string) error {
	sqlDb := stdlib.OpenDBFromPool(db.Pool)
	defer sqlDb.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose set dialect: %w", err)
	}

	if err := goose.Up(sqlDb, migrationDir); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	log.Println(" Migrations applied successfully")
	return nil
}
