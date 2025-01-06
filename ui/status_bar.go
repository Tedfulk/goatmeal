package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
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
	config            *config.Config
	width             int
	conversationTitle string
	height            int
	isSearchMode      bool
	spinner           spinner.Model
	isLoading         bool
	temporaryMessage  string
	temporaryTimer    *time.Timer
	errorMessage      string
	errorTimer        *time.Timer
	isEnhancedSearch bool
}

// NewStatusBar creates a new status bar
func NewStatusBar(cfg *config.Config, convoTitle string) *StatusBar {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(theme.CurrentTheme.Primary.GetColor())

	return &StatusBar{
		config:            cfg,
		conversationTitle: convoTitle,
		width:            0,
		height:           1,
		spinner:          s,
		isLoading:        false,
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

// Add method to set search mode
func (s *StatusBar) SetSearchMode(enabled bool) {
	s.isSearchMode = enabled
}

// SetLoading sets the loading state of the status bar
func (s *StatusBar) SetLoading(loading bool) {
	s.isLoading = loading
}

// Update handles the spinner animation
func (s *StatusBar) Update(msg tea.Msg) tea.Cmd {
	if s.isLoading {
		var cmd tea.Cmd
		s.spinner, cmd = s.spinner.Update(msg)
		return cmd
	}
	return nil
}

// View renders the status bar
func (s *StatusBar) View() string {
	// Create the status bar style
	statusBarStyle := theme.BaseStyle.StatusBar.
		Width(s.width).
		Foreground(theme.CurrentTheme.StatusBar.Text.GetColor())

	// Create the title style
	titleStyle := lipgloss.NewStyle().
		Foreground(theme.CurrentTheme.StatusBar.Title.GetColor()).
		Align(lipgloss.Center).
		Width((s.width / 2))

	// Create the model style
	modelStyle := lipgloss.NewStyle().
		Foreground(theme.CurrentTheme.StatusBar.Model.GetColor())

	// Build right section content
	rightContent := ""
	if s.errorMessage != "" {
		// Show error in red
		rightContent = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Render(s.errorMessage)
	} else if s.temporaryMessage != "" {
		// Show temporary message in green
		rightContent = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Render(s.temporaryMessage)
	} else if s.isLoading {
		if s.isSearchMode {
			rightContent = s.spinner.View() + " Searching..."
		} else {
			rightContent = s.spinner.View() + " Thinking..."
		}
	} else if s.isSearchMode {
		searchIndicator := "🔍"
		if s.isEnhancedSearch {
			searchIndicator = "🔍+"
		}
		rightContent = searchIndicator
	}

	leftSection := lipgloss.JoinHorizontal(
		lipgloss.Left,
		"➕💬",
		" | ",
		modelStyle.Render(s.config.CurrentProvider+"/"+s.config.CurrentModel),
	)

	// Calculate the right section width to ensure proper alignment
	rightStyle := lipgloss.NewStyle().
		Align(lipgloss.Right).
		Width(s.width/4)

	// Render the status bar with three sections
	return statusBarStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			leftSection,
			titleStyle.Render(s.conversationTitle),
			rightStyle.Render(rightContent),
		),
	)
}

func (s *StatusBar) inBounds(x, y int) bool {
	return x >= 0 && x < 4 && y >= 0 && y < s.height
}

func (s *StatusBar) SetTemporaryText(text string) {
	s.temporaryMessage = text
	if s.temporaryTimer != nil {
		s.temporaryTimer.Stop()
	}
	s.temporaryTimer = time.NewTimer(1 * time.Second)
	go func() {
		<-s.temporaryTimer.C
		s.temporaryMessage = ""
	}()
}

func (s *StatusBar) SetError(text string) {
	s.errorMessage = text
	if s.errorTimer != nil {
		s.errorTimer.Stop()
	}
	s.errorTimer = time.NewTimer(3 * time.Second)
	go func() {
		<-s.errorTimer.C
		s.errorMessage = ""
	}()
}

func (s *StatusBar) SetEnhancedSearch(enabled bool) {
	s.isEnhancedSearch = enabled
} 