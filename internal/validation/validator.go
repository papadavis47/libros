package validation

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/papadavis47/libros/internal/constants"
	"github.com/papadavis47/libros/internal/models"
)

// BookValidationError represents validation errors for book data
type BookValidationError struct {
	Field   string
	Message string
}

func (e BookValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// ValidateBook validates all book fields and returns any validation errors
func ValidateBook(book *models.Book) []error {
	var errors []error
	
	// Validate title
	if err := ValidateTitle(book.Title); err != nil {
		errors = append(errors, err)
	}
	
	// Validate author
	if err := ValidateAuthor(book.Author); err != nil {
		errors = append(errors, err)
	}
	
	// Validate notes (optional field, only validate if present)
	if book.Notes != "" {
		if err := ValidateNotes(book.Notes); err != nil {
			errors = append(errors, err)
		}
	}
	
	return errors
}

// ValidateTitle validates the book title field
func ValidateTitle(title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return BookValidationError{
			Field:   "title",
			Message: "title is required",
		}
	}
	if len(title) > constants.TitleMaxLength {
		return BookValidationError{
			Field:   "title",
			Message: "title exceeds maximum length",
		}
	}
	return nil
}

// ValidateAuthor validates the book author field
func ValidateAuthor(author string) error {
	author = strings.TrimSpace(author)
	if author == "" {
		return BookValidationError{
			Field:   "author",
			Message: "author is required",
		}
	}
	if len(author) > constants.AuthorMaxLength {
		return BookValidationError{
			Field:   "author",
			Message: "author exceeds maximum length",
		}
	}
	return nil
}

// ValidateNotes validates the book notes field
func ValidateNotes(notes string) error {
	if len(notes) > constants.NotesMaxLength {
		return BookValidationError{
			Field:   "notes",
			Message: "notes exceed maximum length",
		}
	}
	return nil
}

// ValidateFilePath validates a file path for export operations
func ValidateFilePath(path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("file path is required")
	}
	
	// Check if directory exists
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return errors.New("directory does not exist: " + dir)
	}
	
	return nil
}

// ValidateExportPath validates and normalizes an export path
func ValidateExportPath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", errors.New("export path is required")
	}
	
	// Expand home directory if needed
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", errors.New("unable to determine home directory")
		}
		path = filepath.Join(homeDir, path[2:])
	}
	
	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", errors.New("invalid file path")
	}
	
	return absPath, nil
}

// TrimAndValidateInput trims whitespace and validates non-empty input
func TrimAndValidateInput(input, fieldName string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", errors.New(fieldName + " cannot be empty")
	}
	return trimmed, nil
}