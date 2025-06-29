// Package models defines the core data structures used throughout the Libros application
// It contains the Book model, book types, and screen navigation constants
package models

import "time"

// BookType represents the different formats a book can be in
// This is stored as a string in the database but used as a typed constant
type BookType string

// Constants defining the available book types/formats
const (
	Paperback BookType = "paperback" // Physical paperback book
	Hardback  BookType = "hardback"  // Physical hardcover book
	Audio     BookType = "audio"     // Audiobook format
	Digital   BookType = "digital"   // Digital/eBook format
)

// Book represents a book record in the database
// Contains all the metadata and user data associated with a book entry
type Book struct {
	ID        int       // Unique database identifier
	Title     string    // Book title
	Author    string    // Book author name
	Type      BookType  // Format type (paperback, hardback, etc.)
	Notes     string    // User notes about the book
	CreatedAt time.Time // When the book record was created
	UpdatedAt time.Time // When the book record was last modified
}

// Screen represents the different UI screens/views in the application
// Used for navigation and state management in the Bubble Tea UI
type Screen int

// Constants defining the available screens in the application
const (
	MenuScreen       Screen = iota // Main menu screen
	AddBookScreen                  // Screen for adding new books
	ListBooksScreen               // Screen showing list of all books
	BookDetailScreen              // Screen showing details of a specific book
	EditBookScreen                // Screen for editing existing book details
	UtilitiesScreen               // Screen for utilities menu (Export/Backup)
	ExportScreen                  // Screen for exporting book data
	BackupScreen                  // Screen for backing up book data
)
