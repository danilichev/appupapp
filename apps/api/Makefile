all: build test

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

clean:
	@echo "Cleaning..."
	@rm -f main

docker-down:
	@if command -v docker-compose >/dev/null; then \
		docker-compose down; \
	elif command -v docker >/dev/null && docker compose version >/dev/null 2>&1; then \
		docker compose down; \
	else \
		echo "Neither Docker Compose V1 nor V2 is available."; \
		exit 1; \
	fi

docker-up:
	@if command -v docker-compose >/dev/null; then \
		docker-compose up --build; \
	elif command -v docker >/dev/null && docker compose version >/dev/null 2>&1; then \
		docker compose up --build; \
	else \
		echo "Neither Docker Compose V1 nor V2 is available."; \
		exit 1; \
	fi

run:
	@go run cmd/api/main.go

test:
	@echo "Testing..."
	@go test ./... -v

test-integration:
	@echo "Running integration tests..."
	@go test ./internal/database -v

watch:
	@echo "Starting development server with Air..."
	@if command -v air >/dev/null; then \
	    air; \
    else \
        read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
        if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
            go install github.com/air-verse/air@latest; \
            air; \
        else \
            echo "You chose not to install air. Exiting..."; \
            exit 1; \
        fi; \
    fi

.PHONY: all build clean docker-down docker-up run test test-integration watch