package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
	"github.com/tedfulk/goatmeal/ui/theme"
)

var (
	statusBarStyle = theme.BaseStyle.StatusBar.
		Foreground(theme.CurrentTheme.StatusBar.Text.GetColor())

	titleStyle = theme.BaseStyle.StatusBar.
		Foreground(theme.CurrentTheme.StatusBar.Title.GetColor())

	modelStyle = theme.BaseStyle.StatusBar.
		Foreground(theme.CurrentTheme.StatusBar.Model.GetColor())
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

func (s *StatusBar) UpdateStyle() {
	// This method is called when the theme changes
	// The StatusBar already uses theme.CurrentTheme directly in its View method,
	// so it will automatically pick up the new theme colors
}

// View renders the status bar
func (s *StatusBar) View() string {
	// Create the status bar style
	statusBarStyle := theme.BaseStyle.StatusBar.
		Width(s.width).
		Foreground(theme.CurrentTheme.StatusBar.Text.GetColor())

	// Create the title style
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.CurrentTheme.StatusBar.Title.GetColor())

	// Create the model style
	modelStyle := lipgloss.NewStyle().
		Foreground(theme.CurrentTheme.StatusBar.Model.GetColor())

	// Render the status bar
	return statusBarStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			titleStyle.Render(s.conversationTitle),
			" • ",
			modelStyle.Render(s.config.CurrentProvider+"/"+s.config.CurrentModel),
		),
	)
} 