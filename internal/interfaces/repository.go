package interfaces

import "github.com/papadavis47/libros/internal/models"

// BookRepository defines the interface for book data operations
// This abstraction allows for easier testing and potential database backend changes
type BookRepository interface {
	// Book CRUD operations
	SaveBook(book *models.Book) error
	GetAllBooks() ([]models.Book, error)
	GetBookByID(id int) (*models.Book, error)
	UpdateBook(book *models.Book) error
	DeleteBook(id int) error
	
	// Database management
	Close() error
	GetDatabasePath() string
}

// BackupService defines the interface for backup operations
// Separated from repository to follow single responsibility principle
type BackupService interface {
	ExportToJSON(books []models.Book, filePath string) error
	ExportToMarkdown(books []models.Book, filePath string) error
	BackupDatabase(sourcePath, destPath string) error
}