// Package messages defines Bubble Tea command messages for the Libros application
// These messages are sent between the UI and business logic to communicate
// the results of database operations and other asynchronous actions
package messages

import "github.com/papadavis47/libros/internal/models"

// SaveMsg represents the result of a book save operation
// Contains an error field to indicate success (nil) or failure (error details)
type SaveMsg struct {
	Err error // Error from the save operation, nil if successful
}

// UpdateMsg represents the result of a book update operation
// Contains an error field to indicate success (nil) or failure (error details)
type UpdateMsg struct {
	Err error // Error from the update operation, nil if successful
}

// DeleteMsg represents the result of a book delete operation
// Contains an error field to indicate success (nil) or failure (error details)
type DeleteMsg struct {
	Err error // Error from the delete operation, nil if successful
}

// LoadBooksMsg represents the result of loading books from the database
// Contains both the loaded books data and any error that occurred
type LoadBooksMsg struct {
	Books []models.Book // Slice of books loaded from database
	Err   error         // Error from the load operation, nil if successful
}
