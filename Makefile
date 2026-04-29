.PHONY: build test clean install fmt vet lint

# Build variables
BINARY_NAME=linkmgr
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILD_DATE)"

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/linkmgr

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	go test -v -race -tags=integration ./tests/integration/...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	go clean

# Install to GOBIN or ~/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	@if [ -n "$(GOBIN)" ]; then \
		cp $(BINARY_NAME) $(GOBIN)/; \
	else \
		mkdir -p ~/bin; \
		cp $(BINARY_NAME) ~/bin/; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run linter (if golangci-lint is available)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, skipping..."; \
	fi

# Run all checks
check: fmt vet test

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 ./cmd/linkmgr
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 ./cmd/linkmgr
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 ./cmd/linkmgr
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe ./cmd/linkmgr

# Show help
help:
	@echo "Available targets:"
	@echo "  build          - Build the binary"
	@echo "  test           - Run unit tests"
	@echo "  test-integration - Run integration tests"
	@echo "  clean          - Remove build artifacts"
	@echo "  install        - Install binary to GOBIN or ~/bin"
	@echo "  fmt            - Format code"
	@echo "  vet            - Run go vet"
	@echo "  lint           - Run golangci-lint"
	@echo "  check          - Run fmt, vet, and test"
	@echo "  build-all      - Build for multiple platforms"
