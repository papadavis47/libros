package screens

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/papadavis47/libros/internal/config"
	"github.com/papadavis47/libros/internal/models"
	"github.com/papadavis47/libros/internal/styles"
)

// ThemeModel represents the theme selection screen
type ThemeModel struct {
	options []ThemeOption
	index   int
}

// ThemeOption represents a theme option for display
type ThemeOption struct {
	Name        string
	Value       string
	Color       string
	DisplayName string
}

// ThemeSelectedMsg is sent when a theme is selected
type ThemeSelectedMsg struct {
	Theme string
}

// NewThemeModel creates a new theme selection model
func NewThemeModel() ThemeModel {
	// Get theme options
	themeOptions := config.ThemeOptions()
	
	// Convert to ThemeOption struct - no need for pre-styled display names
	var options []ThemeOption
	for _, option := range themeOptions {
		options = append(options, ThemeOption{
			Name:        option.Name,
			Value:       option.Value,
			Color:       option.Color,
			DisplayName: styles.AddLetterSpacing(option.Name), // Just the spaced name
		})
	}
	
	// Find current theme index
	currentTheme := config.GetCurrentTheme()
	currentIndex := 0
	for i, option := range options {
		if option.Name == currentTheme.Name {
			currentIndex = i
			break
		}
	}
	
	return ThemeModel{
		options: options,
		index:   currentIndex,
	}
}

// Init initializes the theme selection screen
func (m ThemeModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the theme selection screen
func (m ThemeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Return to main menu without saving
			return m, func() tea.Msg {
				return SwitchScreenMsg{Screen: models.MenuScreen}
			}
		case "up", "k":
			// Move selection up
			if m.index > 0 {
				m.index--
			}
		case "down", "j":
			// Move selection down
			if m.index < len(m.options)-1 {
				m.index++
			}
		case "enter":
			// Select current theme
			selectedOption := m.options[m.index]
			selectedTheme := config.GetThemeByValue(selectedOption.Value)
			
			// Save the selected theme
			if err := config.UpdateTheme(selectedTheme); err != nil {
				// Handle error - for now just return to menu
				return m, func() tea.Msg {
					return SwitchScreenMsg{Screen: models.MenuScreen}
				}
			}
			
			// Send theme selected message and return to menu
			return m, tea.Batch(
				func() tea.Msg {
					return ThemeSelectedMsg{Theme: selectedOption.Value}
				},
				func() tea.Msg {
					return SwitchScreenMsg{Screen: models.MenuScreen}
				},
			)
		}
	}
	
	return m, nil
}

// View renders the theme selection screen
func (m ThemeModel) View() string {
	var b strings.Builder

	// Display application title and screen subtitle
	b.WriteString("\n")
	b.WriteString(styles.TitleStyle().Render("Ｌｉｂｒｏｓ　－　Ａ　Ｂｏｏｋ　Ｍａｎａｇｅｒ"))
	b.WriteString("\n\n")
	b.WriteString(styles.BlurredStyle.Render(styles.AddLetterSpacing("Pick Theme")))
	b.WriteString("\n\n")
	b.WriteString(styles.BlurredStyle.Render(styles.AddLetterSpacing("Choose your preferred color theme for the application")))
	b.WriteString("\n\n")

	// Render each theme option with dynamic background colors
	for i, option := range m.options {
		if i == m.index {
			// Selected item: use the theme's color as background with white text
			selectedStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color(option.Color)).
				Padding(0, 1).
				MarginLeft(2).
				PaddingLeft(1)
			
			b.WriteString(selectedStyle.Render(option.DisplayName))
		} else {
			// Non-selected items: show with theme color as text color
			unselectedStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(option.Color)).
				PaddingLeft(3)
			
			b.WriteString(unselectedStyle.Render(option.DisplayName))
		}
		b.WriteString("\n\n")
	}

	// Display help text for user guidance
	b.WriteString("\n" + styles.HelpTextStyle.Render(styles.AddLetterSpacing("Use ↑/↓ or j/k to navigate, Enter to select, Esc to return to menu")))

	return b.String()
}