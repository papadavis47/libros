package main

import (
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/papadavis47/libros/internal/database"
	"github.com/papadavis47/libros/internal/ui"
)

// This initializes the SQLite database, creates the UI model, and starts the Bubble Tea program
func main() {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// Create .libros directory in user's home directory if it doesn't exist
	librosDir := filepath.Join(homeDir, ".libros")
	if err := os.MkdirAll(librosDir, 0755); err != nil {
		log.Fatal(err)
	}

	// Set database path to ~/.libros/books.db
	dbPath := filepath.Join(librosDir, "books.db")

	// Initialize database connection to books.db SQLite file
	// This will create the database file if it doesn't exist
	db, err := database.New(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	// Ensure database connection is closed when the program exits
	defer db.Close()

	// Create the main UI model with database connection
	// This model handles all the application state and UI logic
	model := ui.NewModel(db)

	// Create a new Bubble Tea program with our model
	// Bubble Tea is a framework for building terminal applications
	// WithAltScreen enables alternate screen buffer (clears terminal on start, restores on exit)
	p := tea.NewProgram(model, tea.WithAltScreen())

	// Run the Bubble Tea program and handle any errors
	// This starts the main event loop and renders the UI
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
