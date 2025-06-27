// Package screens contains all the individual screen models for the Libros application
// This file implements the AddBookModel which handles the "Add New Book" functionality
// Users can input book title, author, select book type, and add optional notes
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

// AddBookModel represents the "Add New Book" screen state and UI elements
// It manages form inputs, book type selection, and user interaction
type AddBookModel struct {
	db           *database.DB      // Database connection for saving books
	inputs       []textinput.Model // Text input fields [0]=title, [1]=author
	textarea     textarea.Model    // Multi-line text area for optional notes
	bookTypes    []models.BookType // Available book types (paperback, hardback, etc.)
	selectedType int               // Currently selected book type index
	focused      int               // Index of currently focused UI element
	err          error             // Error from save operation, if any
	saved        bool              // Flag indicating if book was successfully saved
}

// NewAddBookModel creates and initializes a new AddBookModel instance
// It sets up the form with text inputs, textarea, and book type options
// The title field is focused by default for immediate user input
func NewAddBookModel(db *database.DB) AddBookModel {
	m := AddBookModel{
		db:           db,                                                                                 // Store database connection
		inputs:       make([]textinput.Model, 2),                                                         // Create title and author inputs
		bookTypes:    []models.BookType{models.Paperback, models.Hardback, models.Audio, models.Digital}, // All available book types
		selectedType: 0,                                                                                  // Default to first type (Paperback)
		focused:      0,                                                                                  // Start focus on title field
	}

	// Initialize and configure the text input fields
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CharLimit = 255 // Maximum characters per field
		t.Width = 50      // Visual width of input field

		// Configure each input field with specific prompts and placeholders
		switch i {
		case 0: // Title field
			t.Placeholder = "_______________"
			t.Prompt = "   " + styles.AddLetterSpacing("Title:") + "  "
			t.Focus() // Start with title field focused
			t.PromptStyle = styles.FormFocusedStyle
			t.TextStyle = styles.FormFocusedStyle
		case 1: // Author field
			t.Placeholder = "_______________"
			t.Prompt = "   " + styles.AddLetterSpacing("Author:") + " "
			t.PromptStyle = styles.NoStyle // Remove purple styling to prevent double padding
			// Author field starts unfocused (default styling)
		}

		m.inputs[i] = t
	}

	// Initialize the textarea for optional book notes
	ta := textarea.New()
	ta.Placeholder = "Notes about this book (optional)..."
	ta.SetWidth(50)            // Match width of text inputs
	ta.SetHeight(4)            // Multi-line height for longer notes
	ta.CharLimit = 1000        // Reasonable limit for notes length
	ta.ShowLineNumbers = false // Disable line numbers
	ta.Prompt = "   "          // 3-space left padding for alignment
	m.textarea = ta

	return m
}

// Update handles all user input and state changes for the Add Book screen
// It processes keyboard input, form navigation, book type selection, and form submission
// Returns the updated model, any commands to execute, and potential screen transitions
func (m AddBookModel) Update(msg tea.Msg) (AddBookModel, tea.Cmd, models.Screen) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc": // Escape key returns to main menu
			m.err = nil     // Clear any error state
			m.saved = false // Clear saved status
			return m, nil, models.MenuScreen
		case "ctrl+a":
			if m.focused < len(m.inputs) {
				m.inputs[m.focused].CursorStart()
			}
			return m, nil, models.AddBookScreen
		case "ctrl+e":
			if m.focused < len(m.inputs) {
				m.inputs[m.focused].CursorEnd()
			}
			return m, nil, models.AddBookScreen
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focused == len(m.inputs)+2 {
				return m, m.saveBookCmd(), models.AddBookScreen
			}

			// Handle tab within type field to cycle through book types
			if m.focused == len(m.inputs) && (s == "tab" || s == "shift+tab") {
				if s == "shift+tab" {
					m.selectedType--
					if m.selectedType < 0 {
						m.selectedType = len(m.bookTypes) - 1
					}
				} else { // tab
					m.selectedType++
					if m.selectedType >= len(m.bookTypes) {
						m.selectedType = 0
					}
				}
				return m, nil, models.AddBookScreen
			}

			// Only allow up/down/enter to move between fields
			if s == "up" {
				m.focused--
			} else if s == "down" || s == "enter" {
				m.focused++
			}

			if m.focused >= len(m.inputs)+3 {
				m.focused = 0
			} else if m.focused < 0 {
				m.focused = len(m.inputs) + 2
			}

			// Update focus for navigation keys
			cmds := make([]tea.Cmd, len(m.inputs)+1)
			for i := 0; i < len(m.inputs); i++ {
				if i == m.focused {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = styles.FormFocusedStyle
					m.inputs[i].TextStyle = styles.FormFocusedStyle
					m.inputs[i].CursorEnd()
				} else {
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = styles.NoStyle
					m.inputs[i].TextStyle = styles.NoStyle
				}
			}

			// Handle textarea focus
			if m.focused == len(m.inputs)+1 {
				cmds[len(m.inputs)] = m.textarea.Focus()
			} else {
				m.textarea.Blur()
			}

			return m, tea.Batch(cmds...), models.AddBookScreen

		case "left", "right":
			s := msg.String()

			// Handle book type selection with left/right arrows when focused on type field
			if m.focused == len(m.inputs) {
				if s == "left" {
					m.selectedType--
					if m.selectedType < 0 {
						m.selectedType = len(m.bookTypes) - 1
					}
				} else { // right
					m.selectedType++
					if m.selectedType >= len(m.bookTypes) {
						m.selectedType = 0
					}
				}
				return m, nil, models.AddBookScreen
			}
			// For text fields, let the input handle left/right for cursor movement
			// This will be handled by updateInputs() method
		}

	case messages.SaveMsg:
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.saved = true
			for i := range m.inputs {
				m.inputs[i].SetValue("")
			}
			m.textarea.SetValue("")
			m.focused = 0
			m.inputs[0].Focus()
		}
		return m, nil, models.AddBookScreen
	}

	cmd := m.updateInputs(msg)
	return m, cmd, models.AddBookScreen
}

