// Package screens provides the user interface screens for the Libros book manager application.
// This file contains the book editing screen that allows users to modify existing book information
// including title, author, type, and notes. It provides a form-based interface with navigation
// between fields and validation before saving changes to the database.
package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/papadavis47/libros/internal/database"
	"github.com/papadavis47/libros/internal/messages"
	"github.com/papadavis47/libros/internal/models"
	"github.com/papadavis47/libros/internal/styles"
)

// EditModel represents the book editing screen that provides a form interface
// for modifying existing book information. It manages multiple input fields,
// focus navigation, and form validation.
type EditModel struct {
	db           *database.DB      // Database connection for saving changes
	SelectedBook *models.Book      // Book being edited (set by navigation from detail screen)
	inputs       []textinput.Model // Text input fields for title and author
	textarea     textarea.Model    // Multi-line text area for notes
	bookTypes    []models.BookType // Available book types (Paperback, Hardback, etc.)
	selectedType int               // Currently selected book type index
	focused      int               // Currently focused form element (0=title, 1=author, 2=type, 3=notes, 4=button)
	err          error             // Any error from form validation or save operation
}

// NewEditModel creates and initializes a new EditModel instance.
// It sets up all form components including text inputs, textarea, and book type options.
// The model starts with the title field focused and ready for user input.
//
// Parameters:
//   - db: Database connection for saving book changes
//
// Returns:
//   - EditModel: Fully initialized edit model ready for use
func NewEditModel(db *database.DB) EditModel {
	m := EditModel{
		db:     db,
		inputs: make([]textinput.Model, 2), // Title and Author inputs
		// Define available book types in order
		bookTypes:    []models.BookType{models.Paperback, models.Hardback, models.Audio, models.Digital},
		selectedType: 0, // Start with first book type selected
		focused:      0, // Start with title field focused
	}

	// Initialize text input fields
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CharLimit = 255 // Reasonable limit for title and author
		t.Width = 50      // Visual width in terminal

		switch i {
		case 0: // Title field
			t.Placeholder = "Enter book title"
			t.Prompt = "Title:  "
			t.Focus() // Start focused on title
			t.PromptStyle = styles.FocusedStyle
			t.TextStyle = styles.FocusedStyle
		case 1: // Author field
			t.Placeholder = "Enter author name"
			t.Prompt = "Author: "
			// Initially blurred (not focused)
		}

		m.inputs[i] = t
	}

	// Initialize textarea for notes
	ta := textarea.New()
	ta.Placeholder = "Enter notes about this book (optional)..."
	ta.SetWidth(50)            // Match width of text inputs
	ta.SetHeight(4)            // Allow multiple lines
	ta.CharLimit = 1000        // Reasonable limit for notes
	ta.ShowLineNumbers = false // Disable line numbers
	ta.Prompt = " "            // Minimal left padding with single space
	m.textarea = ta

	return m
}

