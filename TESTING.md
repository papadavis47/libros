# Libros Testing Documentation

This document provides an overview of the comprehensive test suite for the Libros book management application.

## Test Structure

The test suite follows Go best practices with tests co-located alongside the code they test:

### Package-Level Tests
Tests are organized by package and located within each respective package directory:

- **`internal/constants/constants_test.go`** - Tests application constants and configuration values
- **`internal/database/database_test.go`** - Tests database operations with real SQLite (integration-style)
- **`internal/database/main_database_test.go`** - Tests full CRUD cycles and validation
- **`internal/factory/factory_test.go`** - Tests UI component factory functions
- **`internal/models/models_test.go`** - Tests core data models (Book, BookType, Screen constants)
- **`internal/services/services_test.go`** - Tests backup and export services with file I/O
- **`internal/ui/ui_test.go`** - Tests UI model initialization and Bubble Tea integration
- **`internal/utils/utils_test.go`** - Tests utility functions like date and book type formatting
- **`internal/validation/validation_test.go`** - Tests input validation functions for data integrity

### Test Package Types

**Same Package Tests** (e.g., `package constants`):
- Test internal/private functions and methods
- Access package-internal variables and functions directly
- Used for most unit tests

**Separate Package Tests** (e.g., `package database_test`):
- Test public API from external perspective (black-box testing)
- Used for integration-style tests that verify public interfaces

## Running Tests

### Quick Commands

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage report
go test -cover ./...

# Run specific package tests
go test ./internal/database/...
go test ./internal/models/...
```

### Package-Specific Testing

```bash
# Run constants tests
go test ./internal/constants/

# Run database tests
go test ./internal/database/

# Run validation tests
go test ./internal/validation/

# Run all tests with coverage
go test -cover ./...
```

## Test Coverage

### Constants (`internal/constants/constants_test.go`)
- UI dimension constants (input widths, text area sizes)
- Character limits for different fields
- File permission constants
- Path resolution functions (`GetAppDir`, `GetDatabasePath`)
- Logical relationships between constants
- Immutability verification

### Database (`internal/database/`)
**database_test.go** - Integration-style database tests:
- CRUD operations (Create, Read, Update, Delete) with real SQLite
- Data persistence and retrieval
- Special character handling (unicode, emojis, quotes)
- Edge cases (nonexistent records, empty data)
- Book counting functionality

**main_database_test.go** - Full workflow tests:
- Complete CRUD cycles with validation
- Book validation during save operations
- Database initialization and cleanup

### Factory (`internal/factory/factory_test.go`)
- Text input creation with consistent styling
- Textarea creation for notes
- Path input creation for export screens
- Component consistency across factory functions
- Focus state management
- Performance benchmarks

### Models (`internal/models/models_test.go`)
- Book model field validation and access
- BookType enum values and string representations
- Screen navigation constants verification
- Data structure integrity

### Services (`internal/services/services_test.go`)
- JSON export with proper formatting and metadata
- Markdown export with headers, formatting, and separators
- Database backup file operations
- File I/O error handling
- Empty data set handling

### UI (`internal/ui/ui_test.go`)
- Bubble Tea model initialization
- UI component creation and setup
- Integration with database layer

### Utils (`internal/utils/utils_test.go`)
- Date formatting with ordinal suffixes (1st, 2nd, 3rd, etc.)
- Book type formatting for both enum and string inputs
- Edge cases with special characters and unicode
- Performance benchmarks for formatting functions

### Validation (`internal/validation/validation_test.go`)
- Book validation with all field combinations
- Title validation (required, length limits, trimming)
- Author validation (required, length limits, trimming)
- Notes validation (optional, length limits)
- File path validation for export operations
- Input trimming and sanitization

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
- **Co-located**: Tests are in the same package as the code they test, following Go conventions

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

1. **Add Tests in Same Package** - Create `*_test.go` files alongside your code
2. **Use Same Package Name** for unit tests testing internal functions
3. **Use `_test` Package Suffix** for integration tests testing public APIs
4. **Follow Naming Conventions**: Use descriptive test function names starting with `Test`
5. **Include Documentation**: Add comments explaining what the test verifies
6. **Test Edge Cases**: Include tests for boundary conditions and error cases
7. **Update This Documentation** if you add new test categories or patterns

### Test Organization Examples

**For internal/mypackage/mycode.go:**
```go
// internal/mypackage/mycode_test.go
package mypackage // Same package for unit tests

func TestInternalFunction(t *testing.T) {
    // Test internal functions directly
}
```

**For integration testing:**
```go
// internal/mypackage/integration_test.go  
package mypackage_test // Separate package for integration tests

func TestPublicAPI(t *testing.T) {
    // Test public API from external perspective
}
```

## Performance Benchmarks

The test suite includes benchmarks for performance-critical functions:
- Date formatting (used frequently in UI)
- Book type formatting (used in list displays)
- Factory function creation (used during screen initialization)

Run benchmarks with:
```bash
go test -bench=. ./...
```

## Test Results Interpretation

- **PASS**: All tests completed successfully
- **FAIL**: One or more tests failed - check output for specific failures
- **Build Failed**: Compilation errors - fix syntax/import issues first
- **Timeout**: Tests took too long - may indicate performance issues or infinite loops

The test suite follows Go best practices and is designed to be reliable and deterministic, providing confidence in the application's correctness and stability.