// updateInputs propagates messages to all input fields and textarea
// This ensures that all form elements receive keyboard input for editing
// Returns a batched command containing all input field commands
func (m *AddBookModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs)+1) // Commands for inputs + textarea

	// Update each text input field (title, author)
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	// Update the notes textarea
	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	cmds[len(m.inputs)] = cmd

	// Return all commands batched together
	return tea.Batch(cmds...)
}

// View renders the Add Book form UI with all input fields, book type selector, and buttons
// It displays the current state including any error or success messages
// The layout includes title, author inputs, book type buttons, notes textarea, and save button
func (m AddBookModel) View() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(styles.TitleStyle.Render("Ｌｉｂｒｏｓ　－　Ａ　Ｂｏｏｋ　Ｍａｎａｇｅｒ"))
	b.WriteString("\n\n")
	b.WriteString(styles.BlurredStyle.Render("Ａｄｄ　Ｎｅｗ　Ｂｏｏｋ"))
	b.WriteString("\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i == 0 {
			b.WriteString("\n\n") // Extra space after Title
		} else {
			b.WriteRune('\n')
		}
	}

	// Add book type selector
	b.WriteString("\n")
	typeLabel := styles.AddLetterSpacing("Type:") + "   "
	if m.focused == len(m.inputs) {
		b.WriteString(styles.FocusedStyle.Render(typeLabel))
	} else {
		b.WriteString(styles.BlurredStyle.Render(typeLabel))
	}

	for i, bookType := range m.bookTypes {
		if i == m.selectedType {
			if m.focused == len(m.inputs) {
				b.WriteString(styles.ButtonStyle.Render(fmt.Sprintf(" %s ", styles.CapitalizeBookType(string(bookType)))))
			} else {
				b.WriteString(styles.FocusedStyle.Render(fmt.Sprintf(" %s ", styles.CapitalizeBookType(string(bookType)))))
			}
		} else {
			b.WriteString(styles.BlurredStyle.Render(fmt.Sprintf(" %s ", styles.CapitalizeBookType(string(bookType)))))
		}
		if i < len(m.bookTypes)-1 {
			b.WriteString(" ")
		}
	}
	b.WriteString("\n")

	// Add notes textarea
	b.WriteString("\n")
	b.WriteString(styles.FocusedStyle.Render(styles.AddLetterSpacing("Notes:") + " "))
	b.WriteString("\n\n")
	b.WriteString(m.textarea.View())

	button := &styles.BlurredStyle
	if m.focused == len(m.inputs)+2 {
		button = &styles.ButtonStyle
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", button.Render("SAVE BOOK"))

	if m.err != nil {
		b.WriteString(styles.ErrorStyle.Render(styles.AddLetterSpacing("Error: " + m.err.Error())))
		b.WriteString("\n")
	}

	if m.saved {
		b.WriteString(styles.SuccessStyle.Render(styles.AddLetterSpacing("Book saved successfully!")))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(styles.HelpTextStyle.Render("Press Esc to return to menu, Ctrl+A/Ctrl+E for start/end of field, q or Ctrl+C to quit"))

	return b.String()
}

// saveBookCmd creates a command that saves the current book data to the database
// It extracts values from all form fields and calls the database save function
// Returns a SaveMsg with either nil (success) or an error
func (m AddBookModel) saveBookCmd() tea.Cmd {
	return func() tea.Msg {
		// Extract values from form fields
		title := m.inputs[0].Value()            // Get title from first input
		author := m.inputs[1].Value()           // Get author from second input
		bookType := m.bookTypes[m.selectedType] // Get selected book type
		notes := m.textarea.Value()             // Get optional notes

		// Attempt to save the book to database
		err := m.db.SaveBook(title, author, bookType, notes)

		// Return result message that will be handled by Update method
		return messages.SaveMsg{Err: err}
	}
}

// Reset clears all form data and returns the screen to its initial state
// This is called when returning to the menu to prepare for the next book entry
// All fields are cleared and focus returns to the title field
func (m *AddBookModel) Reset() {
	// Clear all status flags
	m.err = nil        // Clear any error messages
	m.saved = false    // Clear saved confirmation
	m.focused = 0      // Reset focus to title field
	m.selectedType = 0 // Reset to first book type (Paperback)

	// Clear all text input values
	for i := range m.inputs {
		m.inputs[i].SetValue("")
	}

	// Clear textarea notes
	m.textarea.SetValue("")

	// Reset focus styling - title field focused, others blurred
	m.inputs[0].Focus()
	m.inputs[0].PromptStyle = styles.FormFocusedStyle // Purple for focused
	m.inputs[0].TextStyle = styles.FormFocusedStyle

	// Blur all other input fields
	for i := 1; i < len(m.inputs); i++ {
		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = styles.NoStyle // No special styling
		m.inputs[i].TextStyle = styles.NoStyle
	}

	// Blur the textarea
	m.textarea.Blur()
}
