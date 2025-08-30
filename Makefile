# RUMI Backend Makefile

# Variables
BINARY_NAME=server
BUILD_DIR=bin
CMD_DIR=cmd/server
MIGRATION_DIR=migrations

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: help build run clean test deps migrate migrate-down dev fmt lint vet

# Default target
help: ## Show help message
	@echo "$(BLUE)RUMI Backend - Available commands:$(NC)"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

deps: ## Install dependencies
	@echo "$(BLUE)Installing dependencies...$(NC)"
	@go mod tidy
	@go mod download
	@echo "$(GREEN)✅ Dependencies installed$(NC)"

build: ## Build the application
	@echo "$(BLUE)Building application...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go
	@echo "$(GREEN)✅ Build completed: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

run: ## Run the application in development mode
	@echo "$(BLUE)Starting server in development mode...$(NC)"
	@go run $(CMD_DIR)/main.go

dev: ## Run with air for hot reload (requires air to be installed)
	@echo "$(BLUE)Starting server with hot reload...$(NC)"
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(YELLOW)Air not found. Install it with: go install github.com/cosmtrek/air@latest$(NC)"; \
		echo "$(BLUE)Running without hot reload...$(NC)"; \
		go run $(CMD_DIR)/main.go; \
	fi

migrate: ## Run database migrations
	@echo "$(BLUE)Running database migrations...$(NC)"
	@go run scripts/migrate.go
	@echo "$(GREEN)✅ Migrations completed$(NC)"

migrate-down: ## Rollback database migrations (manual process)
	@echo "$(YELLOW)⚠️  Manual rollback required. Check migration files in $(MIGRATION_DIR)$(NC)"

test: ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	@go test -v ./...
	@echo "$(GREEN)✅ Tests completed$(NC)"

test-coverage: ## Run tests with coverage
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Coverage report generated: coverage.html$(NC)"

fmt: ## Format code
	@echo "$(BLUE)Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✅ Code formatted$(NC)"

lint: ## Run linter (requires golangci-lint)
	@echo "$(BLUE)Running linter...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
		echo "$(GREEN)✅ Linting completed$(NC)"; \
	else \
		echo "$(YELLOW)golangci-lint not found. Install it from: https://golangci-lint.run/usage/install/$(NC)"; \
	fi

vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✅ Vet completed$(NC)"

clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)✅ Cleaned$(NC)"

install-tools: ## Install development tools
	@echo "$(BLUE)Installing development tools...$(NC)"
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)✅ Tools installed$(NC)"

docker-build: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(NC)"
	@docker build -t rumi-backend .
	@echo "$(GREEN)✅ Docker image built$(NC)"

docker-run: ## Run Docker container
	@echo "$(BLUE)Running Docker container...$(NC)"
	@docker run -p 8080:8080 --env-file .env rumi-backend

setup: deps migrate ## Complete setup (install deps and run migrations)
	@echo "$(GREEN)✅ Setup completed! You can now run 'make run' to start the server$(NC)"

# Continuous Integration targets
ci-test: fmt vet test ## Run CI tests (format, vet, test)
	@echo "$(GREEN)✅ All CI tests passed$(NC)"

# Development workflow
quick-start: deps build run ## Quick start for new developers

# Production build
prod-build: ## Build for production
	@echo "$(BLUE)Building for production...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go
	@echo "$(GREEN)✅ Production build completed$(NC)"
