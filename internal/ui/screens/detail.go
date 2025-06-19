// Package screens provides the user interface screens for the Libros book manager application.
// This file contains the book detail screen that displays comprehensive information about a selected book,
// including all metadata, creation/update dates, and full notes. It provides actions for editing,
// deleting, or returning to the book list.
package screens

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/papadavis47/libros/internal/database"
	"github.com/papadavis47/libros/internal/messages"
	"github.com/papadavis47/libros/internal/models"
	"github.com/papadavis47/libros/internal/styles"
)

// DetailModel represents the book detail screen that shows comprehensive information
// about a selected book and provides actions for editing, deleting, or navigation.
type DetailModel struct {
	db           *database.DB  // Database connection for book operations
	SelectedBook *models.Book  // Currently displayed book (set by navigation from list screen)
	actions      []string      // Available actions (Edit, Delete, Back to List)
	index        int           // Currently selected action index (0-based)
	err          error         // Any error from book operations (deletion, etc.)
	updated      bool          // Flag indicating if book was recently updated (for showing success message)
}

// NewDetailModel creates and initializes a new DetailModel instance.
// It sets up the database connection and defines the available actions.
// The selected book will be set later via the SetBook method.
//
// Parameters:
//   - db: Database connection for book operations
//
// Returns:
//   - DetailModel: Initialized detail model ready to display book information
func NewDetailModel(db *database.DB) DetailModel {
	return DetailModel{
		db:      db,
		// Define available actions for the selected book
		actions: []string{"Edit Book", "Delete Book", "Back to List"},
		index:   0, // Start with first action selected
	}
}

// formatDateDetail converts a time.Time to a human-readable date string with ordinal suffix.
// This is identical to formatDate in listbooks.go but kept separate to maintain modularity.
// Examples: "January 1st, 2024", "March 23rd, 2024", "April 11th, 2024"
//
// Parameters:
//   - t: Time value to format
//
// Returns:
//   - string: Formatted date with month name, day with ordinal suffix, and year
func formatDateDetail(t time.Time) string {
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

// wrapText wraps long text to fit within a specified width by breaking at word boundaries.
// This ensures that long notes are displayed properly in the terminal without horizontal scrolling.
// It preserves word integrity by only breaking at spaces.
//
// Parameters:
//   - text: Original text to wrap
//   - width: Maximum line width in characters
//
// Returns:
//   - string: Text with newlines inserted to fit within specified width
func wrapText(text string, width int) string {
	if len(text) <= width {
		return text // No wrapping needed
	}
	
	var result []string
	words := strings.Fields(text) // Split into words
	if len(words) == 0 {
		return text // Handle edge case of empty or whitespace-only text
	}
	
	// Start first line with first word
	currentLine := words[0]
	for _, word := range words[1:] {
		// Check if adding this word would exceed width
		if len(currentLine)+1+len(word) <= width {
			// Add word to current line
			currentLine += " " + word
		} else {
			// Start new line with this word
			result = append(result, currentLine)
			currentLine = word
		}
	}
	// Add the final line
	result = append(result, currentLine)
	
	return strings.Join(result, "\n")
}

// Update handles user input and system messages for the book detail screen.
// It processes navigation between actions, executes selected actions (edit, delete, back),
// and handles responses from update and delete operations.
//
// Parameters:
//   - msg: Message to process (keyboard input or system message)
//
// Returns:
//   - DetailModel: Updated model state
//   - tea.Cmd: Command to execute (if any)
//   - models.Screen: Next screen to display
func (m DetailModel) Update(msg tea.Msg) (DetailModel, tea.Cmd, models.Screen) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc": // Return to book list
			return m, nil, models.ListBooksScreen
		case "up", "k": // Move action selection up
			if m.index > 0 {
				m.index--
			}
		case "down", "j": // Move action selection down
			if m.index < len(m.actions)-1 {
				m.index++
			}
		case "enter": // Execute selected action
			selectedAction := m.actions[m.index]
			switch selectedAction {
			case "Edit Book":
				// Navigate to edit screen
				return m, nil, models.EditBookScreen
			case "Delete Book":
				// Execute delete command and stay on detail screen to show result
				return m, m.deleteBookCmd(), models.BookDetailScreen
			case "Back to List":
				// Return to book list
				return m, nil, models.ListBooksScreen
			}
		}

	case messages.UpdateMsg: // Handle book update result from edit screen
		if msg.Err != nil {
			// Store error for display
			m.err = msg.Err
		} else {
			// Set flag to show success message
			m.updated = true
		}

	case messages.DeleteMsg: // Handle book deletion result
		if msg.Err != nil {
			// Store error for display
			m.err = msg.Err
		} else {
			// Successfully deleted - refresh book list and return to it
			return m, m.loadBooksCmd(), models.ListBooksScreen
		}
	}

	// Stay on detail screen by default
	return m, nil, models.BookDetailScreen
}

