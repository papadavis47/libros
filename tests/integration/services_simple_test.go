package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/papadavis47/libros/internal/models"
	"github.com/papadavis47/libros/internal/services"
)

// TestBackupService_ExportToJSON tests JSON export functionality
func TestBackupService_ExportToJSON(t *testing.T) {
	// Create temporary directory for export testing
	tempDir, err := os.MkdirTemp("", "libros_test_json")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	service := services.NewBackupService()

	// Create test data
	testBooks := []models.Book{
		{
			ID:        1,
			Title:     "The Go Programming Language",
			Author:    "Alan Donovan",
			Type:      models.Paperback,
			Notes:     "Excellent reference book for Go developers",
			CreatedAt: time.Date(2023, 1, 15, 10, 30, 0, 0, time.UTC),
			UpdatedAt: time.Date(2023, 1, 20, 14, 45, 0, 0, time.UTC),
		},
		{
			ID:        2,
			Title:     "Clean Code",
			Author:    "Robert C. Martin",
			Type:      models.Hardback,
			Notes:     "Essential reading for software craftsmanship",
			CreatedAt: time.Date(2023, 2, 10, 9, 15, 0, 0, time.UTC),
			UpdatedAt: time.Date(2023, 2, 10, 9, 15, 0, 0, time.UTC),
		},
	}

	// Test successful export
	t.Run("SuccessfulExport", func(t *testing.T) {
		exportPath := filepath.Join(tempDir, "books_test.json")

		err := service.ExportToJSON(testBooks, exportPath)
		if err != nil {
			t.Fatalf("ExportToJSON failed: %v", err)
		}

		// Verify file was created
		if _, err := os.Stat(exportPath); os.IsNotExist(err) {
			t.Fatal("Export file was not created")
		}

		// Read and parse the exported JSON
		content, err := os.ReadFile(exportPath)
		if err != nil {
			t.Fatalf("Failed to read exported file: %v", err)
		}

		var backupData struct {
			ExportDate time.Time     `json:"export_date"`
			TotalBooks int           `json:"total_books"`
			Books      []models.Book `json:"books"`
		}

		err = json.Unmarshal(content, &backupData)
		if err != nil {
			t.Fatalf("Failed to parse exported JSON: %v", err)
		}

		// Verify metadata
		if backupData.TotalBooks != len(testBooks) {
			t.Errorf("TotalBooks = %d, want %d", backupData.TotalBooks, len(testBooks))
		}

		if len(backupData.Books) != len(testBooks) {
			t.Errorf("Exported books count = %d, want %d", len(backupData.Books), len(testBooks))
		}

		// Verify export date is recent
		if time.Since(backupData.ExportDate) > time.Minute {
			t.Error("Export date should be recent")
		}

		// Verify first book data
		if len(backupData.Books) > 0 {
			book := backupData.Books[0]
			expected := testBooks[0]
			if book.Title != expected.Title {
				t.Errorf("Book Title = %q, want %q", book.Title, expected.Title)
			}
			if book.Author != expected.Author {
				t.Errorf("Book Author = %q, want %q", book.Author, expected.Author)
			}
		}
	})

	// Test export with empty book list
	t.Run("ExportEmptyList", func(t *testing.T) {
		emptyBooks := []models.Book{}
		exportPath := filepath.Join(tempDir, "empty_books.json")

		err := service.ExportToJSON(emptyBooks, exportPath)
		if err != nil {
			t.Fatalf("ExportToJSON with empty list failed: %v", err)
		}

		// Verify file was created and contains valid JSON
		content, err := os.ReadFile(exportPath)
		if err != nil {
			t.Fatalf("Failed to read empty export file: %v", err)
		}

		var backupData struct {
			TotalBooks int `json:"total_books"`
		}

		err = json.Unmarshal(content, &backupData)
		if err != nil {
			t.Fatalf("Failed to parse empty export JSON: %v", err)
		}

		if backupData.TotalBooks != 0 {
			t.Errorf("Empty export TotalBooks = %d, want 0", backupData.TotalBooks)
		}
	})
}

