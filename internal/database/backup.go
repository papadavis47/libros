package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/papadavis47/libros/internal/models"
)

type BackupData struct {
	ExportDate time.Time      `json:"export_date"`
	TotalBooks int            `json:"total_books"`
	Books      []models.Book  `json:"books"`
}

func (db *DB) BackupToJSON(backupDir string) error {
	books, err := db.LoadBooks()
	if err != nil {
		return fmt.Errorf("failed to load books: %w", err)
	}

	backupData := BackupData{
		ExportDate: time.Now(),
		TotalBooks: len(books),
		Books:      books,
	}

	jsonData, err := json.MarshalIndent(backupData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	filename := "books.json"
	filePath := filepath.Join(backupDir, filename)

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}

func (db *DB) BackupToMarkdown(backupDir string) error {
	books, err := db.LoadBooks()
	if err != nil {
		return fmt.Errorf("failed to load books: %w", err)
	}

	var md string
	md += "# Libros - My Book Collection\n\n"

	for _, book := range books {
		md += fmt.Sprintf("## %s\n", book.Title)
		md += fmt.Sprintf("**Author:** %s  \n", book.Author)
		md += fmt.Sprintf("**Type:** %s  \n", formatBookType(book.Type))
		md += fmt.Sprintf("**Added:** %s  \n", book.CreatedAt.Format("January 2, 2006"))
		if !book.UpdatedAt.Equal(book.CreatedAt) {
			md += fmt.Sprintf("**Updated:** %s  \n", book.UpdatedAt.Format("January 2, 2006"))
		}
		md += "\n"

		if book.Notes != "" {
			md += fmt.Sprintf("*%s*\n", book.Notes)
		}
		md += "\n"
	}

	filename := "books.md"
	filePath := filepath.Join(backupDir, filename)

	err = os.WriteFile(filePath, []byte(md), 0644)
	if err != nil {
		return fmt.Errorf("failed to write Markdown file: %w", err)
	}

	return nil
}

func formatBookType(bookType models.BookType) string {
	switch bookType {
	case models.Paperback:
		return "Paperback"
	case models.Hardback:
		return "Hardback"
	case models.Audio:
		return "Audio"
	case models.Digital:
		return "Digital"
	default:
		return string(bookType)
	}
}