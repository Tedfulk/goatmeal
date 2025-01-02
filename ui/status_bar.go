package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
)

var (
	statusBarStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#1a1a1a"))
		// Padding(0, 1)

	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Background(lipgloss.Color("#1a1a1a"))
		// Padding(0, 1)

	modelStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFFF")).
		Background(lipgloss.Color("#1a1a1a"))
		// Padding(0, 1)
)

// StatusBar represents the status bar at the bottom of the screen
type StatusBar struct {
	config           *config.Config
	width            int
	conversationTitle string
}

// NewStatusBar creates a new status bar
func NewStatusBar(cfg *config.Config, convoTitle string) *StatusBar {
	return &StatusBar{
		config:           cfg,
		conversationTitle: convoTitle,
	}
}

func (s *StatusBar) SetWidth(width int) {
	s.width = width
}

func (s *StatusBar) SetConversationTitle(title string) {
	s.conversationTitle = title
}

// View renders the status bar
func (s *StatusBar) View() string {
	// Title section
	title := titleStyle.Render(s.conversationTitle)

	// Model section
	modelInfo := modelStyle.Render(s.config.CurrentModel)

	// Combine sections with spacing
	bar := lipgloss.JoinHorizontal(
		lipgloss.Center,
		title,
		"  |  ", // Add separator
		modelInfo,
	)

	// Ensure the status bar fills the width
	return statusBarStyle.Width(s.width).Render(bar)
} 