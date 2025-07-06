package config

import "github.com/charmbracelet/lipgloss"

// Theme represents a visual theme configuration
type Theme struct {
	Name           string `toml:"name"`
	PrimaryColor   string `toml:"primary_color"`
	SecondaryColor string `toml:"secondary_color"`
	TertiaryColor  string `toml:"tertiary_color"`
}

// ThemeOption represents a theme option for selection
type ThemeOption struct {
	Name  string
	Value string
	Color string
}

// Available themes
var (
	// DefaultTheme is the original purple theme
	DefaultTheme = Theme{
		Name:           "Default",
		PrimaryColor:   "#7D56F4",
		SecondaryColor: "#FFA500",
		TertiaryColor:  "#FFD700",
	}

	// PeachRedTheme is a warm red theme
	PeachRedTheme = Theme{
		Name:         "Peach Red",
		PrimaryColor: "#ff5d62",
		// SecondaryColor: "#78e08f",
		SecondaryColor: "#b8e994",
		TertiaryColor:  "#7bed9f",
	}

	// SurimiOrangeTheme is a bright orange theme
	SurimiOrangeTheme = Theme{
		Name:           "Surimi Orange",
		PrimaryColor:   "#ff9e3b",
		SecondaryColor: "#70a1ff",
		TertiaryColor:  "#1e90ff",
	}

	// SpringBlueTheme is a cool blue theme
	SpringBlueTheme = Theme{
		Name:           "Spring Blue",
		PrimaryColor:   "#7fb4ca",
		SecondaryColor: "#f8a5c2",
		TertiaryColor:  "#f78fb3",
	}
)

// AllThemes returns all available themes
func AllThemes() []Theme {
	return []Theme{
		DefaultTheme,
		PeachRedTheme,
		SurimiOrangeTheme,
		SpringBlueTheme,
	}
}

// ThemeOptions returns theme options for selection interface
func ThemeOptions() []ThemeOption {
	return []ThemeOption{
		{Name: "Default", Value: "default", Color: "#7D56F4"},
		{Name: "Peach Red", Value: "peach_red", Color: "#ff5d62"},
		{Name: "Surimi Orange", Value: "surimi_orange", Color: "#ff9e3b"},
		{Name: "Spring Blue", Value: "spring_blue", Color: "#7fb4ca"},
	}
}

// GetThemeByValue returns a theme by its value identifier
func GetThemeByValue(value string) Theme {
	switch value {
	case "default":
		return DefaultTheme
	case "peach_red":
		return PeachRedTheme
	case "surimi_orange":
		return SurimiOrangeTheme
	case "spring_blue":
		return SpringBlueTheme
	default:
		return DefaultTheme
	}
}

// GetThemeByName returns a theme by its display name
func GetThemeByName(name string) Theme {
	for _, theme := range AllThemes() {
		if theme.Name == name {
			return theme
		}
	}
	return DefaultTheme
}

// CreateThemedStyle creates a lipgloss style with the theme's primary color
func CreateThemedStyle(theme Theme) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(theme.PrimaryColor))
}

// CreateThemedBackgroundStyle creates a lipgloss style with the theme's primary color as background
func CreateThemedBackgroundStyle(theme Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(theme.PrimaryColor)).
		Foreground(lipgloss.Color("#FFFFFF"))
}
