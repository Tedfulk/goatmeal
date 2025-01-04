package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
	"github.com/tedfulk/goatmeal/ui/theme"
)

type UsernameSettings struct {
	textInput textinput.Model
	config    *config.Config
	width     int
	height    int
}

func NewUsernameSettings(cfg *config.Config) UsernameSettings {
	ti := textinput.New()
	ti.Placeholder = "Enter new username"
	ti.Focus()

	return UsernameSettings{
		textInput: ti,
		config:    cfg,
	}
}

func (u UsernameSettings) Update(msg tea.Msg) (UsernameSettings, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if u.textInput.Value() != "" {
				u.config.Settings.Username = u.textInput.Value()
				manager, err := config.NewManager()
				if err == nil {
					manager.UpdateSettings(u.config.Settings)
				}
			}
		}

		u.textInput, cmd = u.textInput.Update(msg)
		return u, cmd
	}

	return u, nil
}

func (u UsernameSettings) View() string {
	menuStyle := theme.BaseStyle.Menu.
		BorderForeground(theme.CurrentTheme.Primary.GetColor())

	titleStyle := theme.BaseStyle.Title.
		Foreground(theme.CurrentTheme.Primary.GetColor())

	helpStyle := lipgloss.NewStyle().
		Foreground(theme.CurrentTheme.Secondary.GetColor()).
		Align(lipgloss.Center)

	menuContent := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Username Settings"),
		"",
		"Current username: " + u.config.Settings.Username,
		"",
		u.textInput.View(),
		"",
		helpStyle.Render("      enter: save • esc: back • q: quit"),
	)

	return lipgloss.Place(
		u.width,
		u.height,
		lipgloss.Center,
		lipgloss.Center,
		menuStyle.Render(menuContent),
	)
}

func (u *UsernameSettings) SetSize(width, height int) {
	u.width = width
	u.height = height
} 