package constants

import (
	"os"
	"path/filepath"
	"testing"

)

// TestConstants_Values tests that all constants are properly defined
// This ensures the application has consistent sizing and limits across all components
func TestConstants_Values(t *testing.T) {
	tests := []struct {
		name     string
		actual   int
		expected int
	}{
		{"InputFieldWidth", InputFieldWidth, 50},
		{"TextAreaWidth", TextAreaWidth, 60},
		{"TitleMaxLength", TitleMaxLength, 255},
		{"AuthorMaxLength", AuthorMaxLength, 255},
		{"NotesMaxLength", NotesMaxLength, 1000},
		{"BooksPerPage", BooksPerPage, 3},
		{"TextWrapWidth", TextWrapWidth, 60},
		{"NoteTruncateLength", NoteTruncateLength, 100},
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
		{"DirPermissions", DirPermissions, 0755},
		{"FilePermissions", FilePermissions, 0644},
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
		{"DefaultAppDir", DefaultAppDir, "~/.libros"},
		{"DatabaseFilename", DatabaseFilename, "libros.db"},
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
	appDir := GetAppDir()
	
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
	dbPath := GetDatabasePath()
	
	if dbPath == "" {
		t.Error("GetDatabasePath() returned empty string")
	}
	
	// Test that the path ends with the database filename
	if filepath.Base(dbPath) != DatabaseFilename {
		t.Errorf("GetDatabasePath() should end with %q, got: %q", DatabaseFilename, filepath.Base(dbPath))
	}
	
	// Test that the directory part matches GetAppDir()
	expectedDir := GetAppDir()
	actualDir := filepath.Dir(dbPath)
	if actualDir != expectedDir {
		t.Errorf("GetDatabasePath() directory = %q, want %q", actualDir, expectedDir)
	}
}


// TestConstants_Logical_Relationships tests logical relationships between constants
// This ensures constants make sense relative to each other
func TestConstants_Logical_Relationships(t *testing.T) {
	// TextAreaWidth should be larger than InputFieldWidth for better notes editing
	if TextAreaWidth <= InputFieldWidth {
		t.Errorf("TextAreaWidth (%d) should be larger than InputFieldWidth (%d)", 
			TextAreaWidth, InputFieldWidth)
	}
	
	// NotesMaxLength should be larger than TitleMaxLength and AuthorMaxLength
	if NotesMaxLength <= TitleMaxLength {
		t.Errorf("NotesMaxLength (%d) should be larger than TitleMaxLength (%d)", 
			NotesMaxLength, TitleMaxLength)
	}
	
	if NotesMaxLength <= AuthorMaxLength {
		t.Errorf("NotesMaxLength (%d) should be larger than AuthorMaxLength (%d)", 
			NotesMaxLength, AuthorMaxLength)
	}
	
	// NoteTruncateLength should be less than NotesMaxLength
	if NoteTruncateLength >= NotesMaxLength {
		t.Errorf("NoteTruncateLength (%d) should be less than NotesMaxLength (%d)", 
			NoteTruncateLength, NotesMaxLength)
	}
	
	// BooksPerPage should be a reasonable number (not too high or too low)
	if BooksPerPage < 1 || BooksPerPage > 10 {
		t.Errorf("BooksPerPage (%d) should be between 1 and 10 for good UX", BooksPerPage)
	}
}

// TestConstants_Immutability tests that constants maintain their values
// This is a regression test to ensure constants aren't accidentally modified
func TestConstants_Immutability(t *testing.T) {
	// Store initial values
	initialInputWidth := InputFieldWidth
	initialTextAreaWidth := TextAreaWidth
	initialTitleMaxLength := TitleMaxLength
	
	// Call functions that might modify internal state
	_ = GetAppDir()
	_ = GetDatabasePath()
	
	// Verify constants haven't changed
	if InputFieldWidth != initialInputWidth {
		t.Errorf("InputFieldWidth changed from %d to %d", initialInputWidth, InputFieldWidth)
	}
	if TextAreaWidth != initialTextAreaWidth {
		t.Errorf("TextAreaWidth changed from %d to %d", initialTextAreaWidth, TextAreaWidth)
	}
	if TitleMaxLength != initialTitleMaxLength {
		t.Errorf("TitleMaxLength changed from %d to %d", initialTitleMaxLength, TitleMaxLength)
	}
}