package screens

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/papadavis47/libros/internal/constants"
	"github.com/papadavis47/libros/internal/database"
	"github.com/papadavis47/libros/internal/factory"
	"github.com/papadavis47/libros/internal/messages"
	"github.com/papadavis47/libros/internal/models"
	"github.com/papadavis47/libros/internal/services"
	"github.com/papadavis47/libros/internal/styles"
)

type SwitchScreenMsg struct {
	Screen models.Screen
}

func SwitchScreenCmd(screen models.Screen) tea.Cmd {
	return func() tea.Msg {
		return SwitchScreenMsg{Screen: screen}
	}
}

// ExportState represents the current state of the export flow
type ExportState int

const (
	PathInput ExportState = iota // Getting file path input from user
	FormatSelection              // Selecting export format (JSON/Markdown)
	Exporting                    // Currently performing export
	ShowResult                   // Showing export result (success/error)
)

type ExportScreen struct {
	db                *database.DB
	state             ExportState
	pathInput         textinput.Model
	exportPath        string
	formatItems       []string
	formatIndex       int
	status            string
	isError           bool
	defaultExportsDir string
	lastExportedFile  string
}

func NewExportScreen(db *database.DB) *ExportScreen {
	// Get default exports directory
	homeDir, _ := os.UserHomeDir()
	defaultExportsDir := filepath.Join(homeDir, ".libros", "exports")

	// Initialize text input for file path using factory function
	pathInput := factory.CreatePathInput(defaultExportsDir)
	pathInput.Focus()

	formatItems := []string{
		"ＪＳＯＮ　Ｆｏｒｍａｔ",
		"Ｍａｒｋｄｏｗｎ　Ｆｏｒｍａｔ",
		"Ｂａｃｋ　ｔｏ　Ｕｔｉｌｉｔｉｅｓ",
		"Ｂａｃｋ　ｔｏ　Ｍａｉｎ　Ｍｅｎｕ",
	}

	return &ExportScreen{
		db:                db,
		state:             PathInput,
		pathInput:         pathInput,
		formatItems:       formatItems,
		formatIndex:       0,
		defaultExportsDir: defaultExportsDir,
	}
}

func (s *ExportScreen) ClearStatus() {
	s.status = ""
	s.isError = false
	s.state = PathInput
	s.pathInput.SetValue("")
	s.pathInput.Prompt = "   " // Ensure proper alignment
	s.pathInput.Focus()
	s.formatIndex = 0
	s.lastExportedFile = ""
}

func (s *ExportScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (s *ExportScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch s.state {
	case PathInput:
		return s.updatePathInput(msg)
	case FormatSelection:
		return s.updateFormatSelection(msg)
	case Exporting:
		return s.updateExporting(msg)
	case ShowResult:
		return s.updateShowResult(msg)
	}
	return s, nil
}

func (s *ExportScreen) updatePathInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Validate and process the path input
			inputPath := strings.TrimSpace(s.pathInput.Value())
			if inputPath == "" {
				// Use default path
				s.exportPath = s.defaultExportsDir
			} else {
				// Validate absolute path
				if !strings.HasPrefix(inputPath, "/") && !strings.HasPrefix(inputPath, "~") {
					s.status = "Please enter an absolute path (starting with / or ~)"
					s.isError = true
					return s, nil
				}
				
				// Expand ~ to home directory
				if strings.HasPrefix(inputPath, "~") {
					homeDir, err := os.UserHomeDir()
					if err != nil {
						s.status = "Error getting home directory: " + err.Error()
						s.isError = true
						return s, nil
					}
					inputPath = strings.Replace(inputPath, "~", homeDir, 1)
				}
				
				// Check if directory exists and is writable
				if err := s.validatePath(inputPath); err != nil {
					s.status = err.Error()
					s.isError = true
					return s, nil
				}
				
				s.exportPath = inputPath
			}
			
			// Move to format selection
			s.state = FormatSelection
			s.status = ""
			s.isError = false
			return s, nil
			
		case "esc":
			return s, SwitchScreenCmd(models.UtilitiesScreen)
		case "q", "ctrl+c":
			return s, tea.Quit
		}
	}

	// Update text input
	s.pathInput, cmd = s.pathInput.Update(msg)
	return s, cmd
}

func (s *ExportScreen) updateFormatSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if s.formatIndex > 0 {
				s.formatIndex--
			}
		case "down", "j":
			if s.formatIndex < len(s.formatItems)-1 {
				s.formatIndex++
			}
		case "enter":
			selectedItem := s.formatItems[s.formatIndex]
			switch selectedItem {
			case "ＪＳＯＮ　Ｆｏｒｍａｔ":
				s.state = Exporting
				s.status = "Exporting to JSON..."
				s.isError = false
				s.lastExportedFile = filepath.Join(s.exportPath, "books.json")
				return s, s.performExport("json")
			case "Ｍａｒｋｄｏｗｎ　Ｆｏｒｍａｔ":
				s.state = Exporting
				s.status = "Exporting to Markdown..."
				s.isError = false
				s.lastExportedFile = filepath.Join(s.exportPath, "books.md")
				return s, s.performExport("markdown")
			case "Ｂａｃｋ　ｔｏ　Ｕｔｉｌｉｔｉｅｓ":
				return s, SwitchScreenCmd(models.UtilitiesScreen)
			case "Ｂａｃｋ　ｔｏ　Ｍａｉｎ　Ｍｅｎｕ":
				return s, SwitchScreenCmd(models.MenuScreen)
			}
		case "esc":
			// Go back to path input
			s.state = PathInput
			s.pathInput.Prompt = "   " // Ensure proper alignment
			s.pathInput.Focus()
			return s, textinput.Blink
		case "q", "ctrl+c":
			return s, tea.Quit
		}
	}
	return s, nil
}

