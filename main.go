package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/mattn/go-sqlite3"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4"))
	
	focusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4"))
	
	blurredStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF"))
	
	noStyle = lipgloss.NewStyle()
)

type model struct {
	db       *sql.DB
	inputs   []textinput.Model
	focused  int
	err      error
	saved    bool
}

func initialModel() model {
	db, err := sql.Open("sqlite3", "books.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	m := model{
		db:     db,
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CharLimit = 255

		switch i {
		case 0:
			t.Placeholder = "Enter book title"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Enter author name"
		}

		m.inputs[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.db.Close()
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focused == len(m.inputs)-1 {
				return m, m.saveBook()
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
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focused {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}

	case saveMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.saved = true
		}
		return m, nil
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("📚 Book Entry"))
	b.WriteString("\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredStyle
	if m.focused == len(m.inputs) {
		button = &focusedStyle
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", button.Render("[ Save Book ]"))

	if m.err != nil {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render("Error: " + m.err.Error()))
		b.WriteString("\n")
	}

	if m.saved {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("✓ Book saved successfully!"))
		b.WriteString("\n")
	}

	b.WriteString(blurredStyle.Render("Press ctrl+c to quit"))

	return b.String()
}

type saveMsg struct {
	err error
}

func (m model) saveBook() tea.Cmd {
	return func() tea.Msg {
		title := strings.TrimSpace(m.inputs[0].Value())
		author := strings.TrimSpace(m.inputs[1].Value())

		if title == "" || author == "" {
			return saveMsg{err: fmt.Errorf("both title and author are required")}
		}

		_, err := m.db.Exec("INSERT INTO books (title, author) VALUES (?, ?)", title, author)
		if err != nil {
			return saveMsg{err: err}
		}

		return saveMsg{}
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
