package main

import (
	"os"
	"testing"

	"github.com/papadavis47/libros/internal/database"
	"github.com/papadavis47/libros/internal/ui"
)

func TestDatabaseOperations(t *testing.T) {
	testDBPath := "test_books.db"
	
	// Clean up any existing test database
	os.Remove(testDBPath)
	defer os.Remove(testDBPath)
	
	// Create a test database
	db, err := database.New(testDBPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()
	
	// Test inserting a book
	title := "Test Book"
	author := "Test Author"
	
	err = db.SaveBook(title, author)
	if err != nil {
		t.Fatalf("Failed to save book: %v", err)
	}
	
	// Test retrieving books
	books, err := db.LoadBooks()
	if err != nil {
		t.Fatalf("Failed to load books: %v", err)
	}
	
	if len(books) != 1 {
		t.Fatalf("Expected 1 book, got %d", len(books))
	}
	
	book := books[0]
	if book.Title != title || book.Author != author {
		t.Errorf("Expected title: %s, author: %s, got title: %s, author: %s", title, author, book.Title, book.Author)
	}
	
	// Test updating the book
	newTitle := "Updated Test Book"
	newAuthor := "Updated Test Author"
	
	err = db.UpdateBook(book.ID, newTitle, newAuthor)
	if err != nil {
		t.Fatalf("Failed to update book: %v", err)
	}
	
	// Verify the update
	books, err = db.LoadBooks()
	if err != nil {
		t.Fatalf("Failed to load books after update: %v", err)
	}
	
	if len(books) != 1 {
		t.Fatalf("Expected 1 book after update, got %d", len(books))
	}
	
	updatedBook := books[0]
	if updatedBook.Title != newTitle || updatedBook.Author != newAuthor {
		t.Errorf("Expected updated title: %s, author: %s, got title: %s, author: %s", newTitle, newAuthor, updatedBook.Title, updatedBook.Author)
	}
	
	// Test deleting the book
	err = db.DeleteBook(book.ID)
	if err != nil {
		t.Fatalf("Failed to delete book: %v", err)
	}
	
	// Verify the deletion
	books, err = db.LoadBooks()
	if err != nil {
		t.Fatalf("Failed to load books after deletion: %v", err)
	}
	
	if len(books) != 0 {
		t.Errorf("Expected 0 books after deletion, got %d", len(books))
	}
}

func TestModelInitialization(t *testing.T) {
	testDBPath := "test_init_books.db"
	defer os.Remove(testDBPath)

	db, err := database.New(testDBPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	model := ui.NewModel(db)
	
	// Test that the model initializes properly
	if model.Init() == nil {
		t.Error("Expected Init() to return a command")
	}
}

func TestSaveBookValidation(t *testing.T) {
	testDBPath := "test_validation_books.db"
	defer os.Remove(testDBPath)

	db, err := database.New(testDBPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	// Test empty title and author
	err = db.SaveBook("", "")
	if err == nil {
		t.Error("Expected validation error for empty fields")
	}

	// Test empty title
	err = db.SaveBook("", "Valid Author")
	if err == nil {
		t.Error("Expected validation error for empty title")
	}

	// Test empty author
	err = db.SaveBook("Valid Title", "")
	if err == nil {
		t.Error("Expected validation error for empty author")
	}

	// Test valid input
	err = db.SaveBook("Valid Title", "Valid Author")
	if err != nil {
		t.Errorf("Expected no error for valid input, got: %v", err)
	}
}

func TestBookCount(t *testing.T) {
	testDBPath := "test_count_books.db"
	defer os.Remove(testDBPath)

	db, err := database.New(testDBPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	// Test initial count
	count, err := db.GetBookCount()
	if err != nil {
		t.Fatalf("Failed to get book count: %v", err)
	}
	
	if count != 0 {
		t.Errorf("Expected 0 books initially, got %d", count)
	}

	// Add a book and test count
	err = db.SaveBook("Test Title", "Test Author")
	if err != nil {
		t.Fatalf("Failed to save book: %v", err)
	}

	count, err = db.GetBookCount()
	if err != nil {
		t.Fatalf("Failed to get book count after save: %v", err)
	}
	
	if count != 1 {
		t.Errorf("Expected 1 book after save, got %d", count)
	}
}
