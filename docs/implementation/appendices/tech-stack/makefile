# Makefile for Contrail

.PHONY: build test lint clean install

# Build variables
BINARY_NAME := contrail
BUILD_DIR := ./bin
MAIN_PATH := ./cmd/contrail

# Go variables
GOFLAGS := -ldflags="-s -w"

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Install to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(GOFLAGS) $(MAIN_PATH)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	go test -v -tags=integration ./...

# Run end-to-end tests
test-e2e:
	@echo "Running e2e tests..."
	go test -v -tags=e2e ./...

# Run all tests
test-all: test test-integration test-e2e

# Lint the code
lint:
	@echo "Linting..."
	golangci-lint run

# Format code
fmt:
	@echo "Formatting..."
	go fmt ./...

# Tidy modules
tidy:
	@echo "Tidying modules..."
	go mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Generate mocks
generate:
	@echo "Generating mocks..."
	go generate ./...

# Build for all platforms
build-all:
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

# Release (requires goreleaser)
release:
	@echo "Creating release..."
	goreleaser release --clean

# Snapshot release (for testing)
release-snapshot:
	@echo "Creating snapshot release..."
	goreleaser release --snapshot --clean

# Run the application
run:
	@go run $(MAIN_PATH) $(ARGS)

# Check for vulnerabilities
vuln:
	@echo "Checking for vulnerabilities..."
	govulncheck ./...

# Help
help:
	@echo "Available targets:"
	@echo "  build            - Build the binary"
	@echo "  install          - Install to GOPATH/bin"
	@echo "  test             - Run unit tests"
	@echo "  test-coverage    - Run tests with coverage"
	@echo "  test-integration - Run integration tests"
	@echo "  test-e2e         - Run end-to-end tests"
	@echo "  test-all         - Run all tests"
	@echo "  lint             - Run linter"
	@echo "  fmt              - Format code"
	@echo "  tidy             - Tidy modules"
	@echo "  clean            - Clean build artifacts"
	@echo "  generate         - Generate mocks"
	@echo "  build-all        - Build for all platforms"
	@echo "  release          - Create release"
	@echo "  release-snapshot - Create snapshot release"
	@echo "  run ARGS=...     - Run the application"
	@echo "  vuln             - Check for vulnerabilities"
	@echo "  help             - Show this help"
