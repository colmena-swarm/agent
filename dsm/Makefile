# DSM (Dynamic Service Manager) Makefile
# Copyright 2002-2025 Barcelona Supercomputing Center (www.bsc.es)

# Variables
BINARY_NAME := dsm
BUILD_DIR := build
UNIT_TESTS_DIR := $(shell go list ./... | grep -v '/tests/integration$$')
INTEGRATION_TESTS_DIR := ./tests/integration

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) .

# Run tests (unit + integration)
.PHONY: test unit integration
 test: unit integration

unit:
	@echo "Running unit tests..."
	@if [ -n "$(UNIT_TESTS_DIR)" ]; then go test -v $(UNIT_TESTS_DIR); else echo "No unit test packages"; fi

integration:
	@echo "Running integration tests..."
	@go test -v -timeout 10m $(INTEGRATION_TESTS_DIR)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  test         - Run unit and integration tests"
	@echo "  unit         - Run unit tests (excludes tests/integration)"
	@echo "  integration  - Run integration tests only"
	@echo "  clean        - Clean build artifacts"
	@echo "  help         - Show this help message"
