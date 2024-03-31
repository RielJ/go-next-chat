# Simple Makefile for a Go project
include app.env

$(eval export $(shell sed -ne 's/ *#.*$$//; /./ s/=.*$$// p' app.env))

# Build the application
all: build

build:
	@echo "Building..."
	
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# migrate up
migrateup:
	@echo "Migrating up..."
	migrate -path internal/db/migration -database "postgres://$$DB_USERNAME:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_DATABASE?sslmode=disable" --verbose up

# migrate up 1
migrateup1:
	@echo "Migrating up..."
	migrate -path internal/db/migration -database "postgres://$$DB_USERNAME:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_DATABASE?sslmode=disable" --verbose up 1

# migrate down
migratedown:
	@echo "Migrating up..."
	migrate -path internal/db/migration -database "postgres://$$DB_USERNAME:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_DATABASE?sslmode=disable" --verbose down

# migrate down
migratedown1:
	@echo "Migrating up..."
	migrate -path internal/db/migration -database "postgres://$$DB_USERNAME:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_DATABASE?sslmode=disable" --verbose down 1

# new migration
new_migration:
	migrate create -ext sql -dir internal/db/migration -seq $(name)

# sqlc generate
sqlc:
	@echo "Generating SQLC..."
	sqlc generate


# Run Docker Detached
ddocker-run:
	@if docker compose up -d 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up -d; \
	fi

# Create DB container
docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload Client
watch-client:
	@cd client && npm run dev

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

.PHONY: all build run test clean
