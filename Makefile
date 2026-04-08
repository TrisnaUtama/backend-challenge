include .env
export

BINARY=tmp/app
PORT ?= 21001

.PHONY: all build run watch clean test itest migrate module docker-run docker-down

all: build test

# Build binary
build:
	@echo "Building..."
	@go build -o $(BINARY) ./cmd/api

# Run binary
run: build
	@echo "Running..."
	@./$(BINARY)

# Live reload using Air
watch:
	@if command -v air > /dev/null; then \
		air; \
	else \
		read -p "Air not installed. Install now? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/cosmtrek/air@latest; \
			air; \
		else \
			echo "Air not installed. Exiting..."; \
			exit 1; \
		fi; \
	fi

# Test unit
test:
	@echo "Running unit tests..."
	@go test ./... -v

# Integration test
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean build
clean:
	@echo "Cleaning..."
	@rm -rf $(BINARY) tmp

# Create migration file
migrate:
	@echo "Creating migration: $(filter-out $@,$(MAKECMDGOALS))"
	@goose -dir ./migrations create $(filter-out $@,$(MAKECMDGOALS)) sql

# Generate module
module:
	@echo "Creating module: $(filter-out $@,$(MAKECMDGOALS))"
	@go run ./cmd/scripts/generate-module.go $(filter-out $@,$(MAKECMDGOALS))

# Docker dev
docker-run:
	@docker compose up --build

docker-down:
	@docker compose down

migrate-up:
	@echo "Running migrations..."
	@goose -dir ./migrations postgres "$(DB_URL)" up

migrate-down:
	@echo "Rolling back..."
	@goose -dir ./migrations postgres "$(DB_URL)" down

migrate-status:
	@echo "Migration status..."
	@goose -dir ./migrations postgres "$(DB_URL)" status