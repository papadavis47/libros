package screens

import (
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/papadavis47/libros/internal/database"
	"github.com/papadavis47/libros/internal/messages"
	"github.com/papadavis47/libros/internal/models"
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

type BackupScreen struct {
	db      *database.DB
	items   []string
	index   int
	status  string
	isError bool
}

func NewBackupScreen(db *database.DB) *BackupScreen {
	items := []string{
		"ＪＳＯＮ　Ｆｏｒｍａｔ",
		"Ｍａｒｋｄｏｗｎ　Ｆｏｒｍａｔ",
		"Ｂａｃｋ　ｔｏ　Ｍａｉｎ　Ｍｅｎｕ",
	}

	return &BackupScreen{
		db:    db,
		items: items,
		index: 0,
	}
}

func (s *BackupScreen) ClearStatus() {
	s.status = ""
	s.isError = false
}

func (s *BackupScreen) Init() tea.Cmd {
	return nil
}

func (s *BackupScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if s.index > 0 {
				s.index--
			}
		case "down", "j":
			if s.index < len(s.items)-1 {
				s.index++
			}
		case "enter":
			selectedItem := s.items[s.index]
			switch selectedItem {
			case "ＪＳＯＮ　Ｆｏｒｍａｔ":
				s.status = ""
				return s, s.performBackup("json")
			case "Ｍａｒｋｄｏｗｎ　Ｆｏｒｍａｔ":
				s.status = ""
				return s, s.performBackup("markdown")
			case "Ｂａｃｋ　ｔｏ　Ｍａｉｎ　Ｍｅｎｕ":
				return s, SwitchScreenCmd(models.MenuScreen)
			}
		case "q", "ctrl+c":
			return s, tea.Quit
		case "esc":
			return s, SwitchScreenCmd(models.MenuScreen)
		}

	case messages.BackupMsg:
		if msg.Err != nil {
			s.status = "Backup failed: " + msg.Err.Error()
			s.isError = true
		} else {
			s.status = "Backup completed successfully!"
			s.isError = false
		}
	}

	return s, nil
}

func (s *BackupScreen) View() string {
	var b strings.Builder

	// Display application title
	b.WriteString("\n")
	b.WriteString(styles.TitleStyle.Render("Ｌｉｂｒｏｓ　－　Ａ　Ｂｏｏｋ　Ｍａｎａｇｅｒ"))
	b.WriteString("\n\n")
	b.WriteString(styles.BlurredStyle.Render("Ｂａｃｋｕｐ"))
	b.WriteString("\n\n")

	// Render each menu item with appropriate styling
	for i, item := range s.items {
		if i == s.index {
			// Highlight currently selected item
			b.WriteString(styles.SelectedStyle.Render(item))
		} else {
			// Dim non-selected items
			b.WriteString(styles.BlurredStyle.Render(item))
		}
		b.WriteString("\n\n")
	}

	// Display status message if present
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

	// Display help text
	b.WriteString("\n\n" + styles.HelpTextStyle.Render("Use ↑/↓ or j/k to navigate, Enter to select, Esc to go back"))

	return b.String()
}

func (s *BackupScreen) performBackup(format string) tea.Cmd {
	return func() tea.Msg {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return messages.BackupMsg{Err: err}
		}

		backupDir := filepath.Join(homeDir, ".libros")

		switch format {
		case "json":
			err = s.db.BackupToJSON(backupDir)
		case "markdown":
			err = s.db.BackupToMarkdown(backupDir)
		}

		return messages.BackupMsg{Err: err}
	}
}
