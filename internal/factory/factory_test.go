package factory

import (
	"testing"

	"github.com/papadavis47/libros/internal/constants"
)

// TestCreateTextInput tests the generic text input factory function
// This function creates standardized text inputs with consistent configuration
func TestCreateTextInput(t *testing.T) {
	placeholder := "Test placeholder"
	maxLength := 100
	
	input := CreateTextInput(placeholder, maxLength)
	
	// Test basic properties
	if input.Placeholder != placeholder {
		t.Errorf("CreateTextInput() placeholder = %q, want %q", input.Placeholder, placeholder)
	}
	
	if input.CharLimit != maxLength {
		t.Errorf("CreateTextInput() CharLimit = %d, want %d", input.CharLimit, maxLength)
	}
	
	if input.Width != constants.InputFieldWidth {
		t.Errorf("CreateTextInput() Width = %d, want %d", input.Width, constants.InputFieldWidth)
	}
	
	// Test that prompt has consistent formatting
	expectedPrompt := "   " // 3-space left padding
	if input.Prompt != expectedPrompt {
		t.Errorf("CreateTextInput() Prompt = %q, want %q", input.Prompt, expectedPrompt)
	}
}

// TestCreateTitleInput tests the title-specific input factory function
// This function creates inputs with title-specific styling and validation limits
func TestCreateTitleInput(t *testing.T) {
	input := CreateTitleInput()
	
	// Test title-specific properties
	if input.CharLimit != constants.TitleMaxLength {
		t.Errorf("CreateTitleInput() CharLimit = %d, want %d", input.CharLimit, constants.TitleMaxLength)
	}
	
	if input.Width != constants.InputFieldWidth {
		t.Errorf("CreateTitleInput() Width = %d, want %d", input.Width, constants.InputFieldWidth)
	}
	
	// Test that placeholder is appropriate for titles
	expectedPlaceholder := "_______________"
	if input.Placeholder != expectedPlaceholder {
		t.Errorf("CreateTitleInput() Placeholder = %q, want %q", input.Placeholder, expectedPlaceholder)
	}
	
	// Test that the input is initially focused (for title input)
	if !input.Focused() {
		t.Error("CreateTitleInput() should create a focused input")
	}
}

// TestCreateAuthorInput tests the author-specific input factory function
// This function creates inputs with author-specific styling and validation limits
func TestCreateAuthorInput(t *testing.T) {
	input := CreateAuthorInput()
	
	// Test author-specific properties
	if input.CharLimit != constants.AuthorMaxLength {
		t.Errorf("CreateAuthorInput() CharLimit = %d, want %d", input.CharLimit, constants.AuthorMaxLength)
	}
	
	if input.Width != constants.InputFieldWidth {
		t.Errorf("CreateAuthorInput() Width = %d, want %d", input.Width, constants.InputFieldWidth)
	}
	
	// Test that placeholder is appropriate for authors
	expectedPlaceholder := "_______________"
	if input.Placeholder != expectedPlaceholder {
		t.Errorf("CreateAuthorInput() Placeholder = %q, want %q", input.Placeholder, expectedPlaceholder)
	}
	
	// Test that the input is not initially focused (author comes after title)
	if input.Focused() {
		t.Error("CreateAuthorInput() should not create a focused input")
	}
}

// TestCreateNotesTextArea tests the notes-specific textarea factory function
// This function creates textareas optimized for longer text input
func TestCreateNotesTextArea(t *testing.T) {
	textarea := CreateNotesTextArea()
	
	// Test notes-specific properties
	if textarea.CharLimit != constants.NotesMaxLength {
		t.Errorf("CreateNotesTextArea() CharLimit = %d, want %d", textarea.CharLimit, constants.NotesMaxLength)
	}
	
	// Test dimensions - the textarea width may be adjusted by internal padding
	// We test that it's reasonable rather than an exact match
	if textarea.Width() < 40 || textarea.Width() > 60 {
		t.Errorf("CreateNotesTextArea() Width = %d, should be between 40 and 60", textarea.Width())
	}
	
	expectedHeight := 4
	if textarea.Height() != expectedHeight {
		t.Errorf("CreateNotesTextArea() Height = %d, want %d", textarea.Height(), expectedHeight)
	}
	
	// Test that placeholder is appropriate for notes
	expectedPlaceholder := "Notes about this book (optional)..."
	if textarea.Placeholder != expectedPlaceholder {
		t.Errorf("CreateNotesTextArea() Placeholder = %q, want %q", textarea.Placeholder, expectedPlaceholder)
	}
	
	// Test that line numbers are disabled for cleaner appearance
	if textarea.ShowLineNumbers {
		t.Error("CreateNotesTextArea() should disable line numbers")
	}
	
	// Test prompt formatting
	expectedPrompt := "   " // 3-space left padding
	if textarea.Prompt != expectedPrompt {
		t.Errorf("CreateNotesTextArea() Prompt = %q, want %q", textarea.Prompt, expectedPrompt)
	}
}

