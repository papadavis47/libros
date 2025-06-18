// Package ui provides message type re-exports for backwards compatibility
// This allows UI components to access message types without importing the messages package directly
// It maintains a clean import structure while providing convenient access to message types
package ui

import "github.com/papadavis47/libros/internal/messages"

// Re-export messages for backwards compatibility and clean imports
// These type aliases allow UI components to use message types without direct imports

// SaveMsg represents the result of a book save operation
// Re-exported from messages package for UI component convenience
type SaveMsg = messages.SaveMsg

// UpdateMsg represents the result of a book update operation
// Re-exported from messages package for UI component convenience
type UpdateMsg = messages.UpdateMsg

// DeleteMsg represents the result of a book delete operation
// Re-exported from messages package for UI component convenience
type DeleteMsg = messages.DeleteMsg

// LoadBooksMsg represents the result of loading books from database
// Re-exported from messages package for UI component convenience
type LoadBooksMsg = messages.LoadBooksMsg
