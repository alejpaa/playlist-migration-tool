# Variables
APP_NAME := playlist-migration-tool
API_NAME := playlist-api
CMD_DIR := cmd/api
BUILD_DIR := build
BIN_DIR := bin
MAIN_PATH := ./$(CMD_DIR)

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := $(GOCMD) fmt

# Build flags
BUILD_FLAGS := -v
LDFLAGS := -w -s

.PHONY: all build clean test deps fmt vet run dev install help

# Default target
all: clean deps test build

# Build the application
build:
	@echo "Building $(API_NAME)..."
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(API_NAME) $(MAIN_PATH)

# Build for multiple platforms
build-all: clean deps
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	# Linux
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(API_NAME)-linux-amd64 $(MAIN_PATH)
	# Windows
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(API_NAME)-windows-amd64.exe $(MAIN_PATH)
	# macOS
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(API_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(API_NAME)-darwin-arm64 $(MAIN_PATH)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -rf $(BIN_DIR)

# OAuth2 authentication
.PHONY: auth
auth: build ## Authenticate with YouTube API
	@echo "🔐 Starting YouTube OAuth2 authentication..."
	@./$(BIN_DIR)/$(APP_NAME)

# Clean authentication tokens
.PHONY: clean-auth
clean-auth: ## Remove saved authentication tokens
	@echo "🧹 Cleaning authentication tokens..."
	@rm -f token.json
	@echo "✅ Tokens cleaned. You'll need to re-authenticate next time."

# Run advanced example
.PHONY: example-advanced
example-advanced: build ## Run advanced playlist example
	@echo "� Running advanced example..."
	@go run examples/advanced_example.go

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...

# Run the application
run:
	@echo "Running $(API_NAME)..."
	$(GOCMD) run $(MAIN_PATH)

# Development mode with hot reload (requires air)
dev:
	@echo "Starting development server..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Installing air..."; \
		$(GOGET) -u github.com/cosmtrek/air; \
		air; \
	fi

# Install the application
install: build
	@echo "Installing $(APP_NAME)..."
	@cp $(BIN_DIR)/$(APP_NAME) $(GOPATH)/bin/

# Lint code (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install it with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2"; \
	fi

# Security scan (requires gosec)
security:
	@echo "Running security scan..."
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Install it with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Show help
help:
	@echo "Available commands:"
	@echo "  build         Build the application"
	@echo "  build-all     Build for multiple platforms"
	@echo "  clean         Clean build artifacts"
	@echo "  test          Run tests"
	@echo "  test-coverage Run tests with coverage report"
	@echo "  deps          Download and tidy dependencies"
	@echo "  fmt           Format code"
	@echo "  vet           Run go vet"
	@echo "  run           Run the application"
	@echo "  dev           Start development server with hot reload"
	@echo "  install       Install the application to GOPATH/bin"
	@echo "  lint          Run linter (requires golangci-lint)"
	@echo "  security      Run security scan (requires gosec)"
	@echo "  help          Show this help message"
