package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/papadavis47/libros/internal/messages"
	"github.com/papadavis47/libros/internal/models"
	"github.com/papadavis47/libros/internal/styles"
)

type ListBooksModel struct {
	books   []models.Book
	index   int
	err     error
	deleted bool
}

func NewListBooksModel() ListBooksModel {
	return ListBooksModel{
		index: 0,
	}
}

func (m ListBooksModel) Update(msg tea.Msg) (ListBooksModel, tea.Cmd, models.Screen, *models.Book) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, nil, models.MenuScreen, nil
		case "up", "k":
			if m.index > 0 {
				m.index--
			}
		case "down", "j":
			if m.index < len(m.books)-1 {
				m.index++
			}
		case "enter":
			if len(m.books) > 0 {
				selectedBook := &m.books[m.index]
				return m, nil, models.BookDetailScreen, selectedBook
			}
		}

	case messages.LoadBooksMsg:
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.books = msg.Books
			if m.index >= len(m.books) && len(m.books) > 0 {
				m.index = len(m.books) - 1
			}
		}

	case messages.DeleteMsg:
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.deleted = true
		}
	}

	return m, nil, models.ListBooksScreen, nil
}

func (m ListBooksModel) View() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("📚 Libros - Davis Family Book Manager"))
	b.WriteString("\n")
	b.WriteString(styles.BlurredStyle.Render("Your Book Collection"))
	b.WriteString("\n\n")

	if len(m.books) == 0 {
		b.WriteString(styles.BlurredStyle.Render("No books found. Add some books first!"))
	} else {
		for i, book := range m.books {
			if i == m.index {
				b.WriteString(styles.SelectedStyle.Render(fmt.Sprintf("%s by %s", book.Title, book.Author)))
			} else {
				b.WriteString(styles.FocusedStyle.Render(book.Title))
				b.WriteString(" by ")
				b.WriteString(styles.BlurredStyle.Render(book.Author))
			}
			b.WriteString("\n\n")
		}
		
		b.WriteString(styles.BlurredStyle.Render(fmt.Sprintf("Total books: %d", len(m.books))))
		b.WriteString("\n")
	}

	if m.deleted {
		b.WriteString("\n")
		b.WriteString(styles.SuccessStyle.Render("✓ Book deleted successfully!"))
		b.WriteString("\n")
	}

	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(styles.ErrorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n")
	}

	if len(m.books) > 0 {
		b.WriteString("\n" + styles.BlurredStyle.Render("Use ↑/↓ or j/k to navigate, Enter to select, Esc to return to menu, q to quit"))
	} else {
		b.WriteString("\n" + styles.BlurredStyle.Render("Press Esc to return to menu, q or Ctrl+C to quit"))
	}

	return b.String()
}

func (m *ListBooksModel) ClearDeleted() {
	m.deleted = false
}
