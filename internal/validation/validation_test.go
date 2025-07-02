package validation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/papadavis47/libros/internal/models"
)

// TestValidateBook tests comprehensive book validation
// This function is critical for ensuring data integrity before database operations
func TestValidateBook(t *testing.T) {
	tests := []struct {
		name          string
		book          *models.Book
		expectedCount int // number of expected validation errors
		shouldFail    bool
	}{
		{
			name: "valid book with all fields",
			book: &models.Book{
				Title:  "Valid Title",
				Author: "Valid Author",
				Type:   models.Paperback,
				Notes:  "Some notes about the book",
			},
			expectedCount: 0,
			shouldFail:    false,
		},
		{
			name: "valid book without notes",
			book: &models.Book{
				Title:  "Another Valid Title",
				Author: "Another Valid Author",
				Type:   models.Digital,
				Notes:  "",
			},
			expectedCount: 0,
			shouldFail:    false,
		},
		{
			name: "book with empty title",
			book: &models.Book{
				Title:  "",
				Author: "Valid Author",
				Type:   models.Hardback,
				Notes:  "Notes",
			},
			expectedCount: 1,
			shouldFail:    true,
		},
		{
			name: "book with empty author",
			book: &models.Book{
				Title:  "Valid Title",
				Author: "",
				Type:   models.Audio,
				Notes:  "Notes",
			},
			expectedCount: 1,
			shouldFail:    true,
		},
		{
			name: "book with both title and author empty",
			book: &models.Book{
				Title:  "",
				Author: "",
				Type:   models.Paperback,
				Notes:  "Notes",
			},
			expectedCount: 2,
			shouldFail:    true,
		},
		{
			name: "book with whitespace-only title",
			book: &models.Book{
				Title:  "   ",
				Author: "Valid Author",
				Type:   models.Digital,
				Notes:  "Notes",
			},
			expectedCount: 1,
			shouldFail:    true,
		},
		{
			name: "book with very long title",
			book: &models.Book{
				Title:  strings.Repeat("a", 300), // Exceeds max length
				Author: "Valid Author",
				Type:   models.Paperback,
				Notes:  "Notes",
			},
			expectedCount: 1,
			shouldFail:    true,
		},
		{
			name: "book with very long notes",
			book: &models.Book{
				Title:  "Valid Title",
				Author: "Valid Author",
				Type:   models.Hardback,
				Notes:  strings.Repeat("a", 1100), // Exceeds max length
			},
			expectedCount: 1,
			shouldFail:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateBook(tt.book)
			
			if len(errors) != tt.expectedCount {
				t.Errorf("ValidateBook() returned %d errors, want %d", len(errors), tt.expectedCount)
				for i, err := range errors {
					t.Errorf("  Error %d: %v", i+1, err)
				}
			}
			
			hasErrors := len(errors) > 0
			if hasErrors != tt.shouldFail {
				t.Errorf("ValidateBook() hasErrors = %v, want %v", hasErrors, tt.shouldFail)
			}
		})
	}
}

// TestValidateTitle tests individual title validation
// This ensures title validation rules are correctly implemented
func TestValidateTitle(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		shouldErr bool
	}{
		{"valid title", "The Great Gatsby", false},
		{"empty title", "", true},
		{"whitespace only title", "   ", true},
		{"title with leading/trailing spaces", "  Valid Title  ", false},
		{"very long title", strings.Repeat("a", 300), true},
		{"title at max length", strings.Repeat("a", 255), false},
		{"title just over max length", strings.Repeat("a", 256), true},
		{"single character title", "A", false},
		{"title with special characters", "Design Patterns: Elements of Reusable Object-Oriented Software", false},
		{"title with numbers", "Catch-22", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTitle(tt.title)
			
			if tt.shouldErr && err == nil {
				t.Errorf("ValidateTitle(%q) should have returned an error", tt.title)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("ValidateTitle(%q) should not have returned an error: %v", tt.title, err)
			}
		})
	}
}

