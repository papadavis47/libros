package screens

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/papadavis47/libros/internal/database"
	"github.com/papadavis47/libros/internal/messages"
	"github.com/papadavis47/libros/internal/models"
	"github.com/papadavis47/libros/internal/styles"
)

type MenuModel struct {
	db        *database.DB
	items     []string
	index     int
}

func NewMenuModel(db *database.DB) MenuModel {
	m := MenuModel{
		db:    db,
		index: 0,
	}
	m.updateMenuItems()
	return m
}

func (m *MenuModel) updateMenuItems() {
	count, err := m.db.GetBookCount()
	if err != nil {
		m.items = []string{"Add Book", "Quit"}
		return
	}

	if count > 0 {
		m.items = []string{"Add Book", "View Books", "Quit"}
	} else {
		m.items = []string{"Add Book", "Quit"}
	}

	if m.index >= len(m.items) {
		m.index = len(m.items) - 1
	}
}

func (m MenuModel) Update(msg tea.KeyMsg) (MenuModel, tea.Cmd, models.Screen) {
	switch msg.String() {
	case "up", "k":
		if m.index > 0 {
			m.index--
		}
	case "down", "j":
		if m.index < len(m.items)-1 {
			m.index++
		}
	case "enter":
		selectedItem := m.items[m.index]
		switch selectedItem {
		case "Add Book":
			return m, nil, models.AddBookScreen
		case "View Books":
			return m, m.LoadBooksCmd(), models.ListBooksScreen
		case "Quit":
			return m, tea.Quit, models.MenuScreen
		}
	}
	return m, nil, models.MenuScreen
}

func (m MenuModel) View() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("📚 Libros - Davis Family Book Manager"))
	b.WriteString("\n\n")

	for i, item := range m.items {
		if i == m.index {
			b.WriteString(styles.SelectedStyle.Render(item))
		} else {
			b.WriteString(styles.BlurredStyle.Render(item))
		}
		b.WriteString("\n\n")
	}

	b.WriteString("\n" + styles.BlurredStyle.Render("Use ↑/↓ or j/k to navigate, Enter to select, q or Ctrl+C to quit"))

	return b.String()
}

func (m *MenuModel) RefreshItems() {
	m.updateMenuItems()
}

func (m MenuModel) LoadBooksCmd() tea.Cmd {
	return func() tea.Msg {
		books, err := m.db.LoadBooks()
		return messages.LoadBooksMsg{Books: books, Err: err}
	}
}
