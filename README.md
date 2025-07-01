# Libros - Personal Book Manager

A beautiful, interactive terminal-based book management application built with Go and Bubble Tea. Track your personal library with support for multiple book formats, detailed notes, and powerful export/backup features.

## Features

- **📚 Book Management**: Add, edit, view, and manage your personal book collection
- **📝 Detailed Records**: Track title, author, format type, and personal notes for each book
- **🎨 Beautiful UI**: Clean, interactive terminal interface powered by Bubble Tea
- **💾 Multiple Formats**: Support for paperback, hardback, audiobook, and digital formats
- **📊 Export Options**: Export your library to JSON or Markdown formats
- **🔄 Backup & Restore**: Create backups of your entire book database
- **🔍 Smart Navigation**: Intuitive menu system with keyboard shortcuts
- **⚡ Fast Performance**: Lightweight SQLite database for quick access

## Installation

### Prerequisites

- Go 1.24.0 or later
- SQLite3

### Build from Source

```bash
# Clone the repository
git clone https://github.com/papadavis47/libros.git
cd libros

# Build the application
make build

# Or build manually
go build -o libros ./cmd/libros
```

## Usage

### Running the Application

```bash
./libros
```

The application will:
- Create a `.libros` directory in your home folder
- Initialize a SQLite database at `~/.libros/books.db`
- Launch the interactive terminal interface

### Navigation

- Use **↑/↓ arrow keys** to navigate menus
- Press **Enter** to select options
- Press **Esc** to go back to previous screens
- Press **q** to quit the application

### Main Features

#### Adding Books
1. Select "Add a new book" from the main menu
2. Fill in the book details:
   - Title (required)
   - Author (required)
   - Format type (paperback/hardback/audio/digital)
   - Personal notes (optional)
3. Save your book to the collection

#### Managing Your Collection
- **View All Books**: Browse your entire library with formatted display
- **Book Details**: View complete information for any book
- **Edit Books**: Update any book's information
- **Delete Books**: Remove books from your collection

#### Export & Backup
- **JSON Export**: Export your library as structured JSON data
- **Markdown Export**: Create readable Markdown documentation of your books
- **Database Backup**: Create complete backups of your book database

## Project Structure

```
libros/
├── cmd/libros/          # Application entry point
├── internal/
│   ├── constants/       # Application constants
│   ├── database/        # SQLite database layer
│   ├── factory/         # UI component factory
│   ├── interfaces/      # Interface definitions
│   ├── messages/        # Bubble Tea messages
│   ├── models/          # Data models and types
│   ├── services/        # Business logic (backup, export)
│   ├── styles/          # UI styling and themes
│   ├── ui/              # Main UI components and screens
│   ├── utils/           # Utility functions
│   └── validation/      # Input validation
├── tests/
│   ├── unit/           # Unit tests
│   └── integration/    # Integration tests
└── Makefile            # Build and test commands
```

## Development

### Testing

The project includes comprehensive unit and integration tests:

```bash
# Run all tests
make test

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration

# Run tests with coverage
make test-coverage

# Run tests with verbose output
make test-verbose
```

### Building

```bash
# Build the application
make build

# Clean build artifacts
make clean

# Show available commands
make help
```

### Code Architecture

- **Clean Architecture**: Separation of concerns with clear interfaces
- **Bubble Tea Framework**: Modern terminal UI framework
- **SQLite Database**: Lightweight, serverless database
- **Comprehensive Testing**: Unit and integration test coverage
- **Type Safety**: Strongly typed models and interfaces

## Database Schema

The application uses a simple SQLite schema:

- **Books Table**: Stores book information with fields for ID, title, author, type, notes, and timestamps
- **Automatic Migrations**: Database schema is created automatically on first run
- **Data Integrity**: Foreign key constraints and validation ensure data consistency

## File Locations

- **Database**: `~/.libros/books.db`
- **Database Backup**: `~/.libros/books.db.bak` (when backup is created)
- **Exports**: User-specified locations

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Run the test suite (`make test`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## Dependencies

- **Bubble Tea**: Terminal UI framework
- **Bubbles**: UI components for Bubble Tea
- **Lipgloss**: Styling for terminal applications
- **SQLite3**: Database driver

## License

This project is open source. See the repository for license details.

## Support

If you encounter any issues or have questions:
1. Check the existing issues in the repository
2. Create a new issue with detailed information
3. Include steps to reproduce any bugs

---

Built with ❤️ using Go and Bubble Tea