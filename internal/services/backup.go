package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/papadavis47/libros/internal/constants"
	"github.com/papadavis47/libros/internal/interfaces"
	"github.com/papadavis47/libros/internal/models"
	"github.com/papadavis47/libros/internal/utils"
)

// BackupService implements the BackupService interface
type BackupService struct{}

// NewBackupService creates a new backup service instance
func NewBackupService() interfaces.BackupService {
	return &BackupService{}
}

// BackupData represents the structure of exported JSON data
type BackupData struct {
	ExportDate time.Time      `json:"export_date"`
	TotalBooks int            `json:"total_books"`
	Books      []models.Book  `json:"books"`
}

// ExportToJSON exports books to a JSON file
func (s *BackupService) ExportToJSON(books []models.Book, filePath string) error {
	// Create backup data structure
	backupData := BackupData{
		ExportDate: time.Now(),
		TotalBooks: len(books),
		Books:      books,
	}

	// Marshal to JSON with proper formatting
	jsonData, err := json.MarshalIndent(backupData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %v", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), constants.DirPermissions); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, jsonData, constants.FilePermissions); err != nil {
		return fmt.Errorf("failed to write JSON file: %v", err)
	}

	return nil
}

// ExportToMarkdown exports books to a Markdown file
func (s *BackupService) ExportToMarkdown(books []models.Book, filePath string) error {
	// Create markdown content
	md := fmt.Sprintf("# Book Collection Export\n\n")
	md += fmt.Sprintf("**Export Date:** %s  \n", time.Now().Format("January 2, 2006"))
	md += fmt.Sprintf("**Total Books:** %d  \n\n", len(books))

	// Add each book
	for i, book := range books {
		md += fmt.Sprintf("## %d. %s\n\n", i+1, book.Title)
		md += fmt.Sprintf("**Author:** %s  \n", book.Author)
		md += fmt.Sprintf("**Type:** %s  \n", utils.FormatBookType(book.Type))
		md += fmt.Sprintf("**Created:** %s  \n", utils.FormatDate(book.CreatedAt))
		md += fmt.Sprintf("**Updated:** %s  \n", utils.FormatDate(book.UpdatedAt))

		if book.Notes != "" {
			md += fmt.Sprintf("\n**Notes:**  \n%s\n", book.Notes)
		}
		md += "\n---\n\n"
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), constants.DirPermissions); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, []byte(md), constants.FilePermissions); err != nil {
		return fmt.Errorf("failed to write markdown file: %v", err)
	}

	return nil
}

// BackupDatabase creates a backup copy of the database file
func (s *BackupService) BackupDatabase(sourcePath, destPath string) error {
	// Read source file
	sourceFile, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source database: %v", err)
	}

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), constants.DirPermissions); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}

	// Write to destination
	if err := os.WriteFile(destPath, sourceFile, constants.FilePermissions); err != nil {
		return fmt.Errorf("failed to write backup file: %v", err)
	}

	return nil
}