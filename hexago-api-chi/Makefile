
.PHONY: build clean test run help

SERVICE_NAME=hexago-api-chi

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building $(SERVICE_NAME)..."
	@CGO_ENABLED=0 go build -o $(SERVICE_NAME) main.go

run: ## Run the application
	@echo "Running $(SERVICE_NAME)..."
	@go run main.go run

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f $(SERVICE_NAME) coverage.out

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	@go tool cover -html=coverage.out

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...

mod-tidy: ## Tidy go modules
	@echo "Tidying modules..."
	@go mod tidy

.DEFAULT_GOAL := help
