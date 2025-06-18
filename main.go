// Package main is the entry point for the Libros book management application
// Libros is a TUI (Terminal User Interface) application built with Bubble Tea
// that allows users to manage their book collection
package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/papadavis47/libros/internal/database"
	"github.com/papadavis47/libros/internal/ui"
)

// main is the entry point of the application
// It initializes the SQLite database, creates the UI model, and starts the Bubble Tea program
func main() {
	// Initialize database connection to books.db SQLite file
	// This will create the database file if it doesn't exist
	db, err := database.New("books.db")
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
	p := tea.NewProgram(model)
	
	// Run the Bubble Tea program and handle any errors
	// This starts the main event loop and renders the UI
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
