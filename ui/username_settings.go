package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
)

type UsernameSettings struct {
	textInput textinput.Model
	config    *config.Config
	width     int
	height    int
}

func NewUsernameSettings(cfg *config.Config) UsernameSettings {
	ti := textinput.New()
	ti.Placeholder = "Enter username"
	ti.Width = 40
	ti.Focus()
	ti.SetValue(cfg.Settings.Username)

	return UsernameSettings{
		textInput: ti,
		config:    cfg,
	}
}

func (u UsernameSettings) Update(msg tea.Msg) (UsernameSettings, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return u, tea.Quit
		case "enter":
			// Update the config
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

func (u UsernameSettings) View() string {
	menuStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 2).
		Width(50)

	titleStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Width(46).
		Align(lipgloss.Center)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
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