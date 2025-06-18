// Package main contains unit tests for the Libros application
// These tests verify the core functionality including database operations,
// model initialization, and data validation
package main

import (
	"os"
	"testing"

	"github.com/papadavis47/libros/internal/database"
	"github.com/papadavis47/libros/internal/models"
	"github.com/papadavis47/libros/internal/ui"
)

// TestDatabaseOperations tests the full CRUD (Create, Read, Update, Delete) cycle
// for book operations in the database. This ensures all database functionality works correctly.
func TestDatabaseOperations(t *testing.T) {
	testDBPath := "test_books.db"
	
	// Clean up any existing test database to ensure a clean test environment
	os.Remove(testDBPath)
	defer os.Remove(testDBPath)
	
	// Create a new test database instance
	db, err := database.New(testDBPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()
	
	// Test inserting a book - CREATE operation
	title := "Test Book"
	author := "Test Author"
	
	err = db.SaveBook(title, author, models.Paperback, "")
	if err != nil {
		t.Fatalf("Failed to save book: %v", err)
	}
	
	// Test retrieving books - READ operation
	books, err := db.LoadBooks()
	if err != nil {
		t.Fatalf("Failed to load books: %v", err)
	}
	
	// Verify that exactly one book was saved
	if len(books) != 1 {
		t.Fatalf("Expected 1 book, got %d", len(books))
	}
	
	// Verify the book data was saved correctly
	book := books[0]
	if book.Title != title || book.Author != author || book.Type != models.Paperback {
		t.Errorf("Expected title: %s, author: %s, type: %s, got title: %s, author: %s, type: %s", title, author, models.Paperback, book.Title, book.Author, book.Type)
	}
	
	// Test updating the book - UPDATE operation
	newTitle := "Updated Test Book"
	newAuthor := "Updated Test Author"
	
	err = db.UpdateBook(book.ID, newTitle, newAuthor, models.Hardback, "")
	if err != nil {
		t.Fatalf("Failed to update book: %v", err)
	}
	
	// Verify the update by reloading books from database
	books, err = db.LoadBooks()
	if err != nil {
		t.Fatalf("Failed to load books after update: %v", err)
	}
	
	// Ensure we still have exactly one book
	if len(books) != 1 {
		t.Fatalf("Expected 1 book after update, got %d", len(books))
	}
	
	// Verify the book was updated with new values
	updatedBook := books[0]
	if updatedBook.Title != newTitle || updatedBook.Author != newAuthor || updatedBook.Type != models.Hardback {
		t.Errorf("Expected updated title: %s, author: %s, type: %s, got title: %s, author: %s, type: %s", newTitle, newAuthor, models.Hardback, updatedBook.Title, updatedBook.Author, updatedBook.Type)
	}
	
	// Test deleting the book - DELETE operation
	err = db.DeleteBook(book.ID)
	if err != nil {
		t.Fatalf("Failed to delete book: %v", err)
	}
	
	// Verify the deletion by checking that no books remain
	books, err = db.LoadBooks()
	if err != nil {
		t.Fatalf("Failed to load books after deletion: %v", err)
	}
	
	// Ensure all books have been deleted
	if len(books) != 0 {
		t.Errorf("Expected 0 books after deletion, got %d", len(books))
	}
}

// TestModelInitialization tests that the UI model initializes correctly
// This ensures the Bubble Tea model can be created and initialized properly
func TestModelInitialization(t *testing.T) {
	testDBPath := "test_init_books.db"
	defer os.Remove(testDBPath)

	// Create a test database for the model
	db, err := database.New(testDBPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	// Create a new UI model with the test database
	model := ui.NewModel(db)
	
	// Test that the model initializes and returns an initial command
	// Bubble Tea models should return a command from Init() to start the program
	if model.Init() == nil {
		t.Error("Expected Init() to return a command")
	}
}

// TestSaveBookValidation tests input validation for saving books
// This ensures that invalid data (empty titles/authors) are rejected properly
func TestSaveBookValidation(t *testing.T) {
	testDBPath := "test_validation_books.db"
	defer os.Remove(testDBPath)

	// Create a test database for validation testing
	db, err := database.New(testDBPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	// Test validation: both title and author are empty (should fail)
	err = db.SaveBook("", "", models.Paperback, "")
	if err == nil {
		t.Error("Expected validation error for empty fields")
	}

	// Test validation: empty title with valid author (should fail)
	err = db.SaveBook("", "Valid Author", models.Paperback, "")
	if err == nil {
		t.Error("Expected validation error for empty title")
	}

	// Test validation: valid title with empty author (should fail)
	err = db.SaveBook("Valid Title", "", models.Paperback, "")
	if err == nil {
		t.Error("Expected validation error for empty author")
	}

	// Test validation: both title and author are valid (should succeed)
	err = db.SaveBook("Valid Title", "Valid Author", models.Paperback, "")
	if err != nil {
		t.Errorf("Expected no error for valid input, got: %v", err)
	}
}

// TestBookCount tests the book counting functionality
// This ensures the database correctly tracks the number of books stored
func TestBookCount(t *testing.T) {
	testDBPath := "test_count_books.db"
	defer os.Remove(testDBPath)

	// Create a test database for counting functionality
	db, err := database.New(testDBPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	// Test initial count - should be zero for a new database
	count, err := db.GetBookCount()
	if err != nil {
		t.Fatalf("Failed to get book count: %v", err)
	}
	
	// Verify the database starts with no books
	if count != 0 {
		t.Errorf("Expected 0 books initially, got %d", count)
	}

	// Add a book to the database
	err = db.SaveBook("Test Title", "Test Author", models.Paperback, "")
	if err != nil {
		t.Fatalf("Failed to save book: %v", err)
	}

	// Test count after adding a book - should be one
	count, err = db.GetBookCount()
	if err != nil {
		t.Fatalf("Failed to get book count after save: %v", err)
	}
	
	// Verify the count increased to 1 after adding a book
	if count != 1 {
		t.Errorf("Expected 1 book after save, got %d", count)
	}
}