// TestValidateAuthor tests individual author validation
// This ensures author validation rules match title validation patterns
func TestValidateAuthor(t *testing.T) {
	tests := []struct {
		name      string
		author    string
		shouldErr bool
	}{
		{"valid author", "F. Scott Fitzgerald", false},
		{"empty author", "", true},
		{"whitespace only author", "   ", true},
		{"author with leading/trailing spaces", "  Valid Author  ", false},
		{"very long author", strings.Repeat("a", 300), true},
		{"author at max length", strings.Repeat("a", 255), false},
		{"author just over max length", strings.Repeat("a", 256), true},
		{"single character author", "X", false},
		{"author with multiple names", "J.R.R. Tolkien", false},
		{"author with Jr./Sr.", "Martin Luther King Jr.", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAuthor(tt.author)
			
			if tt.shouldErr && err == nil {
				t.Errorf("ValidateAuthor(%q) should have returned an error", tt.author)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("ValidateAuthor(%q) should not have returned an error: %v", tt.author, err)
			}
		})
	}
}

// TestValidateNotes tests notes validation
// Notes are optional but have length limits when provided
func TestValidateNotes(t *testing.T) {
	tests := []struct {
		name      string
		notes     string
		shouldErr bool
	}{
		{"empty notes", "", false},
		{"short notes", "Great book!", false},
		{"notes at max length", strings.Repeat("a", 1000), false},
		{"notes just over max length", strings.Repeat("a", 1001), true},
		{"very long notes", strings.Repeat("a", 2000), true},
		{"notes with newlines", "Line 1\nLine 2\nLine 3", false},
		{"notes with special characters", "Book with émojis! 📚 Very good.", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNotes(tt.notes)
			
			if tt.shouldErr && err == nil {
				t.Errorf("ValidateNotes(%q) should have returned an error", tt.notes)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("ValidateNotes(%q) should not have returned an error: %v", tt.notes, err)
			}
		})
	}
}

// TestValidateFilePath tests file path validation for export operations
// This is critical for ensuring export operations don't fail due to invalid paths
func TestValidateFilePath(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "libros_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name      string
		path      string
		shouldErr bool
	}{
		{"empty path", "", true},
		{"whitespace only path", "   ", true},
		{"valid existing directory", tempDir, false},
		{"nonexistent directory", "/nonexistent/path", true},
		{"valid file in existing directory", filepath.Join(tempDir, "test.txt"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFilePath(tt.path)
			
			if tt.shouldErr && err == nil {
				t.Errorf("ValidateFilePath(%q) should have returned an error", tt.path)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("ValidateFilePath(%q) should not have returned an error: %v", tt.path, err)
			}
		})
	}
}

// TestValidateExportPath tests export path validation and normalization
// This function handles path expansion and validation for export operations
func TestValidateExportPath(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "libros_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name      string
		path      string
		shouldErr bool
		checkAbs  bool // whether to check if result is absolute
	}{
		{"empty path", "", true, false},
		{"whitespace only path", "   ", true, false},
		{"relative path", "test/path", false, true},
		{"absolute path", tempDir, false, true},
		{"path with spaces", "path with spaces", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateExportPath(tt.path)
			
			if tt.shouldErr && err == nil {
				t.Errorf("ValidateExportPath(%q) should have returned an error", tt.path)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("ValidateExportPath(%q) should not have returned an error: %v", tt.path, err)
			}
			
			if !tt.shouldErr && tt.checkAbs && !filepath.IsAbs(result) {
				t.Errorf("ValidateExportPath(%q) should have returned absolute path, got: %q", tt.path, result)
			}
		})
	}
}

// TestTrimAndValidateInput tests the utility function for input trimming and validation
// This function is used throughout the UI for cleaning user input
func TestTrimAndValidateInput(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		fieldName string
		expected  string
		shouldErr bool
	}{
		{"valid input", "test input", "field", "test input", false},
		{"input with leading spaces", "  test input", "field", "test input", false},
		{"input with trailing spaces", "test input  ", "field", "test input", false},
		{"input with both leading and trailing spaces", "  test input  ", "field", "test input", false},
		{"empty input", "", "field", "", true},
		{"whitespace only input", "   ", "field", "", true},
		{"single character input", "a", "field", "a", false},
		{"input with internal spaces", "test   input", "field", "test   input", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TrimAndValidateInput(tt.input, tt.fieldName)
			
			if tt.shouldErr && err == nil {
				t.Errorf("TrimAndValidateInput(%q, %q) should have returned an error", tt.input, tt.fieldName)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("TrimAndValidateInput(%q, %q) should not have returned an error: %v", tt.input, tt.fieldName, err)
			}
			if !tt.shouldErr && result != tt.expected {
				t.Errorf("TrimAndValidateInput(%q, %q) = %q, want %q", tt.input, tt.fieldName, result, tt.expected)
			}
		})
	}
}