package styles

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4"))

	FocusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4"))

	BlurredStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF"))

	NoStyle = lipgloss.NewStyle()

	SelectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	ButtonStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFA500"))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00"))
)
