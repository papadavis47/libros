// Package database provides SQLite database operations for the libros book management application.
// It handles all database interactions including table creation, CRUD operations, and connection management.
package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3" // SQLite driver for database/sql
	"github.com/papadavis47/libros/internal/models"
)

// DB wraps a SQL database connection and provides methods for book management operations.
type DB struct {
	conn *sql.DB // SQLite database connection
}

// New creates a new database connection and initializes the books table.
// It takes a database file path and returns a DB instance or an error.
func New(dbPath string) (*DB, error) {
	// Open SQLite database connection
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create DB instance and initialize table schema
	db := &DB{conn: conn}
	if err := db.createTable(); err != nil {
		return nil, err
	}

	return db, nil
}

// Close closes the database connection and releases resources.
func (db *DB) Close() error {
	return db.conn.Close()
}

// createTable creates the books table if it doesn't exist and handles schema migrations.
// It ensures backward compatibility by adding missing columns to existing tables.
func (db *DB) createTable() error {
	// SQL statement to create books table with all required columns
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

	// Execute table creation
	_, err := db.conn.Exec(createTable)
	if err != nil {
		return err
	}

	// Handle schema migration: add type column to existing tables
	// This ensures backward compatibility with databases created before the type column was added
	_, err = db.conn.Exec("ALTER TABLE books ADD COLUMN type TEXT NOT NULL DEFAULT 'paperback'")
	if err != nil && !strings.Contains(err.Error(), "duplicate column name") {
		return err
	}

	return nil
}

// SaveBook inserts a new book record into the database.
// It validates required fields and trims whitespace from input values.
func (db *DB) SaveBook(title, author string, bookType models.BookType, notes string) error {
	// Sanitize input by trimming whitespace
	title = strings.TrimSpace(title)
	author = strings.TrimSpace(author)
	notes = strings.TrimSpace(notes)

	// Validate required fields
	if title == "" || author == "" {
		return fmt.Errorf("title, author, and type are required")
	}

	// Insert book record using parameterized query to prevent SQL injection
	_, err := db.conn.Exec("INSERT INTO books (title, author, type, notes) VALUES (?, ?, ?, ?)", title, author, string(bookType), notes)
	return err
}

// LoadBooks retrieves all books from the database ordered by creation date (newest first).
// It returns a slice of Book models or an error if the query fails.
func (db *DB) LoadBooks() ([]models.Book, error) {
	// Query all books with ordering by creation date
	rows, err := db.conn.Query("SELECT id, title, author, type, notes, created_at, updated_at FROM books ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process query results
	var books []models.Book
	for rows.Next() {
		var b models.Book
		var bookType string
		// Scan row data into book struct
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &bookType, &b.Notes, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		// Convert string type to BookType enum
		b.Type = models.BookType(bookType)
		books = append(books, b)
	}

	// Check for errors that occurred during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

// UpdateBook modifies an existing book record in the database.
// It validates input fields and updates the record's timestamp.
func (db *DB) UpdateBook(id int, title, author string, bookType models.BookType, notes string) error {
	// Sanitize input by trimming whitespace
	title = strings.TrimSpace(title)
	author = strings.TrimSpace(author)
	notes = strings.TrimSpace(notes)

	// Validate required fields
	if title == "" || author == "" {
		return fmt.Errorf("title, author, and type are required")
	}

	// Update book record and set updated_at timestamp
	_, err := db.conn.Exec("UPDATE books SET title = ?, author = ?, type = ?, notes = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", title, author, string(bookType), notes, id)
	return err
}

// DeleteBook removes a book from the database by its ID
// Takes the book ID as parameter and permanently deletes the record
func (db *DB) DeleteBook(id int) error {
	// Execute DELETE statement using parameterized query to prevent SQL injection
	// The ? parameter in the sql is what paramterizes this code
	_, err := db.conn.Exec("DELETE FROM books WHERE id = ?", id)
	return err
}

// GetBookCount returns the total number of books in the database.
// It executes a COUNT query and returns the result or an error.
func (db *DB) GetBookCount() (int, error) {
	var count int
	// Query total number of books in the database
	err := db.conn.QueryRow("SELECT COUNT(*) FROM books").Scan(&count)
	return count, err
}