func (s *ExportScreen) updateExporting(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return s, SwitchScreenCmd(models.UtilitiesScreen)
		case "q", "ctrl+c":
			return s, tea.Quit
		}
	case messages.BackupMsg:
		if msg.Err != nil {
			s.status = "Export failed: " + msg.Err.Error()
			s.isError = true
		} else {
			s.status = "Export completed successfully!\n\nFile saved to: " + s.lastExportedFile
			s.isError = false
		}
		s.state = ShowResult
	}
	return s, nil
}

func (s *ExportScreen) updateShowResult(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "esc":
			// Return to format selection to allow another export
			s.state = FormatSelection
			s.status = ""
			s.isError = false
			return s, nil
		case "q", "ctrl+c":
			return s, tea.Quit
		}
	}
	return s, nil
}

func (s *ExportScreen) validatePath(path string) error {
	// Check if directory exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	
	// Check if we can write to the directory
	testFile := filepath.Join(path, ".libros_test")
	file, err := os.Create(testFile)
	if err != nil {
		return err
	}
	file.Close()
	os.Remove(testFile)
	
	return nil
}

func (s *ExportScreen) View() string {
	var b strings.Builder

	// Display title
	b.WriteString("\n")
	b.WriteString(styles.TitleStyle.Render("Ｅｘｐｏｒｔ　Ｂｏｏｋ　Ｃｏｌｌｅｃｔｉｏｎ"))
	b.WriteString("\n\n")

	switch s.state {
	case PathInput:
		b.WriteString(styles.BlurredStyle.Render(styles.AddLetterSpacing("Enter export directory path (or press Enter for default):")))
		b.WriteString("\n\n")
		b.WriteString(styles.BlurredStyle.Render(styles.AddLetterSpacing("Default: " + s.defaultExportsDir)))
		b.WriteString("\n\n")
		b.WriteString(s.pathInput.View())
		b.WriteString("\n\n")
		
		if s.status != "" && s.isError {
			errorStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF0000")).
				Bold(true).
				Padding(1, 0).
				PaddingLeft(3)
			b.WriteString("\n" + errorStyle.Render(styles.AddLetterSpacing(s.status)))
			b.WriteString("\n")
		}
		
		b.WriteString("\n" + styles.HelpTextStyle.Render(styles.AddLetterSpacing("Enter to continue, Esc to go back")))

	case FormatSelection:
		b.WriteString(styles.BlurredStyle.Render(styles.AddLetterSpacing("Export to: " + s.exportPath)))
		b.WriteString("\n\n")
		b.WriteString(styles.BlurredStyle.Render(styles.AddLetterSpacing("Select export format:")))
		b.WriteString("\n\n")

		// Render format options
		for i, item := range s.formatItems {
			if i == s.formatIndex {
				b.WriteString(styles.SelectedStyle.Render(item))
			} else {
				b.WriteString(styles.BlurredStyle.Render(item))
			}
			b.WriteString("\n\n")
		}

		b.WriteString("\n" + styles.HelpTextStyle.Render(styles.AddLetterSpacing("Use ↑/↓ or j/k to navigate, Enter to select, Esc to go back")))

	case Exporting:
		b.WriteString(styles.BlurredStyle.Render(styles.AddLetterSpacing("Export to: " + s.exportPath)))
		b.WriteString("\n\n")
		
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
		}
		b.WriteString("\n\n")
		b.WriteString("\n" + styles.HelpTextStyle.Render(styles.AddLetterSpacing("Esc to go back")))

	case ShowResult:
		b.WriteString(styles.BlurredStyle.Render(styles.AddLetterSpacing("Export to: " + s.exportPath)))
		b.WriteString("\n\n")
		
		if s.status != "" {
			statusStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00FF00")).
				Bold(true).
				Padding(1, 0).
				PaddingLeft(3)

			if s.isError {
				statusStyle = statusStyle.Foreground(lipgloss.Color("#FF0000"))
			}

			// Handle multi-line status messages properly
			lines := strings.Split(s.status, "\n")
			for i, line := range lines {
				if line != "" {
					if i == 0 {
						b.WriteString("\n" + statusStyle.Render(styles.AddLetterSpacing(line)))
					} else {
						b.WriteString("\n" + statusStyle.Render(styles.AddLetterSpacing(line)))
					}
				} else {
					b.WriteString("\n")
				}
			}
		}
		b.WriteString("\n\n")
		b.WriteString("\n" + styles.HelpTextStyle.Render(styles.AddLetterSpacing("Press Enter or Esc to continue")))
	}

	return b.String()
}

func (s *ExportScreen) performExport(format string) tea.Cmd {
	return func() tea.Msg {
		// Ensure export directory exists
		if err := os.MkdirAll(s.exportPath, constants.DirPermissions); err != nil {
			return messages.BackupMsg{Err: err}
		}

		// Load books from database
		books, err := s.db.LoadBooks()
		if err != nil {
			return messages.BackupMsg{Err: fmt.Errorf("failed to load books: %v", err)}
		}

		// Create backup service and export
		backupService := services.NewBackupService()
		switch format {
		case "json":
			err = backupService.ExportToJSON(books, filepath.Join(s.exportPath, "books.json"))
		case "markdown":
			err = backupService.ExportToMarkdown(books, filepath.Join(s.exportPath, "books.md"))
		}

		return messages.BackupMsg{Err: err}
	}
}