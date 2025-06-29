# Makefile for Libros Book Manager
# Provides convenient commands for testing, building, and development

.PHONY: test test-unit test-integration test-all test-verbose test-coverage build clean help

# Default target
help:
	@echo "Libros Book Manager - Available Commands:"
	@echo ""
	@echo "  test          - Run all tests (unit + integration)"
	@echo "  test-unit     - Run only unit tests"
	@echo "  test-integration - Run only integration tests"
	@echo "  test-verbose  - Run all tests with verbose output"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  build         - Build the application"
	@echo "  clean         - Clean build artifacts"
	@echo "  help          - Show this help message"

# Run all tests
test:
	@echo "🧪 Running all tests..."
	@go test ./tests/unit/... ./tests/integration/...

# Run only unit tests
test-unit:
	@echo "🧪 Running unit tests..."
	@go test ./tests/unit/...

# Run only integration tests  
test-integration:
	@echo "🔗 Running integration tests..."
	@go test ./tests/integration/...

# Run all tests with verbose output
test-verbose:
	@echo "🧪 Running all tests (verbose)..."
	@go test -v ./tests/unit/... ./tests/integration/...

# Run tests with coverage report
test-coverage:
	@echo "📊 Running tests with coverage..."
	@go test -cover ./tests/unit/... ./tests/integration/...
	@echo ""
	@echo "📊 Generating detailed coverage report..."
	@go test -coverprofile=coverage.out ./tests/unit/... ./tests/integration/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build the application
build:
	@echo "🔨 Building Libros..."
	@go build -o libros ./cmd/libros

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -f libros coverage.out coverage.html
	@go clean