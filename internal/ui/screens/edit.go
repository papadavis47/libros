package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/papadavis47/libros/internal/database"
	"github.com/papadavis47/libros/internal/messages"
	"github.com/papadavis47/libros/internal/models"
	"github.com/papadavis47/libros/internal/styles"
)

type EditModel struct {
	db           *database.DB
	SelectedBook *models.Book
	inputs       []textinput.Model
	focused      int
	err          error
}

func NewEditModel(db *database.DB) EditModel {
	m := EditModel{
		db:      db,
		inputs:  make([]textinput.Model, 2),
		focused: 0,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CharLimit = 255
		t.Width = 50

		switch i {
		case 0:
			t.Placeholder = "Enter book title"
			t.Prompt = "Title:  "
			t.Focus()
			t.PromptStyle = styles.FocusedStyle
			t.TextStyle = styles.FocusedStyle
		case 1:
			t.Placeholder = "Enter author name"
			t.Prompt = "Author: "
		}

		m.inputs[i] = t
	}

	return m
}

func (m EditModel) Update(msg tea.Msg) (EditModel, tea.Cmd, models.Screen) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.err = nil
			return m, nil, models.BookDetailScreen
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focused == len(m.inputs) {
				return m, m.updateBookCmd(), models.EditBookScreen
			}

			if s == "up" || s == "shift+tab" {
				m.focused--
			} else {
				m.focused++
			}

			if m.focused > len(m.inputs) {
				m.focused = 0
			} else if m.focused < 0 {
				m.focused = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i < len(m.inputs); i++ {
				if i == m.focused {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = styles.FocusedStyle
					m.inputs[i].TextStyle = styles.FocusedStyle
				} else {
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = styles.NoStyle
					m.inputs[i].TextStyle = styles.NoStyle
				}
			}

			return m, tea.Batch(cmds...), models.EditBookScreen
		}

	case messages.UpdateMsg:
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.SelectedBook.Title = m.inputs[0].Value()
			m.SelectedBook.Author = m.inputs[1].Value()
			return m, nil, models.BookDetailScreen
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd, models.EditBookScreen
}

func (m *EditModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m EditModel) View() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("📚 Libros - Davis Family Book Manager"))
	b.WriteString("\n")
	b.WriteString(styles.BlurredStyle.Render("Edit Book"))
	b.WriteString("\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &styles.BlurredStyle
	if m.focused == len(m.inputs) {
		button = &styles.ButtonStyle
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", button.Render("UPDATE BOOK"))

	if m.err != nil {
		b.WriteString(styles.ErrorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n")
	}

	b.WriteString(styles.BlurredStyle.Render("Press Esc to cancel, q or Ctrl+C to quit"))

	return b.String()
}

func (m *EditModel) SetBook(book *models.Book) {
	m.SelectedBook = book
	m.focused = 0
	m.err = nil
	m.inputs[0].SetValue(book.Title)
	m.inputs[1].SetValue(book.Author)
	m.inputs[0].Focus()
	for i := 1; i < len(m.inputs); i++ {
		m.inputs[i].Blur()
	}
}

func (m EditModel) updateBookCmd() tea.Cmd {
	return func() tea.Msg {
		title := m.inputs[0].Value()
		author := m.inputs[1].Value()
		err := m.db.UpdateBook(m.SelectedBook.ID, title, author)
		return messages.UpdateMsg{Err: err}
	}
}