// Update handles user input and system messages for the book editing screen.
// It manages focus navigation between form fields, handles book type selection,
// processes form submission, and responds to save operations from the database.
//
// The focus order is: Title -> Author -> Book Type -> Notes -> Save Button
//
// Parameters:
//   - msg: Message to process (keyboard input or system message)
//
// Returns:
//   - EditModel: Updated model state
//   - tea.Cmd: Command to execute (if any)
//   - models.Screen: Next screen to display
func (m EditModel) Update(msg tea.Msg) (EditModel, tea.Cmd, models.Screen) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc": // Cancel editing and return to detail screen
			m.err = nil // Clear any errors
			return m, nil, models.BookDetailScreen
		case "ctrl+a": // Move cursor to start of current text input
			if m.focused < len(m.inputs) {
				m.inputs[m.focused].CursorStart()
			}
			return m, nil, models.EditBookScreen
		case "ctrl+e": // Move cursor to end of current text input
			if m.focused < len(m.inputs) {
				m.inputs[m.focused].CursorEnd()
			}
			return m, nil, models.EditBookScreen
		case "tab", "shift+tab", "enter", "up", "down", "left", "right":
			s := msg.String()

			// Handle form submission when save button is focused
			if s == "enter" && m.focused == len(m.inputs)+2 {
				return m, m.updateBookCmd(), models.EditBookScreen
			}

			// Handle book type selection with left/right arrows
			if m.focused == len(m.inputs) && (s == "left" || s == "right") {
				if s == "left" {
					m.selectedType--
					// Wrap around to last type
					if m.selectedType < 0 {
						m.selectedType = len(m.bookTypes) - 1
					}
				} else {
					m.selectedType++
					// Wrap around to first type
					if m.selectedType >= len(m.bookTypes) {
						m.selectedType = 0
					}
				}
				return m, nil, models.EditBookScreen
			}

			// Handle focus navigation between form elements
			if s == "up" || s == "shift+tab" {
				// Move focus backward
				m.focused--
			} else {
				// Move focus forward
				m.focused++
			}

			// Wrap focus around (total elements: inputs + book type + notes + save button)
			if m.focused > len(m.inputs)+2 {
				m.focused = 0 // Wrap to first element
			} else if m.focused < 0 {
				m.focused = len(m.inputs) + 2 // Wrap to last element
			}

			// Update focus states for all text inputs
			cmds := make([]tea.Cmd, len(m.inputs)+1)
			for i := 0; i < len(m.inputs); i++ {
				if i == m.focused {
					// Focus this input
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = styles.FocusedStyle
					m.inputs[i].TextStyle = styles.FocusedStyle
					m.inputs[i].CursorEnd() // Position cursor at end
				} else {
					// Blur this input
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = styles.NoStyle
					m.inputs[i].TextStyle = styles.NoStyle
				}
			}

			// Handle textarea focus (notes field)
			if m.focused == len(m.inputs)+1 {
				cmds[len(m.inputs)] = m.textarea.Focus()
			} else {
				m.textarea.Blur()
			}

			return m, tea.Batch(cmds...), models.EditBookScreen
		}

	case messages.UpdateMsg: // Handle save operation result
		if msg.Err != nil {
			// Store error for display
			m.err = msg.Err
		} else {
			// Success - update the book object with new values and return to detail screen
			m.SelectedBook.Title = m.inputs[0].Value()
			m.SelectedBook.Author = m.inputs[1].Value()
			m.SelectedBook.Type = m.bookTypes[m.selectedType]
			m.SelectedBook.Notes = m.textarea.Value()
			return m, nil, models.BookDetailScreen
		}
	}

	// Update all input components with the message
	cmd := m.updateInputs(msg)
	return m, cmd, models.EditBookScreen
}

// updateInputs propagates messages to all input components (text inputs and textarea).
// This ensures that all form elements receive keyboard input and can update their state.
// It's called for messages that aren't handled by the main Update function.
//
// Parameters:
//   - msg: Message to propagate to input components
//
// Returns:
//   - tea.Cmd: Batched commands from all input components
func (m *EditModel) updateInputs(msg tea.Msg) tea.Cmd {
	// Create command slice for all inputs plus textarea
	cmds := make([]tea.Cmd, len(m.inputs)+1)

	// Update all text inputs
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	// Update textarea
	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	cmds[len(m.inputs)] = cmd

	// Return all commands batched together
	return tea.Batch(cmds...)
}

