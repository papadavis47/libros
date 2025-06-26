// Package screens provides the user interface screens for the Libros book manager application.
// This file contains the book list screen that displays all books in the collection,
// allowing users to browse, navigate, and select books for detailed viewing.
// The screen shows book information including title, author, type, creation date, and truncated notes.
package screens

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/papadavis47/libros/internal/messages"
	"github.com/papadavis47/libros/internal/models"
	"github.com/papadavis47/libros/internal/styles"
)

// ListBooksModel represents the book list screen that displays all books in the collection.
// It manages the list of books, user navigation, error states, and deletion confirmations.
type ListBooksModel struct {
	books      []models.Book // Complete list of books loaded from the database
	index      int           // Currently selected book index (0-based)
	offset     int           // Current scroll offset for viewport
	pageSize   int           // Number of books to display at once
	err        error         // Any error that occurred during book operations
	deleted    bool          // Flag indicating if a book was recently deleted (for showing success message)
}

// NewListBooksModel creates and initializes a new ListBooksModel instance.
// The model starts with an empty book list and the selection index at the first position.
// Books will be loaded asynchronously via LoadBooksMsg messages.
//
// Returns:
//   - ListBooksModel: Initialized list model ready to receive book data
func NewListBooksModel() ListBooksModel {
	return ListBooksModel{
		index:    0, // Start with first item selected
		offset:   0, // Start at top of list
		pageSize: 3, // Show 3 books at a time to prevent overflow
	}
}

// formatDate converts a time.Time to a human-readable date string with ordinal suffix.
// Examples: "January 1st, 2024", "March 23rd, 2024", "April 11th, 2024"
// This provides a more friendly date display than the default Go time formatting.
//
// Parameters:
//   - t: Time value to format
//
// Returns:
//   - string: Formatted date with month name, day with ordinal suffix, and year
func formatDate(t time.Time) string {
	day := t.Day()
	var suffix string
	// Determine appropriate ordinal suffix for the day
	switch {
	case day >= 11 && day <= 13:
		// Special case: 11th, 12th, 13th (not 11st, 12nd, 13rd)
		suffix = "th"
	case day%10 == 1:
		suffix = "st" // 1st, 21st, 31st
	case day%10 == 2:
		suffix = "nd" // 2nd, 22nd
	case day%10 == 3:
		suffix = "rd" // 3rd, 23rd
	default:
		suffix = "th" // 4th, 5th, 6th, 7th, 8th, 9th, 10th, etc.
	}
	return fmt.Sprintf("%s %d%s, %d", t.Format("January"), day, suffix, t.Year())
}

// truncateNotes shortens long note text for display in the book list.
// It attempts to break at word boundaries to avoid cutting words in half,
// and adds an ellipsis (" . . .") to indicate truncation.
//
// Parameters:
//   - notes: Original notes text to potentially truncate
//   - maxLength: Maximum allowed length before truncation
//
// Returns:
//   - string: Original notes if short enough, or truncated version with ellipsis
func truncateNotes(notes string, maxLength int) string {
	if len(notes) <= maxLength {
		return notes // No truncation needed
	}
	// Find a good break point near the limit to avoid cutting words
	truncated := notes[:maxLength]
	// Look for the last space within 20 characters of the limit
	if lastSpace := strings.LastIndex(truncated, " "); lastSpace > maxLength-20 {
		// Break at word boundary if a space is found reasonably close to the limit
		truncated = notes[:lastSpace]
	}
	return truncated + " . . ." // Add ellipsis to indicate truncation
}

// Update handles user input and system messages for the book list screen.
// It processes keyboard navigation, book selection, data loading, and deletion confirmations.
// The function returns the updated model, any commands to execute, the next screen to show,
// and optionally a selected book for the detail screen.
//
// Parameters:
//   - msg: Message to process (keyboard input or system message)
//
// Returns:
//   - ListBooksModel: Updated model state
//   - tea.Cmd: Command to execute (if any)
//   - models.Screen: Next screen to display
//   - *models.Book: Selected book (if navigating to detail screen)
func (m ListBooksModel) Update(msg tea.Msg) (ListBooksModel, tea.Cmd, models.Screen, *models.Book) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc": // Return to main menu
			return m, nil, models.MenuScreen, nil
		case "up", "k": // Move selection up (arrow key or vim key)
			if m.index > 0 {
				m.index--
				// Scroll up if selection moves above viewport
				if m.index < m.offset {
					m.offset = m.index
				}
			}
		case "down", "j": // Move selection down (arrow key or vim key)
			if m.index < len(m.books)-1 {
				m.index++
				// Scroll down if selection moves below viewport
				if m.index >= m.offset+m.pageSize {
					m.offset = m.index - m.pageSize + 1
				}
			}
		case "enter": // Select current book for detailed view
			if len(m.books) > 0 {
				// Get reference to selected book and navigate to detail screen
				selectedBook := &m.books[m.index]
				return m, nil, models.BookDetailScreen, selectedBook
			}
		}

	case messages.LoadBooksMsg: // Handle book data loaded from database
		if msg.Err != nil {
			// Store error for display
			m.err = msg.Err
		} else {
			// Update book list with loaded data
			m.books = msg.Books
			// Ensure selected index is still valid after loading
			if m.index >= len(m.books) && len(m.books) > 0 {
				m.index = len(m.books) - 1
			}
		}

	case messages.DeleteMsg: // Handle book deletion result
		if msg.Err != nil {
			// Store error for display
			m.err = msg.Err
		} else {
			// Set flag to show success message
			m.deleted = true
		}
	}

	// Stay on list screen by default
	return m, nil, models.ListBooksScreen, nil
}

