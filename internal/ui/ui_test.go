package ui_test

import (
	"os"
	"testing"

	"github.com/papadavis47/libros/internal/database"
	"github.com/papadavis47/libros/internal/ui"
)

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