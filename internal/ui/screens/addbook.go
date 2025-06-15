package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/papadavis47/libros/internal/database"
	"github.com/papadavis47/libros/internal/messages"
	"github.com/papadavis47/libros/internal/models"
	"github.com/papadavis47/libros/internal/styles"
)

type AddBookModel struct {
	db       *database.DB
	inputs   []textinput.Model
	textarea textarea.Model
	focused  int
	err      error
	saved    bool
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
			t.Placeholder = "Book title here . . ."
			t.Prompt = "Title:  "
			t.Focus()
			t.PromptStyle = styles.FocusedStyle
			t.TextStyle = styles.FocusedStyle
		case 1:
			t.Placeholder = "Author name here . . ."
			t.Prompt = "Author: "
		}

		m.inputs[i] = t
	}

	// Initialize textarea for notes
	ta := textarea.New()
	ta.Placeholder = "Notes about this book (optional) . . ."
	ta.SetWidth(50)
	ta.SetHeight(4)
	ta.CharLimit = 1000
	m.textarea = ta

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

			if s == "enter" && m.focused == len(m.inputs)+1 {
				return m, m.saveBookCmd(), models.AddBookScreen
			}

			if s == "up" || s == "shift+tab" {
				m.focused--
			} else {
				m.focused++
			}

			if m.focused >= len(m.inputs)+2 {
				m.focused = 0
			} else if m.focused < 0 {
				m.focused = len(m.inputs) + 1
			}

			cmds := make([]tea.Cmd, len(m.inputs)+1)
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

			// Handle textarea focus
			if m.focused == len(m.inputs) {
				cmds[len(m.inputs)] = m.textarea.Focus()
			} else {
				m.textarea.Blur()
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
			m.textarea.SetValue("")
			m.focused = 0
			m.inputs[0].Focus()
		}
		return m, nil, models.AddBookScreen
	}

	cmd := m.updateInputs(msg)
	return m, cmd, models.AddBookScreen
}

func (m *AddBookModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs)+1)

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	cmds[len(m.inputs)] = cmd

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
		b.WriteRune('\n')
	}

	// Add notes textarea
	b.WriteString("\n")
	b.WriteString(styles.FocusedStyle.Render("Notes: "))
	b.WriteString("\n")
	b.WriteString(m.textarea.View())

	button := &styles.BlurredStyle
	if m.focused == len(m.inputs)+1 {
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
		notes := m.textarea.Value()
		err := m.db.SaveBook(title, author, notes)
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
	m.textarea.SetValue("")
	m.inputs[0].Focus()
	m.inputs[0].PromptStyle = styles.FocusedStyle
	m.inputs[0].TextStyle = styles.FocusedStyle
	for i := 1; i < len(m.inputs); i++ {
		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = styles.NoStyle
		m.inputs[i].TextStyle = styles.NoStyle
	}
	m.textarea.Blur()
}