// TestCreatePathInput tests the path-specific input factory function
// This function creates inputs optimized for file path entry
func TestCreatePathInput(t *testing.T) {
	placeholder := "/home/user/documents"
	
	input := CreatePathInput(placeholder)
	
	// Test path-specific properties
	if input.Placeholder != placeholder {
		t.Errorf("CreatePathInput() Placeholder = %q, want %q", input.Placeholder, placeholder)
	}
	
	if input.Width != constants.TextAreaWidth {
		t.Errorf("CreatePathInput() Width = %d, want %d", input.Width, constants.TextAreaWidth)
	}
	
	// Test prompt formatting
	expectedPrompt := "   " // 3-space left padding
	if input.Prompt != expectedPrompt {
		t.Errorf("CreatePathInput() Prompt = %q, want %q", input.Prompt, expectedPrompt)
	}
}

// TestFactory_Consistency tests that factory functions create consistent components
// This ensures all factory functions follow the same patterns and conventions
func TestFactory_Consistency(t *testing.T) {
	titleInput := CreateTitleInput()
	authorInput := CreateAuthorInput()
	pathInput := CreatePathInput("test")
	notesTextArea := CreateNotesTextArea()
	
	// Test that all text inputs have the same prompt format
	expectedPrompt := "   "
	
	// Check title input prompt contains the expected padding
	if !containsPromptPadding(titleInput.Prompt) {
		t.Error("CreateTitleInput() should have consistent prompt padding")
	}
	
	if !containsPromptPadding(authorInput.Prompt) {
		t.Error("CreateAuthorInput() should have consistent prompt padding")
	}
	
	if pathInput.Prompt != expectedPrompt {
		t.Error("CreatePathInput() should have consistent prompt padding")
	}
	
	if notesTextArea.Prompt != expectedPrompt {
		t.Error("CreateNotesTextArea() should have consistent prompt padding")
	}
	
	// Test that character limits are reasonable
	if titleInput.CharLimit <= 0 {
		t.Error("CreateTitleInput() should have positive character limit")
	}
	
	if authorInput.CharLimit <= 0 {
		t.Error("CreateAuthorInput() should have positive character limit")
	}
	
	if notesTextArea.CharLimit <= 0 {
		t.Error("CreateNotesTextArea() should have positive character limit")
	}
	
	// Test that notes textarea is larger than text inputs for usability
	if notesTextArea.CharLimit <= titleInput.CharLimit {
		t.Error("Notes textarea should have larger character limit than title input")
	}
}

// TestFactory_InputStates tests that factory functions create inputs in appropriate states
// This ensures proper focus management and initial configuration
func TestFactory_InputStates(t *testing.T) {
	titleInput := CreateTitleInput()
	authorInput := CreateAuthorInput()
	
	// Title input should be focused by default (first field in forms)
	if !titleInput.Focused() {
		t.Error("Title input should be focused by default")
	}
	
	// Author input should not be focused by default (second field in forms)
	if authorInput.Focused() {
		t.Error("Author input should not be focused by default")
	}
}

// Helper function to check if a prompt contains the expected padding
func containsPromptPadding(prompt string) bool {
	// The title and author prompts contain "   " plus additional styled text
	// We just check that it starts with the padding
	return len(prompt) > 3 && prompt[:3] == "   "
}

// BenchmarkCreateTitleInput benchmarks the title input creation
// This ensures factory function performance is acceptable
func BenchmarkCreateTitleInput(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateTitleInput()
	}
}

// BenchmarkCreateNotesTextArea benchmarks the notes textarea creation
// This ensures factory function performance is acceptable
func BenchmarkCreateNotesTextArea(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateNotesTextArea()
	}
}