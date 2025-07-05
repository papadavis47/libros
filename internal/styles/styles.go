// Package styles defines the visual styling for the Libros TUI application
// It uses the Lip Gloss library to create consistent colors, fonts, and layouts
// across all UI components and screens
package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/papadavis47/libros/internal/config"
	"github.com/papadavis47/libros/internal/utils"
)

// GetTitleStyle returns the themed title style
func GetTitleStyle() lipgloss.Style {
	theme := config.GetCurrentTheme()
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(theme.PrimaryColor)).
		PaddingLeft(3)
}

// Style functions that always return current theme-aware styles
// These should be used instead of global variables to ensure theme changes apply immediately

// TitleStyle returns the current themed title style
func TitleStyle() lipgloss.Style {
	return GetTitleStyle()
}

// FocusedStyle returns the current themed focused style
func FocusedStyle() lipgloss.Style {
	return GetFocusedStyle()
}

// SelectedStyle returns the current themed selected style
func SelectedStyle() lipgloss.Style {
	return GetSelectedStyle()
}

// FormFocusedStyle returns the current themed form focused style
func FormFocusedStyle() lipgloss.Style {
	return GetFormFocusedStyle()
}

// SpacedFocusedStyle returns the current themed spaced focused style
func SpacedFocusedStyle() lipgloss.Style {
	return GetSpacedFocusedStyle()
}

// BoldFocusedStyle returns the current themed bold focused style
func BoldFocusedStyle() lipgloss.Style {
	return GetBoldFocusedStyle()
}

// BookTitleSelectedStyle returns the current themed book title selected style
func BookTitleSelectedStyle() lipgloss.Style {
	return GetBookTitleSelectedStyle()
}

// BookTitleUnselectedStyle returns the current themed book title unselected style
func BookTitleUnselectedStyle() lipgloss.Style {
	return GetBookTitleUnselectedStyle()
}

// BookContainerSelectedStyle returns the current themed book container selected style
func BookContainerSelectedStyle() lipgloss.Style {
	return GetBookContainerSelectedStyle()
}

// BookSeparatorBoldStyle returns the current themed book separator bold style
func BookSeparatorBoldStyle() lipgloss.Style {
	return GetBookSeparatorBoldStyle()
}

// Static styles that don't depend on theme
var (
	// BlurredStyle is applied to UI elements that are not currently focused
	// White color for better accessibility
	BlurredStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")). // White color for accessibility
			Padding(0, 1).                         // Consistent horizontal padding
			PaddingLeft(3)                         // 3-space left indent

	// NoStyle is a plain style with no special formatting
	// Used as a neutral base or to reset styling
	NoStyle = lipgloss.NewStyle()

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

	// BlurredNoPaddingStyle is like BlurredStyle but without left padding
	// Used for inline text that shouldn't have extra spacing
	BlurredNoPaddingStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")) // White color for accessibility

	// SpacedBlurredStyle is BlurredStyle with 1.5x letter spacing for enhanced readability
	SpacedBlurredStyle = lipgloss.NewStyle().
				Bold(true).                            // Bold for better visibility
				Foreground(lipgloss.Color("#F5F5F5")). // Light gray color for accessibility
				Padding(0, 1).                         // Consistent horizontal padding
				PaddingLeft(3)                         // 3-space left indent

	// SpacedNotesStyle is NotesStyle with 1.5x letter spacing for enhanced readability
	SpacedNotesStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")). // White color for accessibility
				Italic(true).                          // Italic formatting for notes
				Padding(0, 1).                         // Consistent horizontal padding
				PaddingLeft(3)                         // 3-space left indent

	// BoldBlurredNoPaddingStyle is BlurredNoPaddingStyle with bold formatting
	BoldBlurredNoPaddingStyle = lipgloss.NewStyle().
					Bold(true).                           // Bold formatting
					Foreground(lipgloss.Color("#FFFFFF")) // White color for accessibility

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

	// BookTypeSelectedStyle for selected book type buttons (no italics)
	BookTypeSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFD700")). // Gold color for contrast
				Padding(0, 1).                         // Consistent horizontal padding
				PaddingLeft(3)                         // Same left alignment as other elements

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

// Theme-aware style functions
// These functions return styles based on the current theme configuration

// GetFocusedStyle returns the themed focused style
func GetFocusedStyle() lipgloss.Style {
	theme := config.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.PrimaryColor)).
		PaddingLeft(3)
}

// GetSelectedStyle returns the themed selected style
func GetSelectedStyle() lipgloss.Style {
	theme := config.GetCurrentTheme()
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color(theme.PrimaryColor)).
		Padding(0, 1).
		MarginLeft(2).
		PaddingLeft(1)
}

// GetFormFocusedStyle returns the themed form focused style
func GetFormFocusedStyle() lipgloss.Style {
	theme := config.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.PrimaryColor))
}

// GetSpacedFocusedStyle returns the themed spaced focused style
func GetSpacedFocusedStyle() lipgloss.Style {
	theme := config.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.PrimaryColor)).
		PaddingLeft(3)
}

// GetBoldFocusedStyle returns the themed bold focused style
func GetBoldFocusedStyle() lipgloss.Style {
	theme := config.GetCurrentTheme()
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(theme.PrimaryColor)).
		PaddingLeft(3)
}

// GetBookTitleSelectedStyle returns the themed book title selected style
func GetBookTitleSelectedStyle() lipgloss.Style {
	theme := config.GetCurrentTheme()
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color(theme.PrimaryColor)).
		Padding(0, 1).
		MarginLeft(2).
		PaddingLeft(1)
}

// GetBookTitleUnselectedStyle returns the themed book title unselected style
func GetBookTitleUnselectedStyle() lipgloss.Style {
	theme := config.GetCurrentTheme()
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(theme.PrimaryColor)).
		PaddingLeft(3)
}

// GetBookContainerSelectedStyle returns the themed book container selected style
func GetBookContainerSelectedStyle() lipgloss.Style {
	theme := config.GetCurrentTheme()
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.PrimaryColor)).
		Padding(1, 2, 1, 0).
		MarginBottom(1)
}

// GetBookSeparatorBoldStyle returns the themed book separator bold style
func GetBookSeparatorBoldStyle() lipgloss.Style {
	theme := config.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.PrimaryColor)).
		Bold(true).
		MarginTop(1).
		MarginBottom(1).
		PaddingLeft(3)
}

