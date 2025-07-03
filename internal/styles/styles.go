// Package styles defines the visual styling for the Libros TUI application
// It uses the Lip Gloss library to create consistent colors, fonts, and layouts
// across all UI components and screens
package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/papadavis47/libros/internal/utils"
)

// These styles provide a consistent look and feel across all UI components
var (
	// TitleStyle is used for main headings and screen titles
	// Purple color with bold text for prominence
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")). // Purple color
			PaddingLeft(3)                         // 3-space left indent

	// FocusedStyle is applied to UI elements that currently have focus
	// Uses the same purple color to indicate active state
	FocusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")). // Purple color for focus
			PaddingLeft(3)                         // 3-space left indent

	// BlurredStyle is applied to UI elements that are not currently focused
	// White color for better accessibility
	BlurredStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")). // White color for accessibility
			Padding(0, 1).                         // Consistent horizontal padding
			PaddingLeft(3)                         // 3-space left indent

	// NoStyle is a plain style with no special formatting
	// Used as a neutral base or to reset styling
	NoStyle = lipgloss.NewStyle()

	// SelectedStyle is used for highlighted/selected items in lists
	// White text on purple background with padding for visual separation
	SelectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")). // White text
			Background(lipgloss.Color("#7D56F4")). // Purple background
			Padding(0, 1).                         // Horizontal padding
			MarginLeft(2).                         // 2-space left margin (creates gap)
			PaddingLeft(1)                         // 1-space left padding inside background

	// ButtonStyle is used for interactive buttons and action items
	// Orange color with bold text to make actions stand out
	ButtonStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFA500")). // Orange color
			PaddingLeft(3)                         // 3-space left indent

	// ErrorStyle is used for error messages and warnings
	// Red color to clearly indicate problems or failures
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")). // Red color
			PaddingLeft(3)                         // 3-space left indent

	// SuccessStyle is used for success messages and confirmations
	// Green color to indicate successful operations
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")). // Green color
			PaddingLeft(3)                         // 3-space left indent

	// NotesStyle is used for displaying book notes with italic formatting
	// White color with italic text for better readability
	NotesStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")). // White color for accessibility
			Italic(true).                          // Italic formatting for notes
			Padding(0, 1).                         // Consistent horizontal padding
			PaddingLeft(3)                         // 3-space left indent

	// FormFocusedStyle is for form inputs that already have padding in their prompts
	// Purple color for focus indication without additional padding
	FormFocusedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")) // Purple color for focus

	// BlurredNoPaddingStyle is like BlurredStyle but without left padding
	// Used for inline text that shouldn't have extra spacing
	BlurredNoPaddingStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")) // White color for accessibility

	// SpacedFocusedStyle is FocusedStyle with 1.5x letter spacing for enhanced readability
	SpacedFocusedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")). // Purple color for focus
				PaddingLeft(3)                         // 3-space left indent

	// SpacedBlurredStyle is BlurredStyle with 1.5x letter spacing for enhanced readability
	SpacedBlurredStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")). // White color for accessibility
				Padding(0, 1).                         // Consistent horizontal padding
				PaddingLeft(3)                         // 3-space left indent

	// SpacedNotesStyle is NotesStyle with 1.5x letter spacing for enhanced readability
	SpacedNotesStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")). // White color for accessibility
				Italic(true).                          // Italic formatting for notes
				Padding(0, 1).                         // Consistent horizontal padding
				PaddingLeft(3)                         // 3-space left indent

	// BoldFocusedStyle is SpacedFocusedStyle with bold formatting for emphasis
	BoldFocusedStyle = lipgloss.NewStyle().
				Bold(true).                            // Bold formatting
				Foreground(lipgloss.Color("#7D56F4")). // Purple color for focus
				PaddingLeft(3)                         // 3-space left indent

	// BoldBlurredNoPaddingStyle is BlurredNoPaddingStyle with bold formatting
	BoldBlurredNoPaddingStyle = lipgloss.NewStyle().
					Bold(true).                           // Bold formatting
					Foreground(lipgloss.Color("#FFFFFF")) // White color for accessibility

	// Enhanced Book Title Styles for better visual emphasis
	// BookTitleSelectedStyle creates a prominent title with background for selected books
	BookTitleSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")). // White text
				Background(lipgloss.Color("#7D56F4")). // Purple background
				Padding(0, 1).                         // Horizontal padding
				MarginLeft(2).                         // Left margin for alignment
				PaddingLeft(1)                         // Internal left padding

	// BookTitleUnselectedStyle creates emphasized title for non-selected books
	BookTitleUnselectedStyle = lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("#7D56F4")). // Purple color
					PaddingLeft(3)                         // Left padding

	// Enhanced Author Styles with complementary colors
	// BookAuthorSelectedStyle for authors of selected books
	BookAuthorSelectedStyle = lipgloss.NewStyle().
				Italic(true).
				Foreground(lipgloss.Color("#FFD700")). // Gold color for contrast
				PaddingLeft(3).                        // Same left alignment as title
				Faint(false)                           // Keep readable on selection

	// BookAuthorUnselectedStyle for authors of non-selected books
	BookAuthorUnselectedStyle = lipgloss.NewStyle().
					Italic(true).
					Foreground(lipgloss.Color("#FFA500")). // Orange color (complementary to purple)
					PaddingLeft(3)                         // Same left alignment as title

	// Book Entry Container Styles
	// BookContainerSelectedStyle creates a bordered container for selected books
	BookContainerSelectedStyle = lipgloss.NewStyle().
					Border(lipgloss.RoundedBorder()).
					BorderForeground(lipgloss.Color("#7D56F4")). // Purple border
					Padding(1, 2, 1, 0).                         // top, right, bottom, left padding
					MarginBottom(1)

	// BookContainerUnselectedStyle creates a subtle container for non-selected books
	BookContainerUnselectedStyle = lipgloss.NewStyle().
					Border(lipgloss.HiddenBorder()). // Invisible border for spacing
					Padding(1, 2, 1, 0).             // top, right, bottom, left padding
					MarginBottom(1)

	// Visual Separator Styles
	// BookSeparatorStyle creates elegant separators between books
	BookSeparatorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")). // White color for accessibility
				MarginTop(1).
				MarginBottom(1).
				PaddingLeft(3)

	// BookSeparatorBoldStyle creates more prominent separators
	BookSeparatorBoldStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")). // Purple separators
				Bold(true).
				MarginTop(1).
				MarginBottom(1).
				PaddingLeft(3)

	// HelpTextStyle creates bold help text for better visibility
	HelpTextStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")). // White color for accessibility
			PaddingLeft(3)                         // Left padding for alignment
)