// View renders the book detail screen with comprehensive book information and available actions.
// It shows all book metadata, creation/update dates, full notes (with text wrapping),
// and a menu of actions. The screen also displays success/error messages as needed.
//
// Returns:
//   - string: Formatted book detail screen ready for terminal display
func (m DetailModel) View() string {
	var b strings.Builder

	// Display application title and screen subtitle
	b.WriteString(styles.TitleStyle.Render("📚 Libros - Davis Family Book Manager"))
	b.WriteString("\n")
	b.WriteString(styles.BlurredStyle.Render("Book Details"))
	b.WriteString("\n\n")

	if m.SelectedBook != nil {
		// Format the creation and update dates
		createdStr := formatDateDetail(m.SelectedBook.CreatedAt)
		updatedStr := formatDateDetail(m.SelectedBook.UpdatedAt)
		
		// Always show creation date
		b.WriteString(styles.BlurredStyle.Render("Added: " + createdStr) + "\n")
		// Only show update date if it's different from creation date
		if !m.SelectedBook.CreatedAt.Truncate(24*time.Hour).Equal(m.SelectedBook.UpdatedAt.Truncate(24*time.Hour)) {
			b.WriteString(styles.BlurredStyle.Render("Last updated: " + updatedStr) + "\n")
		}
		b.WriteString("\n")
		
		// Display all book metadata with labels
		b.WriteString(styles.FocusedStyle.Render("Title: ") + m.SelectedBook.Title + "\n")
		b.WriteString(styles.FocusedStyle.Render("Author: ") + m.SelectedBook.Author + "\n")
		b.WriteString(styles.FocusedStyle.Render("Type: ") + string(m.SelectedBook.Type) + "\n")
		
		// Display notes if they exist, with text wrapping for readability
		if m.SelectedBook.Notes != "" {
			b.WriteString("\n")
			b.WriteString(styles.FocusedStyle.Render("Notes: ") + "\n")
			// Wrap long notes to fit terminal width
			wrappedNotes := wrapText(m.SelectedBook.Notes, 60)
			b.WriteString(styles.NotesStyle.Render(wrappedNotes) + "\n")
		}
		b.WriteString("\n")

		// Display available actions with selection highlighting
		for i, action := range m.actions {
			if i == m.index {
				// Highlight currently selected action
				b.WriteString(styles.SelectedStyle.Render(action))
			} else {
				// Dim non-selected actions
				b.WriteString(styles.BlurredStyle.Render(action))
			}
			b.WriteString("\n\n")
		}
	}

	// Show success message if book was recently updated
	if m.updated {
		b.WriteString("\n")
		b.WriteString(styles.SuccessStyle.Render("✓ Book updated successfully!"))
		b.WriteString("\n")
	}

	// Show any error messages
	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(styles.ErrorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n")
	}

	// Display help text
	b.WriteString("\n" + styles.BlurredStyle.Render("Use ↑/↓ or j/k to navigate, Enter to select, Esc to go back, q to quit"))

	return b.String()
}

// SetBook sets the book to display and resets the model state.
// This is called when navigating to the detail screen from the book list,
// ensuring the screen shows the correct book and clears any previous state.
//
// Parameters:
//   - book: Pointer to the book to display in detail view
func (m *DetailModel) SetBook(book *models.Book) {
	m.SelectedBook = book
	m.index = 0      // Reset to first action
	m.err = nil      // Clear any previous errors
	m.updated = false // Clear any previous update success message
}

// deleteBookCmd creates a command that asynchronously deletes the currently selected book.
// The command executes the database deletion and returns a DeleteMsg with the result.
// This is called when the user selects the "Delete Book" action.
//
// Returns:
//   - tea.Cmd: Command that deletes the book and returns DeleteMsg
func (m DetailModel) deleteBookCmd() tea.Cmd {
	return func() tea.Msg {
		// Delete the book from the database
		err := m.db.DeleteBook(m.SelectedBook.ID)
		// Return message containing the result (success or error)
		return messages.DeleteMsg{Err: err}
	}
}

// ClearUpdated resets the updated flag to hide the success message.
// This is typically called when navigating away from the detail screen
// to ensure the success message doesn't persist across screen transitions.
func (m *DetailModel) ClearUpdated() {
	m.updated = false
}

// loadBooksCmd creates a command that asynchronously reloads all books from the database.
// This is used after successful book deletion to refresh the book list before returning to it.
// It ensures the deleted book no longer appears in the list.
//
// Returns:
//   - tea.Cmd: Command that loads books and returns LoadBooksMsg
func (m DetailModel) loadBooksCmd() tea.Cmd {
	return func() tea.Msg {
		// Reload all books from database
		books, err := m.db.LoadBooks()
		// Return message containing the refreshed book list
		return messages.LoadBooksMsg{Books: books, Err: err}
	}
}
