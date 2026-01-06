# Makefile
# Build and development tasks for Contrail

.PHONY: all build test lint clean install generate help

# Default target
all: build

# Build the binary
build:
	go build -o bin/contrail ./cmd/contrail

# Build with version information
build-release:
	go build -ldflags "-s -w -X main.version=$(VERSION)" -o bin/contrail ./cmd/contrail

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	golangci-lint run ./...

# Format code
fmt:
	go fmt ./...
	gofumpt -w .

# Generate code (shell scripts, etc.)
generate:
	go generate ./...

# Install the binary
install:
	go install ./cmd/contrail

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install development dependencies
dev-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest

# Run the application
run:
	go run ./cmd/contrail

# Build for all platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/contrail-linux-amd64 ./cmd/contrail
	GOOS=linux GOARCH=arm64 go build -o bin/contrail-linux-arm64 ./cmd/contrail
	GOOS=darwin GOARCH=amd64 go build -o bin/contrail-darwin-amd64 ./cmd/contrail
	GOOS=darwin GOARCH=arm64 go build -o bin/contrail-darwin-arm64 ./cmd/contrail
	GOOS=windows GOARCH=amd64 go build -o bin/contrail-windows-amd64.exe ./cmd/contrail

# Run integration tests
test-integration:
	go test -v -tags=integration ./...

# Run end-to-end tests
test-e2e:
	go test -v -tags=e2e ./...

# Run all tests
test-all: test test-integration test-e2e

# Tidy modules
tidy:
	go mod tidy

# Check for vulnerabilities
vuln:
	govulncheck ./...

# Help
help:
	@echo "Available targets:"
	@echo "  build            - Build the binary"
	@echo "  build-release    - Build with version information"
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
	@echo "  generate         - Generate code"
	@echo "  dev-deps         - Install development dependencies"
	@echo "  build-all        - Build for all platforms"
	@echo "  run              - Run the application"
	@echo "  vuln             - Check for vulnerabilities"
	@echo "  help             - Show this help"