// AddLetterSpacing converts text to have 1.5x letter spacing by adding spaces between characters
// Example: "Book Title" becomes "B o o k  T i t l e"
func AddLetterSpacing(text string) string {
	if text == "" {
		return text
	}

	var result strings.Builder
	runes := []rune(text)

	for i, r := range runes {
		result.WriteRune(r)
		// Add space after each character except the last one
		// Add double space after space characters to maintain word separation
		if i < len(runes)-1 {
			if r == ' ' {
				result.WriteString(" ")
			} else {
				result.WriteString(" ")
			}
		}
	}

	return result.String()
}

// CapitalizeBookType converts BookType enum values to capitalized display names
// Example: "paperback" becomes "Paperback", "audio" becomes "Audio"
// Deprecated: Use utils.FormatBookType instead
func CapitalizeBookType(bookType string) string {
	return utils.FormatBookType(bookType)
}

// CreateBookSeparator creates a decorative separator line for visual separation between books
// Uses Unicode box-drawing characters for elegant terminal display
func CreateBookSeparator(width int, style lipgloss.Style) string {
	if width <= 0 {
		width = 50 // Default width
	}
	separator := strings.Repeat("─", width)
	return style.Render(separator)
}

// CreateBookDottedSeparator creates a dotted separator line for subtle visual separation
func CreateBookDottedSeparator(width int, style lipgloss.Style) string {
	if width <= 0 {
		width = 50 // Default width
	}
	separator := strings.Repeat("·", width)
	return style.Render(separator)
}