// View renders the book editing form with all input fields, book type selector,
// notes textarea, and save button. It shows visual focus indicators and displays
// any validation errors. The layout provides a clear, navigable interface for editing.
//
// Returns:
//   - string: Formatted edit screen ready for terminal display
func (m EditModel) View() string {
	var b strings.Builder

	// Display application title and screen subtitle
	b.WriteString(styles.TitleStyle.Render("📚 Libros - Davis Family Book Manager"))
	b.WriteString("\n")
	b.WriteString(styles.BlurredStyle.Render("Edit Book"))
	b.WriteString("\n\n")

	// Render all text input fields (title and author)
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i == 0 {
			b.WriteString("\n\n") // Extra space after Title
		} else {
			b.WriteRune('\n')
		}
	}

	// Add book type selector with focus-aware styling
	b.WriteString("\n")
	typeLabel := "Type:   "
	if m.focused == len(m.inputs) {
		// Book type selector is focused
		b.WriteString(styles.FocusedStyle.Render(typeLabel))
	} else {
		// Book type selector is not focused
		b.WriteString(styles.BlurredStyle.Render(typeLabel))
	}

	// Render each book type option with appropriate styling
	for i, bookType := range m.bookTypes {
		if i == m.selectedType {
			// This is the currently selected book type
			if m.focused == len(m.inputs) {
				// Book type selector is focused - use button style
				b.WriteString(styles.ButtonStyle.Render(fmt.Sprintf(" %s ", string(bookType))))
			} else {
				// Book type selector not focused but this type is selected
				b.WriteString(styles.FocusedStyle.Render(fmt.Sprintf(" %s ", string(bookType))))
			}
		} else {
			// This is not the selected book type
			b.WriteString(styles.BlurredStyle.Render(fmt.Sprintf(" %s ", string(bookType))))
		}
		// Add spacing between book type options
		if i < len(m.bookTypes)-1 {
			b.WriteString(" ")
		}
	}
	b.WriteString("\n")

	// Add notes textarea with label
	b.WriteString("\n")
	b.WriteString(styles.FocusedStyle.Render("Notes: "))
	b.WriteString("\n")
	b.WriteString(m.textarea.View())

	// Add save button with focus-aware styling
	button := &styles.BlurredStyle
	if m.focused == len(m.inputs)+2 {
		// Save button is focused
		button = &styles.ButtonStyle
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", button.Render("UPDATE BOOK"))

	// Show any validation or save errors
	if m.err != nil {
		b.WriteString(styles.ErrorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n")
	}

	// Display help text
	b.WriteString(styles.BlurredStyle.Render("Press Esc to cancel, Ctrl+A/Ctrl+E for start/end of field, q or Ctrl+C to quit"))

	return b.String()
}

// SetBook populates the edit form with the selected book's current information.
// This is called when navigating to the edit screen from the detail screen,
// ensuring all form fields are pre-filled with existing book data.
//
// Parameters:
//   - book: Pointer to the book to edit
func (m *EditModel) SetBook(book *models.Book) {
	m.SelectedBook = book
	m.focused = 0 // Start with title field focused
	m.err = nil   // Clear any previous errors

	// Populate text input fields with current book data
	m.inputs[0].SetValue(book.Title)
	m.inputs[1].SetValue(book.Author)
	m.textarea.SetValue(book.Notes)

	// Find and select the current book type in the selector
	for i, bookType := range m.bookTypes {
		if bookType == book.Type {
			m.selectedType = i
			break // Found matching type
		}
	}

	// Set initial focus state - title field focused, others blurred
	m.inputs[0].Focus()
	m.inputs[0].CursorEnd() // Position cursor at end of title
	for i := 1; i < len(m.inputs); i++ {
		m.inputs[i].Blur() // Ensure other inputs are not focused
	}
	m.textarea.Blur() // Ensure textarea is not focused
}

// updateBookCmd creates a command that asynchronously saves the edited book to the database.
// It collects all form values and calls the database update method with the current book's ID.
// The command returns an UpdateMsg with the result (success or error).
//
// Returns:
//   - tea.Cmd: Command that updates the book and returns UpdateMsg
func (m EditModel) updateBookCmd() tea.Cmd {
	return func() tea.Msg {
		// Extract values from all form fields
		title := m.inputs[0].Value()            // Title from first input
		author := m.inputs[1].Value()           // Author from second input
		bookType := m.bookTypes[m.selectedType] // Selected book type
		notes := m.textarea.Value()             // Notes from textarea

		// Update the book in the database
		err := m.db.UpdateBook(m.SelectedBook.ID, title, author, bookType, notes)

		// Return message containing the result
		return messages.UpdateMsg{Err: err}
	}
}
