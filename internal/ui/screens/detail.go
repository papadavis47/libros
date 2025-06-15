package screens

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/papadavis47/libros/internal/database"
	"github.com/papadavis47/libros/internal/messages"
	"github.com/papadavis47/libros/internal/models"
	"github.com/papadavis47/libros/internal/styles"
)

type DetailModel struct {
	db           *database.DB
	SelectedBook *models.Book
	actions      []string
	index        int
	err          error
	updated      bool
}

func NewDetailModel(db *database.DB) DetailModel {
	return DetailModel{
		db:      db,
		actions: []string{"Edit Book", "Delete Book", "Back to List"},
		index:   0,
	}
}

func (m DetailModel) Update(msg tea.Msg) (DetailModel, tea.Cmd, models.Screen) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, nil, models.ListBooksScreen
		case "up", "k":
			if m.index > 0 {
				m.index--
			}
		case "down", "j":
			if m.index < len(m.actions)-1 {
				m.index++
			}
		case "enter":
			selectedAction := m.actions[m.index]
			switch selectedAction {
			case "Edit Book":
				return m, nil, models.EditBookScreen
			case "Delete Book":
				return m, m.deleteBookCmd(), models.BookDetailScreen
			case "Back to List":
				return m, nil, models.ListBooksScreen
			}
		}

	case messages.UpdateMsg:
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.updated = true
		}

	case messages.DeleteMsg:
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			return m, m.loadBooksCmd(), models.ListBooksScreen
		}
	}

	return m, nil, models.BookDetailScreen
}

func (m DetailModel) View() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("📚 Libros - Davis Family Book Manager"))
	b.WriteString("\n")
	b.WriteString(styles.BlurredStyle.Render("Book Details"))
	b.WriteString("\n\n")

	if m.SelectedBook != nil {
		b.WriteString(styles.FocusedStyle.Render("Title: ") + m.SelectedBook.Title + "\n")
		b.WriteString(styles.FocusedStyle.Render("Author: ") + m.SelectedBook.Author + "\n\n")

		for i, action := range m.actions {
			if i == m.index {
				b.WriteString(styles.SelectedStyle.Render(action))
			} else {
				b.WriteString(styles.BlurredStyle.Render(action))
			}
			b.WriteString("\n\n")
		}
	}

	if m.updated {
		b.WriteString("\n")
		b.WriteString(styles.SuccessStyle.Render("✓ Book updated successfully!"))
		b.WriteString("\n")
	}

	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(styles.ErrorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n")
	}

	b.WriteString("\n" + styles.BlurredStyle.Render("Use ↑/↓ or j/k to navigate, Enter to select, Esc to go back, q to quit"))

	return b.String()
}

func (m *DetailModel) SetBook(book *models.Book) {
	m.SelectedBook = book
	m.index = 0
	m.err = nil
	m.updated = false
}

func (m DetailModel) deleteBookCmd() tea.Cmd {
	return func() tea.Msg {
		err := m.db.DeleteBook(m.SelectedBook.ID)
		return messages.DeleteMsg{Err: err}
	}
}

func (m *DetailModel) ClearUpdated() {
	m.updated = false
}

func (m DetailModel) loadBooksCmd() tea.Cmd {
	return func() tea.Msg {
		books, err := m.db.LoadBooks()
		return messages.LoadBooksMsg{Books: books, Err: err}
	}
}
