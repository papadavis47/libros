// Package ui provides the main Bubble Tea model for the Libros application
// This is the root model that coordinates between different screen models
// and manages the overall application state and navigation
package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/papadavis47/libros/internal/database"
	"github.com/papadavis47/libros/internal/models"
	"github.com/papadavis47/libros/internal/ui/screens"
)

// Model represents the main application model that coordinates all screen models
// It implements the Bubble Tea model interface and manages navigation between screens
type Model struct {
	db            *database.DB    // Database connection shared across all screens
	currentScreen models.Screen   // Current active screen being displayed
	
	// Screen models - each screen has its own model that handles specific functionality
	menu      screens.MenuModel      // Main menu screen model
	addBook   screens.AddBookModel   // Add new book screen model
	listBooks screens.ListBooksModel // Book list display screen model
	detail    screens.DetailModel    // Book detail view screen model
	edit      screens.EditModel      // Book editing screen model
	backup    *screens.BackupScreen  // Backup data screen model
}

// NewModel creates and initializes a new main application model
// It takes a database connection and creates all the individual screen models
// The application starts on the MenuScreen by default
func NewModel(db *database.DB) Model {
	return Model{
		db:            db,                              // Store database connection
		currentScreen: models.MenuScreen,               // Start at main menu
		menu:          screens.NewMenuModel(db),        // Initialize menu screen
		addBook:       screens.NewAddBookModel(db),     // Initialize add book screen
		listBooks:     screens.NewListBooksModel(),     // Initialize book list screen
		detail:        screens.NewDetailModel(db),      // Initialize detail view screen
		edit:          screens.NewEditModel(db),        // Initialize edit screen
		backup:        screens.NewBackupScreen(db),     // Initialize backup screen
	}
}

// Init initializes the Bubble Tea model and returns the initial command
// This is called once when the program starts to set up the initial state
func (m Model) Init() tea.Cmd {
	// Return the text input blink command to start cursor blinking
	return textinput.Blink
}

// Update handles incoming messages and updates the model state
// It processes global key commands and delegates screen-specific updates to individual screen models
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle global key commands that work across all screens
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Ctrl+C always quits the application immediately
		if msg.String() == "ctrl+c" {
			m.db.Close() // Clean up database connection
			return m, tea.Quit
		}
		// 'q' quits the application, but not from input screens (add/edit)
		// This prevents accidental quits while typing
		if msg.String() == "q" && m.currentScreen != models.AddBookScreen && m.currentScreen != models.EditBookScreen {
			m.db.Close() // Clean up database connection
			return m, tea.Quit
		}
	}

	// Variables to track the command to execute and potential screen changes
	var cmd tea.Cmd
	var newScreen models.Screen

	// Delegate message handling to the appropriate screen model based on current screen
	switch m.currentScreen {
	case models.MenuScreen:
		// Menu only handles key messages, other messages are ignored
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			var menuCmd tea.Cmd
			// Update menu model and get any screen transition
			m.menu, menuCmd, newScreen = m.menu.Update(keyMsg)
			cmd = menuCmd
		} else {
			// No screen change if message isn't a key press
			newScreen = m.currentScreen
		}
		
	case models.AddBookScreen:
		var addBookCmd tea.Cmd
		// Update add book screen model
		m.addBook, addBookCmd, newScreen = m.addBook.Update(msg)
		cmd = addBookCmd
		// Clean up state when returning to menu
		if newScreen == models.MenuScreen {
			m.addBook.Reset()         // Clear form inputs
			m.menu.RefreshItems()     // Refresh menu to show new book count
		}
		
	case models.ListBooksScreen:
		var listCmd tea.Cmd
		var selectedBook *models.Book
		// Update list screen and check if a book was selected
		m.listBooks, listCmd, newScreen, selectedBook = m.listBooks.Update(msg)
		cmd = listCmd
		// If a book was selected, prepare it for the detail screen
		if selectedBook != nil {
			m.detail.SetBook(selectedBook)
		}
		// Refresh menu when returning (in case books were deleted)
		if newScreen == models.MenuScreen {
			m.menu.RefreshItems()
		}
		
	case models.BookDetailScreen:
		var detailCmd tea.Cmd
		// Update detail screen model
		m.detail, detailCmd, newScreen = m.detail.Update(msg)
		cmd = detailCmd
		// If transitioning to edit screen, pass the current book data
		if newScreen == models.EditBookScreen {
			m.edit.SetBook(m.detail.SelectedBook)
		}
		
	case models.EditBookScreen:
		var editCmd tea.Cmd
		// Update edit screen model
		m.edit, editCmd, newScreen = m.edit.Update(msg)
		cmd = editCmd
		// Clear update status when returning to detail screen
		if newScreen == models.BookDetailScreen {
			m.detail.ClearUpdated()
		}
		
	case models.BackupScreen:
		var backupModel tea.Model
		var backupCmd tea.Cmd
		// Update backup screen model
		backupModel, backupCmd = m.backup.Update(msg)
		m.backup = backupModel.(*screens.BackupScreen)
		cmd = backupCmd
		// Handle screen transitions from backup screen
		if switchMsg, ok := msg.(screens.SwitchScreenMsg); ok {
			newScreen = switchMsg.Screen
		} else {
			newScreen = m.currentScreen
		}
	}

	// Handle screen transitions and perform any necessary cleanup
	if newScreen != m.currentScreen {
		m.currentScreen = newScreen
		
		// Perform screen-specific cleanup when transitioning
		if newScreen == models.ListBooksScreen {
			// Clear any delete confirmation state when entering list screen
			m.listBooks.ClearDeleted()
		}
		if newScreen == models.BackupScreen {
			// Clear any backup status when entering backup screen
			m.backup.ClearStatus()
		}
	}

	// Return updated model and any command to execute
	return m, cmd
}

// View renders the current screen by delegating to the appropriate screen model
// It returns the string representation of the UI for the current screen
func (m Model) View() string {
	// Delegate view rendering to the current screen's model
	switch m.currentScreen {
	case models.MenuScreen:
		return m.menu.View()      // Render main menu
	case models.AddBookScreen:
		return m.addBook.View()   // Render add book form
	case models.ListBooksScreen:
		return m.listBooks.View() // Render book list
	case models.BookDetailScreen:
		return m.detail.View()    // Render book details
	case models.EditBookScreen:
		return m.edit.View()      // Render edit book form
	case models.BackupScreen:
		return m.backup.View()    // Render backup screen
	default:
		// Fallback for unknown screen states
		return ""
	}
}
