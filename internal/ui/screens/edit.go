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

type EditModel struct {
	db           *database.DB
	SelectedBook *models.Book
	inputs       []textinput.Model
	textarea     textarea.Model
	bookTypes    []models.BookType
	selectedType int
	focused      int
	err          error
}

func NewEditModel(db *database.DB) EditModel {
	m := EditModel{
		db:           db,
		inputs:       make([]textinput.Model, 2),
		bookTypes:    []models.BookType{models.Paperback, models.Hardback, models.Audio, models.Digital},
		selectedType: 0,
		focused:      0,
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

	// Initialize textarea for notes
	ta := textarea.New()
	ta.Placeholder = "Enter notes about this book (optional) . . ."
	ta.SetWidth(50)
	ta.SetHeight(4)
	ta.CharLimit = 1000
	m.textarea = ta

	return m
}

func (m EditModel) Update(msg tea.Msg) (EditModel, tea.Cmd, models.Screen) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.err = nil
			return m, nil, models.BookDetailScreen
		case "ctrl+a":
			if m.focused < len(m.inputs) {
				m.inputs[m.focused].CursorStart()
			}
			return m, nil, models.EditBookScreen
		case "ctrl+e":
			if m.focused < len(m.inputs) {
				m.inputs[m.focused].CursorEnd()
			}
			return m, nil, models.EditBookScreen
		case "tab", "shift+tab", "enter", "up", "down", "left", "right":
			s := msg.String()

			if s == "enter" && m.focused == len(m.inputs)+2 {
				return m, m.updateBookCmd(), models.EditBookScreen
			}

			// Handle book type selection
			if m.focused == len(m.inputs) && (s == "left" || s == "right") {
				if s == "left" {
					m.selectedType--
					if m.selectedType < 0 {
						m.selectedType = len(m.bookTypes) - 1
					}
				} else {
					m.selectedType++
					if m.selectedType >= len(m.bookTypes) {
						m.selectedType = 0
					}
				}
				return m, nil, models.EditBookScreen
			}

			if s == "up" || s == "shift+tab" {
				m.focused--
			} else {
				m.focused++
			}

			if m.focused > len(m.inputs)+2 {
				m.focused = 0
			} else if m.focused < 0 {
				m.focused = len(m.inputs) + 2
			}

			cmds := make([]tea.Cmd, len(m.inputs)+1)
			for i := 0; i < len(m.inputs); i++ {
				if i == m.focused {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = styles.FocusedStyle
					m.inputs[i].TextStyle = styles.FocusedStyle
					m.inputs[i].CursorEnd()
				} else {
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = styles.NoStyle
					m.inputs[i].TextStyle = styles.NoStyle
				}
			}

			// Handle textarea focus
			if m.focused == len(m.inputs)+1 {
				cmds[len(m.inputs)] = m.textarea.Focus()
			} else {
				m.textarea.Blur()
			}

			return m, tea.Batch(cmds...), models.EditBookScreen
		}

	case messages.UpdateMsg:
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			m.SelectedBook.Title = m.inputs[0].Value()
			m.SelectedBook.Author = m.inputs[1].Value()
			m.SelectedBook.Type = m.bookTypes[m.selectedType]
			m.SelectedBook.Notes = m.textarea.Value()
			return m, nil, models.BookDetailScreen
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd, models.EditBookScreen
}

func (m *EditModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs)+1)

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	cmds[len(m.inputs)] = cmd

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
		b.WriteRune('\n')
	}

	// Add book type selector
	b.WriteString("\n")
	typeLabel := "Type:   "
	if m.focused == len(m.inputs) {
		b.WriteString(styles.FocusedStyle.Render(typeLabel))
	} else {
		b.WriteString(styles.BlurredStyle.Render(typeLabel))
	}
	
	for i, bookType := range m.bookTypes {
		if i == m.selectedType {
			if m.focused == len(m.inputs) {
				b.WriteString(styles.ButtonStyle.Render(fmt.Sprintf(" %s ", string(bookType))))
			} else {
				b.WriteString(styles.FocusedStyle.Render(fmt.Sprintf(" %s ", string(bookType))))
			}
		} else {
			b.WriteString(styles.BlurredStyle.Render(fmt.Sprintf(" %s ", string(bookType))))
		}
		if i < len(m.bookTypes)-1 {
			b.WriteString(" ")
		}
	}
	b.WriteString("\n")

	// Add notes textarea
	b.WriteString("\n")
	b.WriteString(styles.FocusedStyle.Render("Notes: "))
	b.WriteString("\n")
	b.WriteString(m.textarea.View())

	button := &styles.BlurredStyle
	if m.focused == len(m.inputs)+2 {
		button = &styles.ButtonStyle
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", button.Render("UPDATE BOOK"))

	if m.err != nil {
		b.WriteString(styles.ErrorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n")
	}

	b.WriteString(styles.BlurredStyle.Render("Press Esc to cancel, Ctrl+A/Ctrl+E for start/end of field, q or Ctrl+C to quit"))

	return b.String()
}

func (m *EditModel) SetBook(book *models.Book) {
	m.SelectedBook = book
	m.focused = 0
	m.err = nil
	m.inputs[0].SetValue(book.Title)
	m.inputs[1].SetValue(book.Author)
	m.textarea.SetValue(book.Notes)
	
	// Set the selected type based on the book's type
	for i, bookType := range m.bookTypes {
		if bookType == book.Type {
			m.selectedType = i
			break
		}
	}
	
	m.inputs[0].Focus()
	m.inputs[0].CursorEnd()
	for i := 1; i < len(m.inputs); i++ {
		m.inputs[i].Blur()
	}
	m.textarea.Blur()
}

func (m EditModel) updateBookCmd() tea.Cmd {
	return func() tea.Msg {
		title := m.inputs[0].Value()
		author := m.inputs[1].Value()
		bookType := m.bookTypes[m.selectedType]
		notes := m.textarea.Value()
		err := m.db.UpdateBook(m.SelectedBook.ID, title, author, bookType, notes)
		return messages.UpdateMsg{Err: err}
	}
}
