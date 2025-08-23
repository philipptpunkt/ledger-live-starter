# Ledger Live Starter Makefile

.PHONY: build clean install dev test release help

# Variables
BINARY_NAME=ledger-live
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)

# Default target
help: ## Show this help message
	@echo "Ledger Live Starter Build Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@go build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME) ./cmd/ledger-live
	@echo "✓ Built $(BINARY_NAME)"

dev: ## Build and run for development
	@go run ./cmd/ledger-live start

test: ## Run tests
	@go test -v ./...

clean: ## Clean build artifacts
	@rm -f $(BINARY_NAME)
	@echo "✓ Cleaned build artifacts"

install: build ## Build and install to ~/.ledger-live/
	@echo "Installing to ~/.ledger-live/..."
	@mkdir -p ~/.ledger-live
	@cp $(BINARY_NAME) ~/.ledger-live/
	@chmod +x ~/.ledger-live/$(BINARY_NAME)
	@echo "✓ Installed to ~/.ledger-live/"
	@echo "Add ~/.ledger-live to your PATH to use globally"

# Cross-compilation targets
build-all: ## Build for all supported platforms
	@echo "Building for all platforms..."
	@mkdir -p dist
	@GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dist/$(BINARY_NAME)-linux-amd64 ./cmd/ledger-live
	@GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o dist/$(BINARY_NAME)-linux-arm64 ./cmd/ledger-live
	@GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dist/$(BINARY_NAME)-darwin-amd64 ./cmd/ledger-live
	@GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o dist/$(BINARY_NAME)-darwin-arm64 ./cmd/ledger-live
	@GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o dist/$(BINARY_NAME)-windows-amd64.exe ./cmd/ledger-live
	@echo "✓ Built all platform binaries in dist/"

setup-dev: ## Set up development environment
	@echo "Setting up development environment..."
	@go mod download
	@go mod tidy
	@echo "✓ Development environment ready"

release: clean build-all ## Prepare release builds
	@echo "✓ Release builds ready in dist/"
	@ls -la dist/
