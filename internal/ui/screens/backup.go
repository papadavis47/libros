package screens

import (
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/papadavis47/libros/internal/constants"
	"github.com/papadavis47/libros/internal/database"
	"github.com/papadavis47/libros/internal/models"
	"github.com/papadavis47/libros/internal/styles"
)

type BackupScreen struct {
	db      *database.DB
	status  string
	isError bool
	done    bool
}

func NewBackupScreen(db *database.DB) *BackupScreen {
	screen := &BackupScreen{
		db: db,
	}
	// Perform backup immediately when screen is created
	screen.performBackupSync()
	return screen
}

func (s *BackupScreen) ClearStatus() {
	s.status = ""
	s.isError = false
	s.done = false
}

func (s *BackupScreen) Init() tea.Cmd {
	return nil
}

func (s *BackupScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "enter":
			// Return to utilities screen
			return s, SwitchScreenCmd(models.UtilitiesScreen)
		case "q", "ctrl+c":
			return s, tea.Quit
		}
	}

	return s, nil
}

func (s *BackupScreen) View() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(styles.TitleStyle.Render("Ｄａｔａｂａｓｅ　Ｂａｃｋｕｐ"))
	b.WriteString("\n\n")

	// Show backup result
	if s.status != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true).
			Padding(1, 0).
			PaddingLeft(3)

		if s.isError {
			statusStyle = statusStyle.Foreground(lipgloss.Color("#FF0000"))
		}

		b.WriteString("\n" + statusStyle.Render(styles.AddLetterSpacing(s.status)))
		b.WriteString("\n")
	}

	if s.done {
		b.WriteString("\n" + styles.HelpTextStyle.Render(styles.AddLetterSpacing("Press Enter or Esc to return to Utilities")))
	}

	return b.String()
}

func (s *BackupScreen) performBackupSync() {
	s.done = true

	// Get the database file path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		s.status = "Database backup failed: " + err.Error()
		s.isError = true
		return
	}

	dbPath := filepath.Join(homeDir, ".libros", "books.db")
	backupPath := filepath.Join(homeDir, ".libros", "books.db.bak")

	// Check if source database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		s.status = "Database backup failed: " + err.Error()
		s.isError = true
		return
	}

	// Copy the database file to backup location
	err = copyFile(dbPath, backupPath)
	if err != nil {
		s.status = "Database backup failed: " + err.Error()
		s.isError = true
		return
	}

	s.status = "Database backed up successfully to ~/.libros/books.db.bak"
	s.isError = false
}

// copyFile copies a file from src to dst, overwriting dst if it exists
func copyFile(src, dst string) error {
	sourceFile, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, sourceFile, constants.FilePermissions)
}
