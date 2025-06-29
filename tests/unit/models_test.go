package unit

import (
	"testing"
	"time"

	"github.com/papadavis47/libros/internal/models"
)

// TestBook_Fields tests that Book model fields can be set and retrieved correctly
// This test ensures the model can handle typical use cases without data corruption
func TestBook_Fields(t *testing.T) {
	book := models.Book{
		ID:     1,
		Title:  "The Go Programming Language",
		Author: "Alan Donovan",
		Type:   models.Paperback,
		Notes:  "Excellent reference book",
	}

	// Test that all fields are properly set and can be accessed
	if book.ID != 1 {
		t.Errorf("Expected ID = 1, got %d", book.ID)
	}
	if book.Title != "The Go Programming Language" {
		t.Errorf("Expected Title = 'The Go Programming Language', got %q", book.Title)
	}
	if book.Author != "Alan Donovan" {
		t.Errorf("Expected Author = 'Alan Donovan', got %q", book.Author)
	}
	if book.Type != models.Paperback {
		t.Errorf("Expected Type = Paperback, got %v", book.Type)
	}
	if book.Notes != "Excellent reference book" {
		t.Errorf("Expected Notes = 'Excellent reference book', got %q", book.Notes)
	}
}

// TestBookType_Values tests that all BookType enum values are properly defined
// This ensures the type system correctly represents all supported book formats
func TestBookType_Values(t *testing.T) {
	tests := []struct {
		name     string
		bookType models.BookType
		expected string
	}{
		{"paperback type", models.Paperback, "paperback"},
		{"hardback type", models.Hardback, "hardback"},
		{"audio type", models.Audio, "audio"},
		{"digital type", models.Digital, "digital"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.bookType) != tt.expected {
				t.Errorf("BookType %v = %q, want %q", tt.bookType, string(tt.bookType), tt.expected)
			}
		})
	}
}

// TestBook_Validation tests that Book model fields accept valid data
// This test ensures the model can handle typical use cases without data corruption
func TestBook_Validation(t *testing.T) {
	now := time.Now()
	
	book := models.Book{
		ID:          1,
		Title:       "Test Book",
		Author:      "Test Author",
		Type:        models.Paperback,
		Notes:       "Test notes",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Test that all fields are properly set
	if book.ID != 1 {
		t.Errorf("Expected ID = 1, got %d", book.ID)
	}
	if book.Title != "Test Book" {
		t.Errorf("Expected Title = 'Test Book', got %q", book.Title)
	}
	if book.Author != "Test Author" {
		t.Errorf("Expected Author = 'Test Author', got %q", book.Author)
	}
	if book.Type != models.Paperback {
		t.Errorf("Expected Type = Paperback, got %v", book.Type)
	}
	if book.Notes != "Test notes" {
		t.Errorf("Expected Notes = 'Test notes', got %q", book.Notes)
	}
	if book.CreatedAt != now {
		t.Errorf("Expected CreatedAt = %v, got %v", now, book.CreatedAt)
	}
	if book.UpdatedAt != now {
		t.Errorf("Expected UpdatedAt = %v, got %v", now, book.UpdatedAt)
	}
}

// TestScreenType_Values tests that all screen type constants are properly defined
// This ensures the navigation system has all required screen states
func TestScreenType_Values(t *testing.T) {
	tests := []struct {
		name       string
		screenType models.Screen
		expected   models.Screen
	}{
		{"menu screen", models.MenuScreen, 0},
		{"add book screen", models.AddBookScreen, 1},
		{"list books screen", models.ListBooksScreen, 2},
		{"book detail screen", models.BookDetailScreen, 3},
		{"edit book screen", models.EditBookScreen, 4},
		{"utilities screen", models.UtilitiesScreen, 5},
		{"export screen", models.ExportScreen, 6},
		{"backup screen", models.BackupScreen, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.screenType != tt.expected {
				t.Errorf("Screen type %s = %d, want %d", tt.name, int(tt.screenType), int(tt.expected))
			}
		})
	}
}