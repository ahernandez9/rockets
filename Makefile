.PHONY: help build run clean lint swagger test install-tools generate-mocks

help: ## Display this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

install-tools: ## Install required development tools (optional - tools run via 'go run')
	@echo "Note: Tools are run via 'go run' and will be downloaded automatically when needed"
	@echo "If you prefer to install them globally:"
	@echo "  go install github.com/swaggo/swag/cmd/swag@latest"
	@echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
	@echo "  go install go.uber.org/mock/mockgen@latest  # Note: We use go.uber.org/mock (not github.com/golang/mock)"

generate-mocks: ## Generate mock implementations for all interfaces
	@echo "Generating mocks..."
	@go generate ./...
	@echo "Mocks generated successfully!"

swagger: ## Generate swagger documentation
	@echo "Generating swagger docs..."
	@go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go -o docs
	@echo "Swagger docs generated successfully!"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...
	@echo "Linting completed!"

build: swagger ## Build the application
	@echo "Building application..."
	@go build -o bin/rockets cmd/server/main.go
	@echo "Build completed! Binary: bin/rockets"

run: ## Run the application
	@go run cmd/server/main.go

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf docs/
	@echo "Clean completed!"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated!"

