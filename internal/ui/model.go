package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/papadavis47/libros/internal/database"
	"github.com/papadavis47/libros/internal/models"
	"github.com/papadavis47/libros/internal/ui/screens"
)

type Model struct {
	db            *database.DB
	currentScreen models.Screen
	
	// Screen models
	menu      screens.MenuModel
	addBook   screens.AddBookModel
	listBooks screens.ListBooksModel
	detail    screens.DetailModel
	edit      screens.EditModel
}

func NewModel(db *database.DB) Model {
	return Model{
		db:            db,
		currentScreen: models.MenuScreen,
		menu:          screens.NewMenuModel(db),
		addBook:       screens.NewAddBookModel(db),
		listBooks:     screens.NewListBooksModel(),
		detail:        screens.NewDetailModel(db),
		edit:          screens.NewEditModel(db),
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.db.Close()
			return m, tea.Quit
		}
		if msg.String() == "q" && m.currentScreen != models.AddBookScreen && m.currentScreen != models.EditBookScreen {
			m.db.Close()
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	var newScreen models.Screen

	switch m.currentScreen {
	case models.MenuScreen:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			var menuCmd tea.Cmd
			m.menu, menuCmd, newScreen = m.menu.Update(keyMsg)
			cmd = menuCmd
		} else {
			newScreen = m.currentScreen
		}
		
	case models.AddBookScreen:
		var addBookCmd tea.Cmd
		m.addBook, addBookCmd, newScreen = m.addBook.Update(msg)
		cmd = addBookCmd
		if newScreen == models.MenuScreen {
			m.addBook.Reset()
			m.menu.RefreshItems()
		}
		
	case models.ListBooksScreen:
		var listCmd tea.Cmd
		var selectedBook *models.Book
		m.listBooks, listCmd, newScreen, selectedBook = m.listBooks.Update(msg)
		cmd = listCmd
		if selectedBook != nil {
			m.detail.SetBook(selectedBook)
		}
		if newScreen == models.MenuScreen {
			m.menu.RefreshItems()
		}
		
	case models.BookDetailScreen:
		var detailCmd tea.Cmd
		m.detail, detailCmd, newScreen = m.detail.Update(msg)
		cmd = detailCmd
		if newScreen == models.EditBookScreen {
			m.edit.SetBook(m.detail.SelectedBook)
		}
		
	case models.EditBookScreen:
		var editCmd tea.Cmd
		m.edit, editCmd, newScreen = m.edit.Update(msg)
		cmd = editCmd
		if newScreen == models.BookDetailScreen {
			m.detail.ClearUpdated()
		}
	}

	// Handle screen transitions
	if newScreen != m.currentScreen {
		m.currentScreen = newScreen
		
		// Clear state when transitioning
		if newScreen == models.ListBooksScreen {
			m.listBooks.ClearDeleted()
		}
	}

	return m, cmd
}

func (m Model) View() string {
	switch m.currentScreen {
	case models.MenuScreen:
		return m.menu.View()
	case models.AddBookScreen:
		return m.addBook.View()
	case models.ListBooksScreen:
		return m.listBooks.View()
	case models.BookDetailScreen:
		return m.detail.View()
	case models.EditBookScreen:
		return m.edit.View()
	default:
		return ""
	}
}
