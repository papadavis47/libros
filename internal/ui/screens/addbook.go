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

type AddBookModel struct {
	db      *database.DB
	inputs  []textinput.Model
	focused int
	err     error
	saved   bool
}

func NewAddBookModel(db *database.DB) AddBookModel {
	m := AddBookModel{
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

func (m AddBookModel) Update(msg tea.Msg) (AddBookModel, tea.Cmd, models.Screen) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.err = nil
			m.saved = false
			return m, nil, models.MenuScreen
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focused == len(m.inputs) {
				return m, m.saveBookCmd(), models.AddBookScreen
			}

			if s == "up" || s == "shift+tab" {
				m.focused--
			} else {
				m.focused++
			}

			if m.focused >= len(m.inputs)+1 {
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

			return m, tea.Batch(cmds...), models.AddBookScreen
		}

	case messages.SaveMsg:
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.saved = true
			for i := range m.inputs {
				m.inputs[i].SetValue("")
			}
			m.focused = 0
			m.inputs[0].Focus()
		}
		return m, nil, models.AddBookScreen
	}

	cmd := m.updateInputs(msg)
	return m, cmd, models.AddBookScreen
}

func (m *AddBookModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m AddBookModel) View() string {
	var b strings.Builder

	b.WriteString(styles.TitleStyle.Render("📚 Libros - Davis Family Book Manager"))
	b.WriteString("\n")
	b.WriteString(styles.BlurredStyle.Render("Add New Book"))
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
	fmt.Fprintf(&b, "\n\n%s\n\n", button.Render("SAVE BOOK"))

	if m.err != nil {
		b.WriteString(styles.ErrorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n")
	}

	if m.saved {
		b.WriteString(styles.SuccessStyle.Render("✓ Book saved successfully!"))
		b.WriteString("\n")
	}

	b.WriteString(styles.BlurredStyle.Render("Press Esc to return to menu, q or Ctrl+C to quit"))

	return b.String()
}

func (m AddBookModel) saveBookCmd() tea.Cmd {
	return func() tea.Msg {
		title := m.inputs[0].Value()
		author := m.inputs[1].Value()
		err := m.db.SaveBook(title, author)
		return messages.SaveMsg{Err: err}
	}
}

func (m *AddBookModel) Reset() {
	m.err = nil
	m.saved = false
	m.focused = 0
	for i := range m.inputs {
		m.inputs[i].SetValue("")
	}
	m.inputs[0].Focus()
	m.inputs[0].PromptStyle = styles.FocusedStyle
	m.inputs[0].TextStyle = styles.FocusedStyle
	for i := 1; i < len(m.inputs); i++ {
		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = styles.NoStyle
		m.inputs[i].TextStyle = styles.NoStyle
	}
}
