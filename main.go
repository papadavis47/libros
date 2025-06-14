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

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)
)

type screen int

const (
	menuScreen screen = iota
	addBookScreen
	listBooksScreen
	bookDetailScreen
	editBookScreen
)

type book struct {
	ID     int
	Title  string
	Author string
}

type model struct {
	db            *sql.DB
	inputs        []textinput.Model
	focused       int
	err           error
	saved         bool
	deleted       bool
	updated       bool
	currentScreen screen
	menuItems     []string
	menuIndex     int
	books         []book
	selectedBook  *book
	bookListIndex int
	detailActions []string
	detailIndex   int
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
		db:            db,
		inputs:        make([]textinput.Model, 2),
		currentScreen: menuScreen,
		menuIndex:     0,
		detailActions: []string{"Edit Book", "Delete Book", "Back to List"},
		detailIndex:   0,
	}

	m.updateMenuItems()

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

func (m *model) updateMenuItems() {
	var count int
	err := m.db.QueryRow("SELECT COUNT(*) FROM books").Scan(&count)
	if err != nil {
		m.menuItems = []string{"Add Book", "Quit"}
		return
	}

	if count > 0 {
		m.menuItems = []string{"Add Book", "View Books", "Quit"}
	} else {
		m.menuItems = []string{"Add Book", "Quit"}
	}

	if m.menuIndex >= len(m.menuItems) {
		m.menuIndex = len(m.menuItems) - 1
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.db.Close()
			return m, tea.Quit
		}

		switch m.currentScreen {
		case menuScreen:
			return m.updateMenu(msg)
		case addBookScreen:
			return m.updateAddBook(msg)
		case listBooksScreen:
			return m.updateListBooks(msg)
		case bookDetailScreen:
			return m.updateBookDetail(msg)
		case editBookScreen:
			return m.updateEditBook(msg)
		}

	case saveMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.saved = true
			for i := range m.inputs {
				m.inputs[i].SetValue("")
			}
			m.focused = 0
			m.inputs[0].Focus()
			m.updateMenuItems()
		}
		return m, nil

	case updateMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.updated = true
			m.currentScreen = bookDetailScreen
		}
		return m, nil

	case deleteMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.deleted = true
			m.currentScreen = listBooksScreen
			m.updateMenuItems()
			return m, m.loadBooks()
		}
		return m, nil

	case loadBooksMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.books = msg.books
		}
		return m, nil
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.menuIndex > 0 {
			m.menuIndex--
		}
	case "down", "j":
		if m.menuIndex < len(m.menuItems)-1 {
			m.menuIndex++
		}
	case "enter":
		selectedItem := m.menuItems[m.menuIndex]
		switch selectedItem {
		case "Add Book":
			m.currentScreen = addBookScreen
			m.focused = 0
			m.err = nil
			m.saved = false
			m.inputs[0].Focus()
			for i := 1; i < len(m.inputs); i++ {
				m.inputs[i].Blur()
			}
		case "View Books":
			m.currentScreen = listBooksScreen
			return m, m.loadBooks()
		case "Quit":
			m.db.Close()
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) updateAddBook(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.currentScreen = menuScreen
		m.err = nil
		m.saved = false
		return m, nil
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

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m model) updateListBooks(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.currentScreen = menuScreen
		return m, nil
	case "up", "k":
		if m.bookListIndex > 0 {
			m.bookListIndex--
		}
	case "down", "j":
		if m.bookListIndex < len(m.books)-1 {
			m.bookListIndex++
		}
	case "enter":
		if len(m.books) > 0 {
			m.selectedBook = &m.books[m.bookListIndex]
			m.currentScreen = bookDetailScreen
			m.detailIndex = 0
			m.err = nil
			m.deleted = false
			m.updated = false
		}
	}
	return m, nil
}

func (m model) updateBookDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.currentScreen = listBooksScreen
		return m, nil
	case "up", "k":
		if m.detailIndex > 0 {
			m.detailIndex--
		}
	case "down", "j":
		if m.detailIndex < len(m.detailActions)-1 {
			m.detailIndex++
		}
	case "enter":
		selectedAction := m.detailActions[m.detailIndex]
		switch selectedAction {
		case "Edit Book":
			m.currentScreen = editBookScreen
			m.focused = 0
			m.err = nil
			m.updated = false
			m.inputs[0].SetValue(m.selectedBook.Title)
			m.inputs[1].SetValue(m.selectedBook.Author)
			m.inputs[0].Focus()
			for i := 1; i < len(m.inputs); i++ {
				m.inputs[i].Blur()
			}
		case "Delete Book":
			return m, m.deleteBook()
		case "Back to List":
			m.currentScreen = listBooksScreen
		}
	}
	return m, nil
}

