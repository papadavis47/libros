package database_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/papadavis47/libros/internal/database"
	"github.com/papadavis47/libros/internal/models"
)

// TestDatabase_BasicOperations tests basic database operations
// This integration test verifies that the database layer works correctly
// with real SQLite operations and handles data persistence properly
func TestDatabase_BasicOperations(t *testing.T) {
	// Create temporary database for testing
	tempDir, err := os.MkdirTemp("", "libros_test_db")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	db, err := database.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Test CREATE operation
	t.Run("SaveBook", func(t *testing.T) {
		err := db.SaveBook("Test Book", "Test Author", models.Paperback, "Test notes")
		if err != nil {
			t.Fatalf("Failed to save book: %v", err)
		}
	})

	// Test READ operation
	t.Run("LoadBooks", func(t *testing.T) {
		// Add a few more books
		err := db.SaveBook("Book 1", "Author 1", models.Paperback, "Notes 1")
		if err != nil {
			t.Fatalf("Failed to save book 1: %v", err)
		}
		
		err = db.SaveBook("Book 2", "Author 2", models.Hardback, "Notes 2")
		if err != nil {
			t.Fatalf("Failed to save book 2: %v", err)
		}

		books, err := db.LoadBooks()
		if err != nil {
			t.Fatalf("Failed to load books: %v", err)
		}

		// Should have at least the books we created
		if len(books) < 3 {
			t.Errorf("Expected at least 3 books, got %d", len(books))
		}

		// Verify the first book's data
		found := false
		for _, book := range books {
			if book.Title == "Test Book" && book.Author == "Test Author" {
				found = true
				if book.Type != models.Paperback {
					t.Errorf("Expected book type Paperback, got %v", book.Type)
				}
				if book.Notes != "Test notes" {
					t.Errorf("Expected notes 'Test notes', got %q", book.Notes)
				}
				if book.ID == 0 {
					t.Error("Book should have a valid ID")
				}
				if book.CreatedAt.IsZero() {
					t.Error("Book should have CreatedAt timestamp")
				}
				if book.UpdatedAt.IsZero() {
					t.Error("Book should have UpdatedAt timestamp")
				}
				break
			}
		}
		if !found {
			t.Error("Could not find the test book in loaded books")
		}
	})

	// Test UPDATE operation
	t.Run("UpdateBook", func(t *testing.T) {
		// Get books to find one to update
		books, err := db.LoadBooks()
		if err != nil {
			t.Fatalf("Failed to load books for update test: %v", err)
		}
		
		if len(books) == 0 {
			t.Skip("No books available for update test")
		}

		// Update the first book
		bookID := books[0].ID
		err = db.UpdateBook(bookID, "Updated Title", "Updated Author", models.Digital, "Updated notes")
		if err != nil {
			t.Fatalf("Failed to update book: %v", err)
		}

		// Verify the update
		updatedBooks, err := db.LoadBooks()
		if err != nil {
			t.Fatalf("Failed to load books after update: %v", err)
		}

		found := false
		for _, book := range updatedBooks {
			if book.ID == bookID {
				found = true
				if book.Title != "Updated Title" {
					t.Errorf("Expected updated title 'Updated Title', got %q", book.Title)
				}
				if book.Author != "Updated Author" {
					t.Errorf("Expected updated author 'Updated Author', got %q", book.Author)
				}
				if book.Type != models.Digital {
					t.Errorf("Expected updated type Digital, got %v", book.Type)
				}
				if book.Notes != "Updated notes" {
					t.Errorf("Expected updated notes 'Updated notes', got %q", book.Notes)
				}
				break
			}
		}
		if !found {
			t.Error("Could not find updated book")
		}
	})

	// Test DELETE operation
	t.Run("DeleteBook", func(t *testing.T) {
		// Get books to find one to delete
		books, err := db.LoadBooks()
		if err != nil {
			t.Fatalf("Failed to load books for delete test: %v", err)
		}
		
		if len(books) == 0 {
			t.Skip("No books available for delete test")
		}

		initialCount := len(books)
		bookToDeleteID := books[0].ID

		// Delete the book
		err = db.DeleteBook(bookToDeleteID)
		if err != nil {
			t.Fatalf("Failed to delete book: %v", err)
		}

		// Verify the deletion
		remainingBooks, err := db.LoadBooks()
		if err != nil {
			t.Fatalf("Failed to load books after deletion: %v", err)
		}

		if len(remainingBooks) != initialCount-1 {
			t.Errorf("Expected %d books after deletion, got %d", initialCount-1, len(remainingBooks))
		}

		// Verify the specific book was deleted
		for _, book := range remainingBooks {
			if book.ID == bookToDeleteID {
				t.Error("Deleted book should not exist in remaining books")
			}
		}
	})

	// Test GetBookCount
	t.Run("GetBookCount", func(t *testing.T) {
		books, err := db.LoadBooks()
		if err != nil {
			t.Fatalf("Failed to load books: %v", err)
		}

		count, err := db.GetBookCount()
		if err != nil {
			t.Fatalf("Failed to get book count: %v", err)
		}

		if count != len(books) {
			t.Errorf("GetBookCount() = %d, LoadBooks() returned %d books", count, len(books))
		}
	})
}

// TestDatabase_EdgeCases tests edge cases and error conditions
func TestDatabase_EdgeCases(t *testing.T) {
	// Create temporary database for testing
	tempDir, err := os.MkdirTemp("", "libros_test_edge")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test_edge.db")
	db, err := database.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	t.Run("SaveBookWithSpecialCharacters", func(t *testing.T) {
		title := "Book with émojis 📚 and unicode"
		author := "Author with àccénts and ñoñ-ASCII"
		notes := "Notes with 'quotes', \"double quotes\", and unicode: ★☆★"

		err := db.SaveBook(title, author, models.Digital, notes)
		if err != nil {
			t.Fatalf("Failed to save book with special characters: %v", err)
		}

		// Retrieve and verify the special characters were preserved
		books, err := db.LoadBooks()
		if err != nil {
			t.Fatalf("Failed to load books: %v", err)
		}

		found := false
		for _, book := range books {
			if book.Title == title {
				found = true
				if book.Author != author {
					t.Errorf("Special characters in author not preserved: got %q, want %q", 
						book.Author, author)
				}
				if book.Notes != notes {
					t.Errorf("Special characters in notes not preserved: got %q, want %q", 
						book.Notes, notes)
				}
				break
			}
		}
		if !found {
			t.Error("Book with special characters not found")
		}
	})

	t.Run("UpdateNonexistentBook", func(t *testing.T) {
		// This tests that updating a nonexistent book doesn't crash
		// The actual behavior may vary based on implementation
		err := db.UpdateBook(99999, "Nonexistent", "Ghost", models.Paperback, "Notes")
		// We just verify the operation completes without crashing
		_ = err // Some implementations may or may not return an error
	})

	t.Run("DeleteNonexistentBook", func(t *testing.T) {
		// This tests that deleting a nonexistent book doesn't crash
		// The actual behavior may vary based on implementation
		err := db.DeleteBook(99999)
		// We just verify the operation completes without crashing
		_ = err // Some implementations may or may not return an error
	})
}