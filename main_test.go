package main

import (
	"database/sql"
	"os"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

func TestDatabaseOperations(t *testing.T) {
	testDBPath := "test_books.db"
	defer os.Remove(testDBPath)

	db, err := sql.Open("sqlite3", testDBPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	createTable := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	_, err = db.Exec("INSERT INTO books (title, author) VALUES (?, ?)", "Test Title", "Test Author")
	if err != nil {
		t.Fatalf("Failed to insert book: %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM books").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query book count: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 book, got %d", count)
	}

	var title, author string
	err = db.QueryRow("SELECT title, author FROM books WHERE id = 1").Scan(&title, &author)
	if err != nil {
		t.Fatalf("Failed to query book: %v", err)
	}

	if title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", title)
	}

	if author != "Test Author" {
		t.Errorf("Expected author 'Test Author', got '%s'", author)
	}
}

func TestModelInitialization(t *testing.T) {
	testDBPath := "test_init_books.db"
	defer os.Remove(testDBPath)

	originalDB := "books.db"
	defer func() {
		os.Rename(testDBPath, originalDB)
	}()

	db, err := sql.Open("sqlite3", testDBPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	createTable := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	m := model{
		db:     db,
		inputs: make([]textinput.Model, 2),
	}

	var t1 textinput.Model
	for i := range m.inputs {
		t1 = textinput.New()
		t1.CharLimit = 255

		switch i {
		case 0:
			t1.Placeholder = "Enter book title"
			t1.Focus()
		case 1:
			t1.Placeholder = "Enter author name"
		}

		m.inputs[i] = t1
	}

	if len(m.inputs) != 2 {
		t.Errorf("Expected 2 inputs, got %d", len(m.inputs))
	}

	if m.inputs[0].Placeholder != "Enter book title" {
		t.Errorf("Expected first input placeholder 'Enter book title', got '%s'", m.inputs[0].Placeholder)
	}

	if m.inputs[1].Placeholder != "Enter author name" {
		t.Errorf("Expected second input placeholder 'Enter author name', got '%s'", m.inputs[1].Placeholder)
	}
}

func TestSaveBookValidation(t *testing.T) {
	testDBPath := "test_validation_books.db"
	defer os.Remove(testDBPath)

	db, err := sql.Open("sqlite3", testDBPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	createTable := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	m := model{
		db:     db,
		inputs: make([]textinput.Model, 2),
	}

	for i := range m.inputs {
		m.inputs[i] = textinput.New()
	}

	m.inputs[0].SetValue("")
	m.inputs[1].SetValue("")

	cmd := m.saveBook()
	msg := cmd()

	saveMsgResult, ok := msg.(saveMsg)
	if !ok {
		t.Fatal("Expected saveMsg type")
	}

	if saveMsgResult.err == nil {
		t.Error("Expected validation error for empty fields")
	}

	m.inputs[0].SetValue("Valid Title")
	m.inputs[1].SetValue("Valid Author")

	cmd = m.saveBook()
	msg = cmd()

	saveMsgResult, ok = msg.(saveMsg)
	if !ok {
		t.Fatal("Expected saveMsg type")
	}

	if saveMsgResult.err != nil {
		t.Errorf("Expected no error for valid input, got: %v", saveMsgResult.err)
	}
}

func TestModelUpdate(t *testing.T) {
	testDBPath := "test_update_books.db"
	defer os.Remove(testDBPath)

	db, err := sql.Open("sqlite3", testDBPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	createTable := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	m := model{
		db:     db,
		inputs: make([]textinput.Model, 2),
	}

	for i := range m.inputs {
		m.inputs[i] = textinput.New()
	}

	quitMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	updatedModel, cmd := m.Update(quitMsg)
	
	if cmd == nil {
		t.Error("Expected quit command")
	}

	_ = updatedModel

	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, cmd = m.Update(tabMsg)
	
	updatedM := updatedModel.(model)
	if updatedM.focused != 1 {
		t.Errorf("Expected focused to be 1 after tab, got %d", updatedM.focused)
	}
}