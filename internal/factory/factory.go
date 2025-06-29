package factory

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/papadavis47/libros/internal/constants"
	"github.com/papadavis47/libros/internal/styles"
)

// CreateTextInput creates a standardized text input field with consistent styling
func CreateTextInput(placeholder string, maxLength int) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = maxLength
	ti.Width = constants.InputFieldWidth
	ti.Prompt = "   " // 3-space left padding for alignment
	return ti
}

// CreateTitleInput creates a text input specifically for book titles with custom styling
func CreateTitleInput() textinput.Model {
	ti := textinput.New()
	ti.CharLimit = constants.TitleMaxLength
	ti.Width = constants.InputFieldWidth
	ti.Placeholder = "_______________"
	ti.Prompt = "   " + styles.AddLetterSpacing("Title:") + "  "
	ti.Focus() // Start focused
	ti.PromptStyle = styles.FormFocusedStyle
	ti.TextStyle = styles.FormFocusedStyle
	return ti
}

// CreateAuthorInput creates a text input specifically for book authors with custom styling
func CreateAuthorInput() textinput.Model {
	ti := textinput.New()
	ti.CharLimit = constants.AuthorMaxLength
	ti.Width = constants.InputFieldWidth
	ti.Placeholder = "_______________"
	ti.Prompt = "   " + styles.AddLetterSpacing("Author:") + "  "
	ti.PromptStyle = styles.NoStyle // Remove purple styling to prevent double padding
	return ti
}

// CreateNotesTextArea creates a standardized textarea for book notes
func CreateNotesTextArea() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "Notes about this book (optional)..."
	ta.CharLimit = constants.NotesMaxLength
	ta.SetWidth(constants.InputFieldWidth)
	ta.SetHeight(4)
	ta.ShowLineNumbers = false
	ta.Prompt = "   " // 3-space left padding for alignment
	// Note: Custom styles can be applied from the calling screen if needed
	return ta
}

// CreatePathInput creates a text input for file paths (used in export screen)
func CreatePathInput(placeholder string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Width = constants.TextAreaWidth
	ti.Prompt = "   " // 3-space left padding for alignment
	return ti
}