func (m model) updateEditBook(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.currentScreen = bookDetailScreen
		m.err = nil
		m.updated = false
		return m, nil
	case "tab", "shift+tab", "enter", "up", "down":
		s := msg.String()

		if s == "enter" && m.focused == len(m.inputs)-1 {
			return m, m.updateBook()
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
	switch m.currentScreen {
	case menuScreen:
		return m.viewMenu()
	case addBookScreen:
		return m.viewAddBook()
	case listBooksScreen:
		return m.viewListBooks()
	case bookDetailScreen:
		return m.viewBookDetail()
	case editBookScreen:
		return m.viewEditBook()
	default:
		return ""
	}
}

func (m model) viewMenu() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("📚 Libros - Davis Family Book Manager"))
	b.WriteString("\n\n")

	for i, item := range m.menuItems {
		if i == m.menuIndex {
			b.WriteString(selectedStyle.Render("> " + item))
		} else {
			b.WriteString(blurredStyle.Render("  " + item))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n" + blurredStyle.Render("Use ↑/↓ or j/k to navigate, Enter to select, Ctrl+C to quit"))

	return b.String()
}

func (m model) viewAddBook() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("📚 Libros - Davis Family Book Manager"))
	b.WriteString("\n")
	b.WriteString(blurredStyle.Render("Add New Book"))
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

	b.WriteString(blurredStyle.Render("Press Esc to return to menu, Ctrl+C to quit"))

	return b.String()
}

func (m model) viewListBooks() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("📚 Libros - Davis Family Book Manager"))
	b.WriteString("\n")
	b.WriteString(blurredStyle.Render("Your Book Collection"))
	b.WriteString("\n\n")

	if len(m.books) == 0 {
		b.WriteString(blurredStyle.Render("No books found. Add some books first!"))
	} else {
		for i, book := range m.books {
			if i == m.bookListIndex {
				b.WriteString(selectedStyle.Render(fmt.Sprintf("> %d. %s by %s", i+1, book.Title, book.Author)))
			} else {
				b.WriteString(fmt.Sprintf("  %d. ", i+1))
				b.WriteString(focusedStyle.Render(book.Title))
				b.WriteString(" by ")
				b.WriteString(blurredStyle.Render(book.Author))
			}
			b.WriteString("\n")
		}
	}

	if m.deleted {
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("✓ Book deleted successfully!"))
		b.WriteString("\n")
	}

	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render("Error: " + m.err.Error()))
		b.WriteString("\n")
	}

	if len(m.books) > 0 {
		b.WriteString("\n" + blurredStyle.Render("Use ↑/↓ or j/k to navigate, Enter to select, Esc to return to menu"))
	} else {
		b.WriteString("\n" + blurredStyle.Render("Press Esc or q to return to menu, Ctrl+C to quit"))
	}

	return b.String()
}

func (m model) viewBookDetail() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("📚 Libros - Davis Family Book Manager"))
	b.WriteString("\n")
	b.WriteString(blurredStyle.Render("Book Details"))
	b.WriteString("\n\n")

	if m.selectedBook != nil {
		b.WriteString(focusedStyle.Render("Title: ") + m.selectedBook.Title + "\n")
		b.WriteString(focusedStyle.Render("Author: ") + m.selectedBook.Author + "\n\n")

		for i, action := range m.detailActions {
			if i == m.detailIndex {
				b.WriteString(selectedStyle.Render("> " + action))
			} else {
				b.WriteString(blurredStyle.Render("  " + action))
			}
			b.WriteString("\n")
		}
	}

	if m.updated {
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("✓ Book updated successfully!"))
		b.WriteString("\n")
	}

	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render("Error: " + m.err.Error()))
		b.WriteString("\n")
	}

	b.WriteString("\n" + blurredStyle.Render("Use ↑/↓ or j/k to navigate, Enter to select, Esc to go back"))

	return b.String()
}

func (m model) viewEditBook() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("📚 Libros - Davis Family Book Manager"))
	b.WriteString("\n")
	b.WriteString(blurredStyle.Render("Edit Book"))
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
	fmt.Fprintf(&b, "\n\n%s\n\n", button.Render("[ Update Book ]"))

	if m.err != nil {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render("Error: " + m.err.Error()))
		b.WriteString("\n")
	}

	b.WriteString(blurredStyle.Render("Press Esc to cancel, Ctrl+C to quit"))

	return b.String()
}

type saveMsg struct {
	err error
}

type updateMsg struct {
	err error
}

type deleteMsg struct {
	err error
}

type loadBooksMsg struct {
	books []book
	err   error
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

func (m model) updateBook() tea.Cmd {
	return func() tea.Msg {
		title := strings.TrimSpace(m.inputs[0].Value())
		author := strings.TrimSpace(m.inputs[1].Value())

		if title == "" || author == "" {
			return updateMsg{err: fmt.Errorf("both title and author are required")}
		}

		_, err := m.db.Exec("UPDATE books SET title = ?, author = ? WHERE id = ?", title, author, m.selectedBook.ID)
		if err != nil {
			return updateMsg{err: err}
		}

		m.selectedBook.Title = title
		m.selectedBook.Author = author

		return updateMsg{}
	}
}

func (m model) deleteBook() tea.Cmd {
	return func() tea.Msg {
		_, err := m.db.Exec("DELETE FROM books WHERE id = ?", m.selectedBook.ID)
		if err != nil {
			return deleteMsg{err: err}
		}

		return deleteMsg{}
	}
}

func (m model) loadBooks() tea.Cmd {
	return func() tea.Msg {
		rows, err := m.db.Query("SELECT id, title, author FROM books ORDER BY created_at DESC")
		if err != nil {
			return loadBooksMsg{err: err}
		}
		defer rows.Close()

		var books []book
		for rows.Next() {
			var b book
			err := rows.Scan(&b.ID, &b.Title, &b.Author)
			if err != nil {
				return loadBooksMsg{err: err}
			}
			books = append(books, b)
		}

		if err = rows.Err(); err != nil {
			return loadBooksMsg{err: err}
		}

		return loadBooksMsg{books: books}
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
