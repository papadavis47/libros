package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/papadavis47/libros/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

func New(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	db := &DB{conn: conn}
	if err := db.createTable(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) createTable() error {
	createTable := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		type TEXT NOT NULL DEFAULT 'paperback',
		notes TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.conn.Exec(createTable)
	if err != nil {
		return err
	}

	// Add type column to existing tables if it doesn't exist
	_, err = db.conn.Exec("ALTER TABLE books ADD COLUMN type TEXT NOT NULL DEFAULT 'paperback'")
	if err != nil && !strings.Contains(err.Error(), "duplicate column name") {
		return err
	}

	return nil
}

func (db *DB) SaveBook(title, author string, bookType models.BookType, notes string) error {
	title = strings.TrimSpace(title)
	author = strings.TrimSpace(author)
	notes = strings.TrimSpace(notes)

	if title == "" || author == "" {
		return fmt.Errorf("both title and author are required")
	}

	_, err := db.conn.Exec("INSERT INTO books (title, author, type, notes) VALUES (?, ?, ?, ?)", title, author, string(bookType), notes)
	return err
}

func (db *DB) LoadBooks() ([]models.Book, error) {
	rows, err := db.conn.Query("SELECT id, title, author, type, notes, created_at, updated_at FROM books ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		var bookType string
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &bookType, &b.Notes, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		b.Type = models.BookType(bookType)
		books = append(books, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (db *DB) UpdateBook(id int, title, author string, bookType models.BookType, notes string) error {
	title = strings.TrimSpace(title)
	author = strings.TrimSpace(author)
	notes = strings.TrimSpace(notes)

	if title == "" || author == "" {
		return fmt.Errorf("both title and author are required")
	}

	_, err := db.conn.Exec("UPDATE books SET title = ?, author = ?, type = ?, notes = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", title, author, string(bookType), notes, id)
	return err
}

func (db *DB) DeleteBook(id int) error {
	_, err := db.conn.Exec("DELETE FROM books WHERE id = ?", id)
	return err
}

func (db *DB) GetBookCount() (int, error) {
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM books").Scan(&count)
	return count, err
}
