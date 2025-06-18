// Package ui provides style re-exports for backwards compatibility and convenience
// This allows UI components to access styling without importing the styles package directly
// It maintains a clean import structure while providing convenient access to all styles
package ui

import "github.com/papadavis47/libros/internal/styles"

// Re-export styles for backwards compatibility and clean imports
// These variable aliases allow UI components to use styles without direct imports
var (
	TitleStyle    = styles.TitleStyle    // Bold purple style for titles and headings
	FocusedStyle  = styles.FocusedStyle  // Purple style for focused UI elements
	BlurredStyle  = styles.BlurredStyle  // Gray style for unfocused UI elements
	NoStyle       = styles.NoStyle       // Plain style with no formatting
	SelectedStyle = styles.SelectedStyle // White on purple background for selected items
	ButtonStyle   = styles.ButtonStyle   // Bold orange style for interactive buttons
	ErrorStyle    = styles.ErrorStyle    // Red style for error messages
	SuccessStyle  = styles.SuccessStyle  // Green style for success messages
)
