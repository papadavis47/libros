package constants

import "os"

// UI constants for consistent sizing and layout
const (
	// Input field dimensions
	InputFieldWidth     = 50
	TextAreaWidth       = 60
	TitleMaxLength      = 255
	AuthorMaxLength     = 255 
	NotesMaxLength      = 1000
	
	// List and pagination
	BooksPerPage        = 3
	
	// Text wrapping and truncation
	TextWrapWidth       = 60
	NoteTruncateLength  = 100
	
	// File permissions
	DirPermissions      = 0755
	FilePermissions     = 0644
)

// Application paths and directories
var (
	// Default application directory
	DefaultAppDir = "~/.libros"
	
	// Database filename
	DatabaseFilename = "libros.db"
	
	// Backup directory
	BackupDir = "backups"
)

// GetAppDir returns the application directory path, expanding ~ if necessary
func GetAppDir() string {
	if DefaultAppDir == "~/.libros" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ".libros" // fallback to current directory
		}
		return homeDir + "/.libros"
	}
	return DefaultAppDir
}

// GetDatabasePath returns the full path to the database file
func GetDatabasePath() string {
	return GetAppDir() + "/" + DatabaseFilename
}

// GetBackupDir returns the full path to the backup directory
func GetBackupDir() string {
	return GetAppDir() + "/" + BackupDir
}