// TestBackupService_ExportToMarkdown tests Markdown export functionality
func TestBackupService_ExportToMarkdown(t *testing.T) {
	// Create temporary directory for export testing
	tempDir, err := os.MkdirTemp("", "libros_test_md")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	service := services.NewBackupService()

	// Create test data
	testBooks := []models.Book{
		{
			ID:        1,
			Title:     "The Go Programming Language",
			Author:    "Alan Donovan",
			Type:      models.Paperback,
			Notes:     "Excellent reference book",
			CreatedAt: time.Date(2023, 1, 15, 10, 30, 0, 0, time.UTC),
			UpdatedAt: time.Date(2023, 1, 20, 14, 45, 0, 0, time.UTC),
		},
		{
			ID:        2,
			Title:     "Clean Code",
			Author:    "Robert C. Martin",
			Type:      models.Hardback,
			Notes:     "", // Empty notes to test this case
			CreatedAt: time.Date(2023, 2, 10, 9, 15, 0, 0, time.UTC),
			UpdatedAt: time.Date(2023, 2, 10, 9, 15, 0, 0, time.UTC),
		},
	}

	// Test successful export
	t.Run("SuccessfulExport", func(t *testing.T) {
		exportPath := filepath.Join(tempDir, "books_test.md")

		err := service.ExportToMarkdown(testBooks, exportPath)
		if err != nil {
			t.Fatalf("ExportToMarkdown failed: %v", err)
		}

		// Verify file was created
		if _, err := os.Stat(exportPath); os.IsNotExist(err) {
			t.Fatal("Export file was not created")
		}

		// Read and verify the exported Markdown
		content, err := os.ReadFile(exportPath)
		if err != nil {
			t.Fatalf("Failed to read exported file: %v", err)
		}

		contentStr := string(content)

		// Verify header
		if !strings.Contains(contentStr, "# Book Collection Export") {
			t.Error("Markdown should contain main header")
		}

		// Verify metadata
		if !strings.Contains(contentStr, "**Export Date:**") {
			t.Error("Markdown should contain export date")
		}
		if !strings.Contains(contentStr, "**Total Books:** 2") {
			t.Error("Markdown should contain correct total books count")
		}

		// Verify book entries
		for _, book := range testBooks {
			if !strings.Contains(contentStr, book.Title) {
				t.Errorf("Markdown should contain book title: %q", book.Title)
			}
			if !strings.Contains(contentStr, book.Author) {
				t.Errorf("Markdown should contain book author: %q", book.Author)
			}
		}

		// Verify separators between books
		separatorCount := strings.Count(contentStr, "---")
		if separatorCount < len(testBooks) {
			t.Errorf("Markdown should contain separators between books, found %d", separatorCount)
		}
	})
}

// TestBackupService_BackupDatabase tests database file backup functionality
func TestBackupService_BackupDatabase(t *testing.T) {
	// Create temporary directories for testing
	tempDir, err := os.MkdirTemp("", "libros_test_backup")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	service := services.NewBackupService()

	// Create a test database file with some content
	sourceContent := []byte("This is test database content with some data")
	sourcePath := filepath.Join(tempDir, "source.db")
	err = os.WriteFile(sourcePath, sourceContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Test successful backup
	t.Run("SuccessfulBackup", func(t *testing.T) {
		destPath := filepath.Join(tempDir, "backup.db")

		err := service.BackupDatabase(sourcePath, destPath)
		if err != nil {
			t.Fatalf("BackupDatabase failed: %v", err)
		}

		// Verify backup file was created
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			t.Fatal("Backup file was not created")
		}

		// Verify content integrity
		backupContent, err := os.ReadFile(destPath)
		if err != nil {
			t.Fatalf("Failed to read backup file: %v", err)
		}

		if string(backupContent) != string(sourceContent) {
			t.Error("Backup content does not match source content")
		}
	})

	// Test backup of nonexistent source
	t.Run("BackupNonexistentSource", func(t *testing.T) {
		nonexistentSource := filepath.Join(tempDir, "nonexistent.db")
		destPath := filepath.Join(tempDir, "backup_fail.db")

		err := service.BackupDatabase(nonexistentSource, destPath)
		if err == nil {
			t.Error("BackupDatabase should fail with nonexistent source")
		}
	})
}