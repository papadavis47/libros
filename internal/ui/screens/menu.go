// Package screens provides the user interface screens for the Libros book manager application.
// This file contains the main menu screen that serves as the entry point and navigation hub
// for the application, allowing users to choose between adding books, viewing their collection,
// or quitting the application.
package screens

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/papadavis47/libros/internal/database"
	"github.com/papadavis47/libros/internal/messages"
	"github.com/papadavis47/libros/internal/models"
	"github.com/papadavis47/libros/internal/styles"
)

// MenuModel represents the main menu screen of the application.
// It provides navigation options based on the current state of the book collection.
type MenuModel struct {
	db    *database.DB // Database connection for checking book count and loading books
	items []string     // Menu items to display (dynamically generated based on book count)
	index int          // Currently selected menu item index (0-based)
}

// NewMenuModel creates and initializes a new MenuModel instance.
// It sets up the database connection, initializes the selected index to 0,
// and dynamically generates menu items based on whether books exist in the collection.
//
// Parameters:
//   - db: Database connection used to check book count and load books
//
// Returns:
//   - MenuModel: Fully initialized menu model ready for use
func NewMenuModel(db *database.DB) MenuModel {
	m := MenuModel{
		db:    db,
		index: 0,
	}
	// Generate initial menu items based on current book count
	m.updateMenuItems()
	return m
}

// updateMenuItems dynamically generates menu options based on the current book count.
// If books exist in the collection, it shows "View Books" option; otherwise, it hides it.
// This prevents users from trying to view an empty collection and provides a cleaner UX.
func (m *MenuModel) updateMenuItems() {
	// Get current book count to determine available menu options
	count, err := m.db.GetBookCount()
	if err != nil {
		// On database error, provide minimal menu options
		m.items = []string{"Ａｄｄ　Ｂｏｏｋ", "Ｔｈｅｍｅ", "Ｑｕｉｔ"}
		return
	}

	if count > 0 {
		// Books exist - show all menu options including View Books and Utilities
		m.items = []string{"Ａｄｄ　Ｂｏｏｋ", "Ｖｉｅｗ　Ｂｏｏｋｓ", "Ｕｔｉｌｉｔｉｅｓ", "Ｔｈｅｍｅ", "Ｑｕｉｔ"}
	} else {
		// No books exist - hide View Books and Utilities options
		m.items = []string{"Ａｄｄ　Ｂｏｏｋ", "Ｔｈｅｍｅ", "Ｑｕｉｔ"}
	}

	// Ensure selected index is still valid after menu items change
	// This prevents index out of bounds when menu shrinks
	if m.index >= len(m.items) {
		m.index = len(m.items) - 1
	}
}

// Update handles keyboard input and user interactions for the menu screen.
// It processes navigation commands (up/down, vim keys j/k) and selection (enter).
// Based on the selected menu item, it transitions to appropriate screens or quits the app.
//
// Parameters:
//   - msg: Keyboard message containing the pressed key
//
// Returns:
//   - MenuModel: Updated model state
//   - tea.Cmd: Command to execute (if any)
//   - models.Screen: Next screen to display
func (m MenuModel) Update(msg tea.KeyMsg) (MenuModel, tea.Cmd, models.Screen) {
	switch msg.String() {
	case "up", "k": // Move selection up (arrow key or vim key)
		if m.index > 0 {
			m.index--
		}
	case "down", "j": // Move selection down (arrow key or vim key)
		if m.index < len(m.items)-1 {
			m.index++
		}
	case "enter": // Activate selected menu item
		selectedItem := m.items[m.index]
		switch selectedItem {
		case "Ａｄｄ　Ｂｏｏｋ":
			// Navigate to book creation screen
			return m, nil, models.AddBookScreen
		case "Ｖｉｅｗ　Ｂｏｏｋｓ":
			// Load books from database and navigate to list screen
			// The LoadBooksCmd will fetch data asynchronously
			return m, m.LoadBooksCmd(), models.ListBooksScreen
		case "Ｕｔｉｌｉｔｉｅｓ":
			// Navigate to utilities screen
			return m, nil, models.UtilitiesScreen
		case "Ｔｈｅｍｅ":
			// Navigate to theme selection screen
			return m, nil, models.ThemeScreen
		case "Ｑｕｉｔ":
			// Exit the application
			return m, tea.Quit, models.MenuScreen
		}
	}
	// Return to menu screen if no action taken
	return m, nil, models.MenuScreen
}

// View renders the main menu screen with title, menu options, and help text.
// The currently selected menu item is highlighted using the selected style,
// while other items use the blurred style for visual distinction.
//
// Returns:
//   - string: Formatted menu screen ready for terminal display
func (m MenuModel) View() string {
	var b strings.Builder

	// Display application title with emoji and branding
	b.WriteString("\n")
	b.WriteString(styles.TitleStyle().Render("Ｌｉｂｒｏｓ　－　Ａ　Ｂｏｏｋ　Ｍａｎａｇｅｒ"))
	b.WriteString("\n\n")

	// Render each menu item with appropriate styling
	for i, item := range m.items {
		if i == m.index {
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

// RefreshItems updates the menu items based on the current database state.
// This is typically called when returning to the menu from other screens
// to ensure the "View Books" option appears/disappears based on book count.
func (m *MenuModel) RefreshItems() {
	m.updateMenuItems()
}

// LoadBooksCmd creates a command that asynchronously loads all books from the database.
// This command is executed when the user selects "View Books" from the menu.
// It returns a LoadBooksMsg that will be processed by the list books screen.
//
// Returns:
//   - tea.Cmd: Command that loads books and returns LoadBooksMsg
func (m MenuModel) LoadBooksCmd() tea.Cmd {
	return func() tea.Msg {
		// Load all books from database
		books, err := m.db.LoadBooks()
		// Return message containing books data and any error
		return messages.LoadBooksMsg{Books: books, Err: err}
	}
}
