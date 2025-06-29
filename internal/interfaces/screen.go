package interfaces

import tea "github.com/charmbracelet/bubbletea"

// Screen represents a common interface for all screen models in the application
// This allows for consistent handling of screen lifecycle and behavior
type Screen interface {
	// Init initializes the screen and returns any initial commands
	Init() tea.Cmd
	
	// Update handles messages and returns updated model, optional command, and next screen
	Update(tea.Msg) (Screen, tea.Cmd, int)
	
	// View renders the current state of the screen
	View() string
}