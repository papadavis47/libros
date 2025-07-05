package screens

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/papadavis47/libros/internal/database"
	"github.com/papadavis47/libros/internal/models"
	"github.com/papadavis47/libros/internal/styles"
)

// UtilitiesModel represents the utilities menu screen that provides options for
// Export and Backup functionality.
type UtilitiesModel struct {
	db    *database.DB // Database connection for utilities operations
	items []string     // Menu items to display
	index int          // Currently selected menu item index (0-based)
}

// NewUtilitiesModel creates and initializes a new UtilitiesModel instance.
// It sets up the database connection and initializes the menu items.
//
// Parameters:
//   - db: Database connection used for utilities operations
//
// Returns:
//   - UtilitiesModel: Fully initialized utilities model ready for use
func NewUtilitiesModel(db *database.DB) UtilitiesModel {
	items := []string{
		"Ｅｘｐｏｒｔ",
		"Ｂａｃｋｕｐ",
		"Ｂａｃｋ　ｔｏ　Ｍａｉｎ　Ｍｅｎｕ",
	}

	return UtilitiesModel{
		db:    db,
		items: items,
		index: 0,
	}
}

// Update handles keyboard input and user interactions for the utilities screen.
// It processes navigation commands (up/down, vim keys j/k) and selection (enter).
// Based on the selected menu item, it transitions to appropriate screens.
//
// Parameters:
//   - msg: Keyboard message containing the pressed key
//
// Returns:
//   - UtilitiesModel: Updated model state
//   - tea.Cmd: Command to execute (if any)
//   - models.Screen: Next screen to display
func (u UtilitiesModel) Update(msg tea.KeyMsg) (UtilitiesModel, tea.Cmd, models.Screen) {
	switch msg.String() {
	case "up", "k": // Move selection up (arrow key or vim key)
		if u.index > 0 {
			u.index--
		}
	case "down", "j": // Move selection down (arrow key or vim key)
		if u.index < len(u.items)-1 {
			u.index++
		}
	case "enter": // Activate selected menu item
		selectedItem := u.items[u.index]
		switch selectedItem {
		case "Ｅｘｐｏｒｔ":
			// Navigate to export screen for export functionality
			return u, nil, models.ExportScreen
		case "Ｂａｃｋｕｐ":
			// Navigate to database backup functionality
			return u, nil, models.BackupScreen
		case "Ｂａｃｋ　ｔｏ　Ｍａｉｎ　Ｍｅｎｕ":
			// Return to main menu
			return u, nil, models.MenuScreen
		}
	}
	// Return to utilities screen if no action taken
	return u, nil, models.UtilitiesScreen
}

// View renders the utilities menu screen with title, menu options, and help text.
// The currently selected menu item is highlighted using the selected style,
// while other items use the blurred style for visual distinction.
//
// Returns:
//   - string: Formatted utilities screen ready for terminal display
func (u UtilitiesModel) View() string {
	var b strings.Builder

	// Display utilities title
	b.WriteString("\n")
	b.WriteString(styles.TitleStyle().Render("Ｕｔｉｌｉｔｉｅｓ"))
	b.WriteString("\n\n")

	// Render each menu item with appropriate styling
	for i, item := range u.items {
		if i == u.index {
			// Highlight currently selected item
			b.WriteString(styles.SelectedStyle().Render(item))
		} else {
			// Dim non-selected items
			b.WriteString(styles.BlurredStyle.Render(item))
		}
		b.WriteString("\n\n")
	}

	// Display help text for user guidance
	b.WriteString("\n" + styles.HelpTextStyle.Render(styles.AddLetterSpacing("Use ↑/↓ or j/k to navigate, Enter to select, q or Ctrl+C to quit")))

	return b.String()
}