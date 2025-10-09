# Makefile for role-selector
# Copyright 2002-2025 Barcelona Supercomputing Center (www.bsc.es)

# Variables
BINARY_NAME=roleselectorexec
DOCKER_IMAGE=role-selector
DOCKER_TAG=latest
GO_VERSION=1.24.1

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Build flags
LDFLAGS=-ldflags "-s -w"
BUILD_FLAGS=-a -installsuffix cgo $(LDFLAGS)

.PHONY: all build clean test coverage fmt vet mod-tidy docker-build docker-run docker-push help install run

# Default target
all: clean fmt vet mod-tidy build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME) .
	@echo "Build completed: $(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	@echo "Clean completed"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# Tidy go modules
mod-tidy:
	@echo "Tidying go modules..."
	$(GOMOD) tidy

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download

# Install the binary to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	$(GOCMD) install .
	@echo "Installation completed"

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

# Docker targets
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)"

docker-push: docker-build
	@echo "Pushing Docker image..."
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)

# Lint and format check
lint-check:
	@echo "Checking code formatting..."
	@if [ "$$($(GOFMT) . | wc -l)" -ne 0 ]; then \
		echo "Code is not formatted. Run 'make fmt' to fix."; \
		$(GOFMT) -l .; \
		exit 1; \
	fi
	@echo "Code formatting check passed"

# Help target
help:
	@echo "Available targets:"
	@echo "  all          - Clean, format, vet, tidy, and build"
	@echo "  build        - Build the binary"
	@echo "  build-linux  - Build for Linux (Docker compatible)"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  coverage     - Run tests with coverage report"
	@echo "  fmt          - Format code"
	@echo "  vet          - Run go vet"
	@echo "  mod-tidy     - Tidy go modules"
	@echo "  deps         - Download dependencies"
	@echo "  install      - Install binary to GOPATH/bin"
	@echo "  run          - Build and run the application"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Build and run Docker container"
	@echo "  docker-push  - Build and push Docker image"
	@echo "  dev-setup    - Set up development environment"
	@echo "  lint-check   - Check code formatting"
	@echo "  help         - Show this help message"
