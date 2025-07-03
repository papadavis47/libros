# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Building and Running
- `make build` - Build the libros binary
- `make clean` - Remove build artifacts
- `./libros` - Run the built application

### Testing
- `make test` - Run all tests (unit + integration)
- `make test-unit` - Run only unit tests
- `make test-integration` - Run only integration tests  
- `make test-verbose` - Run tests with verbose output
- `make test-coverage` - Run tests with coverage report
- `go test ./internal/[package]/` - Run specific package tests

### Manual Testing
- `go test -v ./internal/database/` - Test database operations
- `go test -bench=. ./...` - Run performance benchmarks

## Architecture Overview

### Core Structure
- **Bubble Tea TUI**: Terminal user interface using Charmbracelet's Bubble Tea framework
- **SQLite Database**: Local database stored at `~/.libros/books.db`
- **Clean Architecture**: Separated concerns with interfaces, models, and services
- **Screen-Based Navigation**: Each UI screen is a separate model with its own state

### Key Components
- **Models**: Core data structures (Book, BookType, Screen constants)
- **Interfaces**: Repository and BackupService abstractions for testability
- **Database**: SQLite operations with CRUD functionality
- **UI Screens**: Individual screen models for each application view
- **Services**: Business logic for backup, export (JSON/Markdown)
- **Factory**: Consistent UI component creation (text inputs, textareas)

### Navigation Flow
Application uses a main `ui.Model` that coordinates between screen models:
- MenuScreen → AddBookScreen/ListBooksScreen/UtilitiesScreen
- ListBooksScreen → BookDetailScreen → EditBookScreen
- UtilitiesScreen → ExportScreen/BackupScreen

### Database Schema
- Books table with fields: ID, Title, Author, Type, Notes, CreatedAt, UpdatedAt
- BookType enum: paperback, hardback, audio, digital
- Database path: `~/.libros/books.db`

## Development Patterns

### Testing Approach
- Tests co-located with source code in same packages
- Integration tests use real SQLite databases (not mocks)
- Same-package tests for internal functions, separate `_test` packages for public API
- Comprehensive coverage including edge cases, unicode, and performance benchmarks

### Code Style
- Extensive inline documentation for all public interfaces
- Clear separation of concerns between UI, business logic, and data access
- Error handling with proper cleanup (database connections, file operations)
- Consistent naming conventions following Go standards

### UI Component Creation
- Use `internal/factory` for consistent UI component creation
- All text inputs and textareas created through factory functions
- Shared styling through `internal/styles` package

### State Management
- Each screen maintains its own state and handles its own updates
- Main model coordinates screen transitions and shared state
- Database connection shared across all screens
- Proper cleanup on screen transitions and application exit

## Key Files
- `cmd/libros/main.go` - Application entry point
- `internal/ui/model.go` - Main Bubble Tea model coordinating all screens
- `internal/database/database.go` - Database operations and connection management
- `internal/models/book.go` - Core data structures and constants
- `internal/interfaces/` - Repository and service interfaces
- `internal/ui/screens/` - Individual screen implementations