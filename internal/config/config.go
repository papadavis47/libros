package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config represents the application configuration
type Config struct {
	Theme Theme `toml:"theme"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Theme: DefaultTheme,
	}
}

// LoadConfig loads the configuration from the user's ~/.libros directory
func LoadConfig() (Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return DefaultConfig(), err
	}

	// If config file doesn't exist, create it with default values
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := DefaultConfig()
		if saveErr := SaveConfig(config); saveErr != nil {
			return config, saveErr
		}
		return config, nil
	}

	// Load existing config
	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return DefaultConfig(), err
	}

	return config, nil
}

// SaveConfig saves the configuration to the user's ~/.libros directory
func SaveConfig(config Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Ensure the directory exists
	if err := ensureConfigDir(); err != nil {
		return err
	}

	// Create the file
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode the config to TOML
	encoder := toml.NewEncoder(file)
	return encoder.Encode(config)
}

// UpdateTheme updates the theme in the configuration and saves it
func UpdateTheme(theme Theme) error {
	config, err := LoadConfig()
	if err != nil {
		// If we can't load config, create a new one
		config = DefaultConfig()
	}

	config.Theme = theme
	return SaveConfig(config)
}

// GetCurrentTheme returns the current theme from the configuration
func GetCurrentTheme() Theme {
	config, err := LoadConfig()
	if err != nil {
		return DefaultTheme
	}
	return config.Theme
}

// getConfigPath returns the path to the configuration file
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".libros", "theme.toml"), nil
}

// ensureConfigDir ensures the .libros directory exists
func ensureConfigDir() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	librosDir := filepath.Join(homeDir, ".libros")
	return os.MkdirAll(librosDir, 0755)
}