// View renders the book list screen with all books and their details.
// It displays each book's title, author, type, creation date, and truncated notes.
// The currently selected book is highlighted, and the screen shows total count,
// success/error messages, and appropriate help text.
//
// Returns:
//   - string: Formatted book list screen ready for terminal display
func (m ListBooksModel) View() string {
	var b strings.Builder

	// Display application title and screen subtitle
	b.WriteString("\n")
	b.WriteString(styles.TitleStyle.Render("Ｌｉｂｒｏｓ　－　Ａ　Ｂｏｏｋ　Ｍａｎａｇｅｒ"))
	b.WriteString("\n\n")
	b.WriteString(styles.BlurredStyle.Render("Ｙｏｕｒ　Ｂｏｏｋ　Ｃｏｌｌｅｃｔｉｏｎ"))
	b.WriteString("\n\n")

	if len(m.books) == 0 {
		// Show empty state message when no books exist
		b.WriteString(styles.BlurredStyle.Render("No books found. Add some books first!"))
	} else {
		// Calculate visible books based on current offset and page size
		endIndex := m.offset + m.pageSize
		if endIndex > len(m.books) {
			endIndex = len(m.books)
		}
		
		// Display only visible books
		for i := m.offset; i < endIndex; i++ {
			book := m.books[i]
			dateStr := formatDate(book.CreatedAt)
			
			// Create book content with enhanced styling
			var bookContent strings.Builder
			
			if i == m.index {
				// Currently selected book - use enhanced selected styles
				bookContent.WriteString(styles.BookTitleSelectedStyle.Render(styles.AddLetterSpacing(book.Title)))
				bookContent.WriteString("\n\n")
				bookContent.WriteString(styles.BookAuthorSelectedStyle.Render(fmt.Sprintf("%s  %s", styles.AddLetterSpacing("Author:"), styles.AddLetterSpacing(book.Author))))
				bookContent.WriteString("\n\n")
				bookContent.WriteString(styles.SpacedBlurredStyle.Render(fmt.Sprintf("%s %s | %s %s", styles.AddLetterSpacing("Type:"), styles.AddLetterSpacing(styles.CapitalizeBookType(string(book.Type))), styles.AddLetterSpacing("Added:"), styles.AddLetterSpacing(dateStr))))
				if book.Notes != "" {
					// Show truncated notes for selected book
					bookContent.WriteString("\n\n")
					bookContent.WriteString(styles.SpacedNotesStyle.Render(styles.AddLetterSpacing(truncateNotes(book.Notes, 60))))
				}
				
				// Wrap selected book in container
				b.WriteString(styles.BookContainerSelectedStyle.Render(bookContent.String()))
			} else {
				// Non-selected book - use enhanced unselected styles
				bookContent.WriteString(styles.BookTitleUnselectedStyle.Render(styles.AddLetterSpacing(book.Title)))
				bookContent.WriteString("\n\n")
				bookContent.WriteString(styles.BookAuthorUnselectedStyle.Render(fmt.Sprintf("%s  %s", styles.AddLetterSpacing("Author:"), styles.AddLetterSpacing(book.Author))))
				bookContent.WriteString("\n\n")
				bookContent.WriteString(styles.SpacedBlurredStyle.Render(fmt.Sprintf("%s %s | %s %s", styles.AddLetterSpacing("Type:"), styles.AddLetterSpacing(styles.CapitalizeBookType(string(book.Type))), styles.AddLetterSpacing("Added:"), styles.AddLetterSpacing(dateStr))))
				if book.Notes != "" {
					// Show truncated notes for non-selected book too
					bookContent.WriteString("\n\n")
					bookContent.WriteString(styles.SpacedNotesStyle.Render(styles.AddLetterSpacing(truncateNotes(book.Notes, 60))))
				}
				
				// Wrap unselected book in subtle container
				b.WriteString(styles.BookContainerUnselectedStyle.Render(bookContent.String()))
			}
			
			// Add minimal spacing between books
			if i < endIndex-1 {
				b.WriteString("\n")
			}
		}

		// Display total book count and scroll position
		b.WriteString("\n")
		currentPage := (m.offset / m.pageSize) + 1
		totalPages := (len(m.books) + m.pageSize - 1) / m.pageSize
		b.WriteString(styles.BlurredStyle.Render(fmt.Sprintf("   %s %d | %s %d/%d", 
			styles.AddLetterSpacing("Total books:"), len(m.books),
			styles.AddLetterSpacing("Page:"), currentPage, totalPages)))
		b.WriteString("\n")
	}

	// Show success message if a book was recently deleted
	if m.deleted {
		b.WriteString("\n")
		b.WriteString(styles.SuccessStyle.Render(styles.AddLetterSpacing("✓ Book deleted successfully!")))
		b.WriteString("\n")
	}

	// Show any error messages
	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(styles.ErrorStyle.Render(styles.AddLetterSpacing("Error: " + m.err.Error())))
		b.WriteString("\n")
	}

	// Display appropriate help text based on whether books exist
	if len(m.books) > 0 {
		b.WriteString("\n" + styles.HelpTextStyle.Render("Use ↑/↓ or j/k to navigate, Enter to select, Esc to return to menu, q to quit"))
	} else {
		b.WriteString("\n" + styles.HelpTextStyle.Render("Press Esc to return to menu, q or Ctrl+C to quit"))
	}

	return b.String()
}

// ClearDeleted resets the deleted flag to hide the success message.
// This is typically called when navigating away from the list screen
// to ensure the success message doesn't persist across screen transitions.
func (m *ListBooksModel) ClearDeleted() {
	m.deleted = false
}
