package unit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/papadavis47/libros/internal/constants"
)

// TestConstants_Values tests that all constants are properly defined
// This ensures the application has consistent sizing and limits across all components
func TestConstants_Values(t *testing.T) {
	tests := []struct {
		name     string
		actual   int
		expected int
	}{
		{"InputFieldWidth", constants.InputFieldWidth, 50},
		{"TextAreaWidth", constants.TextAreaWidth, 60},
		{"TitleMaxLength", constants.TitleMaxLength, 255},
		{"AuthorMaxLength", constants.AuthorMaxLength, 255},
		{"NotesMaxLength", constants.NotesMaxLength, 1000},
		{"BooksPerPage", constants.BooksPerPage, 3},
		{"TextWrapWidth", constants.TextWrapWidth, 60},
		{"NoteTruncateLength", constants.NoteTruncateLength, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.actual != tt.expected {
				t.Errorf("Constant %s = %d, want %d", tt.name, tt.actual, tt.expected)
			}
		})
	}
}

// TestConstants_FilePermissions tests that file permission constants are correctly set
// This ensures proper security settings for created files and directories
func TestConstants_FilePermissions(t *testing.T) {
	tests := []struct {
		name     string
		actual   os.FileMode
		expected os.FileMode
	}{
		{"DirPermissions", constants.DirPermissions, 0755},
		{"FilePermissions", constants.FilePermissions, 0644},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.actual != tt.expected {
				t.Errorf("File permission constant %s = %o, want %o", tt.name, tt.actual, tt.expected)
			}
		})
	}
}

// TestConstants_DefaultPaths tests that default path constants are properly set
// This ensures the application knows where to store its data and configuration
func TestConstants_DefaultPaths(t *testing.T) {
	tests := []struct {
		name     string
		actual   string
		expected string
	}{
		{"DefaultAppDir", constants.DefaultAppDir, "~/.libros"},
		{"DatabaseFilename", constants.DatabaseFilename, "libros.db"},
		{"BackupDir", constants.BackupDir, "backups"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.actual != tt.expected {
				t.Errorf("Path constant %s = %q, want %q", tt.name, tt.actual, tt.expected)
			}
		})
	}
}

// TestGetAppDir tests the application directory resolution function
// This function handles home directory expansion and fallback behavior
func TestGetAppDir(t *testing.T) {
	// Test that GetAppDir returns a valid directory path
	appDir := constants.GetAppDir()
	
	if appDir == "" {
		t.Error("GetAppDir() returned empty string")
	}
	
	// Test that the path is absolute (either starts with / or contains expanded home)
	if !filepath.IsAbs(appDir) && appDir != ".libros" {
		t.Errorf("GetAppDir() returned non-absolute path: %q", appDir)
	}
	
	// Test that the path ends with .libros
	if filepath.Base(appDir) != ".libros" {
		t.Errorf("GetAppDir() should end with .libros, got: %q", appDir)
	}
}

// TestGetDatabasePath tests the database path resolution function
// This function combines the app directory with the database filename
func TestGetDatabasePath(t *testing.T) {
	dbPath := constants.GetDatabasePath()
	
	if dbPath == "" {
		t.Error("GetDatabasePath() returned empty string")
	}
	
	// Test that the path ends with the database filename
	if filepath.Base(dbPath) != constants.DatabaseFilename {
		t.Errorf("GetDatabasePath() should end with %q, got: %q", constants.DatabaseFilename, filepath.Base(dbPath))
	}
	
	// Test that the directory part matches GetAppDir()
	expectedDir := constants.GetAppDir()
	actualDir := filepath.Dir(dbPath)
	if actualDir != expectedDir {
		t.Errorf("GetDatabasePath() directory = %q, want %q", actualDir, expectedDir)
	}
}

// TestGetBackupDir tests the backup directory resolution function
// This function combines the app directory with the backup subdirectory
func TestGetBackupDir(t *testing.T) {
	backupDir := constants.GetBackupDir()
	
	if backupDir == "" {
		t.Error("GetBackupDir() returned empty string")
	}
	
	// Test that the path ends with the backup directory name
	if filepath.Base(backupDir) != constants.BackupDir {
		t.Errorf("GetBackupDir() should end with %q, got: %q", constants.BackupDir, filepath.Base(backupDir))
	}
	
	// Test that the parent directory matches GetAppDir()
	expectedParent := constants.GetAppDir()
	actualParent := filepath.Dir(backupDir)
	if actualParent != expectedParent {
		t.Errorf("GetBackupDir() parent directory = %q, want %q", actualParent, expectedParent)
	}
}

// TestConstants_Logical_Relationships tests logical relationships between constants
// This ensures constants make sense relative to each other
func TestConstants_Logical_Relationships(t *testing.T) {
	// TextAreaWidth should be larger than InputFieldWidth for better notes editing
	if constants.TextAreaWidth <= constants.InputFieldWidth {
		t.Errorf("TextAreaWidth (%d) should be larger than InputFieldWidth (%d)", 
			constants.TextAreaWidth, constants.InputFieldWidth)
	}
	
	// NotesMaxLength should be larger than TitleMaxLength and AuthorMaxLength
	if constants.NotesMaxLength <= constants.TitleMaxLength {
		t.Errorf("NotesMaxLength (%d) should be larger than TitleMaxLength (%d)", 
			constants.NotesMaxLength, constants.TitleMaxLength)
	}
	
	if constants.NotesMaxLength <= constants.AuthorMaxLength {
		t.Errorf("NotesMaxLength (%d) should be larger than AuthorMaxLength (%d)", 
			constants.NotesMaxLength, constants.AuthorMaxLength)
	}
	
	// NoteTruncateLength should be less than NotesMaxLength
	if constants.NoteTruncateLength >= constants.NotesMaxLength {
		t.Errorf("NoteTruncateLength (%d) should be less than NotesMaxLength (%d)", 
			constants.NoteTruncateLength, constants.NotesMaxLength)
	}
	
	// BooksPerPage should be a reasonable number (not too high or too low)
	if constants.BooksPerPage < 1 || constants.BooksPerPage > 10 {
		t.Errorf("BooksPerPage (%d) should be between 1 and 10 for good UX", constants.BooksPerPage)
	}
}

// TestConstants_Immutability tests that constants maintain their values
// This is a regression test to ensure constants aren't accidentally modified
func TestConstants_Immutability(t *testing.T) {
	// Store initial values
	initialInputWidth := constants.InputFieldWidth
	initialTextAreaWidth := constants.TextAreaWidth
	initialTitleMaxLength := constants.TitleMaxLength
	
	// Call functions that might modify internal state
	_ = constants.GetAppDir()
	_ = constants.GetDatabasePath()
	_ = constants.GetBackupDir()
	
	// Verify constants haven't changed
	if constants.InputFieldWidth != initialInputWidth {
		t.Errorf("InputFieldWidth changed from %d to %d", initialInputWidth, constants.InputFieldWidth)
	}
	if constants.TextAreaWidth != initialTextAreaWidth {
		t.Errorf("TextAreaWidth changed from %d to %d", initialTextAreaWidth, constants.TextAreaWidth)
	}
	if constants.TitleMaxLength != initialTitleMaxLength {
		t.Errorf("TitleMaxLength changed from %d to %d", initialTitleMaxLength, constants.TitleMaxLength)
	}
}