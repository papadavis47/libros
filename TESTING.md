# Libros Testing Documentation

This document provides an overview of the comprehensive test suite for the Libros book management application.

## Test Structure

The test suite is organized into two main categories:

### Unit Tests (`tests/unit/`)
Unit tests focus on testing individual components in isolation:

- **`models_test.go`** - Tests the core data models (Book, BookType, Screen constants)
- **`utils_test.go`** - Tests utility functions like date formatting and book type formatting
- **`validation_test.go`** - Tests input validation functions for forms and data integrity
- **`constants_test.go`** - Tests application constants and configuration values
- **`factory_test.go`** - Tests UI component factory functions for consistent input creation

### Integration Tests (`tests/integration/`)
Integration tests verify that components work correctly together:

- **`database_simple_test.go`** - Tests database operations with real SQLite
- **`services_simple_test.go`** - Tests backup and export services with file I/O

## Running Tests

### Quick Commands

```bash
# Run all tests
make test

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration

# Run tests with verbose output
make test-verbose

# Run tests with coverage report
make test-coverage
```

### Manual Commands

```bash
# Run all tests
go test ./tests/unit/... ./tests/integration/...

# Run unit tests only
go test ./tests/unit/...

# Run integration tests only
go test ./tests/integration/...

# Run with verbose output
go test -v ./tests/unit/... ./tests/integration/...

# Run with coverage
go test -cover ./tests/unit/... ./tests/integration/...
```

## Test Coverage

### Unit Tests Cover:

#### Models (`models_test.go`)
- Book model field validation and access
- BookType enum values and string representations
- Screen navigation constants verification
- Data structure integrity

#### Utils (`utils_test.go`)
- Date formatting with ordinal suffixes (1st, 2nd, 3rd, etc.)
- Book type formatting for both enum and string inputs
- Edge cases with special characters and unicode
- Performance benchmarks for formatting functions

#### Validation (`validation_test.go`)
- Book validation with all field combinations
- Title validation (required, length limits, trimming)
- Author validation (required, length limits, trimming)
- Notes validation (optional, length limits)
- File path validation for export operations
- Input trimming and sanitization

#### Constants (`constants_test.go`)
- UI dimension constants (input widths, text area sizes)
- Character limits for different fields
- File permission constants
- Path resolution functions
- Logical relationships between constants

#### Factory (`factory_test.go`)
- Text input creation with consistent styling
- Textarea creation for notes
- Path input creation for export screens
- Component consistency across factory functions
- Focus state management

### Integration Tests Cover:

#### Database (`database_simple_test.go`)
- CRUD operations (Create, Read, Update, Delete) with real SQLite
- Data persistence and retrieval
- Special character handling (unicode, emojis, quotes)
- Edge cases (nonexistent records, empty data)
- Book counting functionality

#### Services (`services_simple_test.go`)
- JSON export with proper formatting and metadata
- Markdown export with headers, formatting, and separators
- Database backup file operations
- File I/O error handling
- Empty data set handling

## Test Quality Features

### Comprehensive Coverage
- **Edge Cases**: Tests handle empty inputs, special characters, unicode, and boundary conditions
- **Error Conditions**: Tests verify proper error handling for invalid inputs and operations
- **Data Integrity**: Tests ensure data is preserved correctly through all operations
- **Performance**: Benchmarks verify acceptable performance for UI operations

### Real Environment Testing
- **Temporary Files**: Integration tests use temporary directories and files
- **Real Database**: Tests use actual SQLite databases, not mocks
- **File Operations**: Tests perform real file I/O for export and backup functions

### Maintainability
- **Clear Documentation**: Each test function includes comments explaining what it tests and why
- **Descriptive Names**: Test names clearly indicate what functionality is being verified
- **Structured Organization**: Tests are logically grouped by component and functionality
- **Isolated Tests**: Each test is independent and can run in any order

## Test Dependencies

The test suite requires:
- Go 1.19+ (for testing framework features)
- SQLite3 (for database integration tests)
- Temporary file system access (for integration tests)

## Continuous Integration

The test suite is designed to run in CI environments:
- No external dependencies beyond Go standard library and project dependencies
- Temporary file cleanup ensures no test artifacts remain
- Reasonable execution time (typically under 30 seconds)
- Clear pass/fail status for automation

## Adding New Tests

When adding new functionality:

1. **Add Unit Tests** for new functions in the appropriate `*_test.go` file
2. **Add Integration Tests** if the functionality involves database or file operations
3. **Follow Naming Conventions**: Use descriptive test function names starting with `Test`
4. **Include Documentation**: Add comments explaining what the test verifies
5. **Test Edge Cases**: Include tests for boundary conditions and error cases
6. **Update This Documentation** if you add new test files or major test categories

## Performance Benchmarks

The test suite includes benchmarks for performance-critical functions:
- Date formatting (used frequently in UI)
- Book type formatting (used in list displays)
- Factory function creation (used during screen initialization)

Run benchmarks with:
```bash
go test -bench=. ./tests/unit/...
```

## Test Results Interpretation

- **PASS**: All tests completed successfully
- **FAIL**: One or more tests failed - check output for specific failures
- **Build Failed**: Compilation errors - fix syntax/import issues first
- **Timeout**: Tests took too long - may indicate performance issues or infinite loops

The test suite is designed to be reliable and deterministic, providing confidence in the application's correctness and stability.