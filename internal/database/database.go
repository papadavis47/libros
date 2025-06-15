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
		notes TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.conn.Exec(createTable)
	return err
}

func (db *DB) SaveBook(title, author, notes string) error {
	title = strings.TrimSpace(title)
	author = strings.TrimSpace(author)
	notes = strings.TrimSpace(notes)

	if title == "" || author == "" {
		return fmt.Errorf("both title and author are required")
	}

	_, err := db.conn.Exec("INSERT INTO books (title, author, notes) VALUES (?, ?, ?)", title, author, notes)
	return err
}

func (db *DB) LoadBooks() ([]models.Book, error) {
	rows, err := db.conn.Query("SELECT id, title, author, notes, created_at, updated_at FROM books ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Notes, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (db *DB) UpdateBook(id int, title, author, notes string) error {
	title = strings.TrimSpace(title)
	author = strings.TrimSpace(author)
	notes = strings.TrimSpace(notes)

	if title == "" || author == "" {
		return fmt.Errorf("both title and author are required")
	}

	_, err := db.conn.Exec("UPDATE books SET title = ?, author = ?, notes = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", title, author, notes, id)
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
