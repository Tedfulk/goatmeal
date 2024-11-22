package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
)

type UsernameEditorModel struct {
	textInput textinput.Model
	width     int
	height    int
	err       error
	colors    config.ThemeColors
}

type UsernameUpdatedMsg struct {
	newUsername string
}

func NewUsernameEditor(currentUsername string, colors config.ThemeColors) UsernameEditorModel {
	ti := textinput.New()
	ti.Placeholder = "Enter Username"
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 70
	if currentUsername != "" {
		ti.SetValue(currentUsername)
	}

	return UsernameEditorModel{
		textInput: ti,
		colors:    colors,
	}
}

func (m UsernameEditorModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m UsernameEditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return ChangeViewMsg(settingsView) }
		case "enter":
			if m.textInput.Value() != "" {
				return m, func() tea.Msg {
					return UsernameUpdatedMsg{newUsername: m.textInput.Value()}
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m UsernameEditorModel) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.colors.UserText)).
		Padding(0, 1).
		MarginBottom(0).
		Width(80).
		Align(lipgloss.Center)

	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.colors.UserBubble)).
		Padding(1).
		Width(80)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.colors.MenuNormal)).
		Align(lipgloss.Center)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render("Username"),
		baseStyle.Render(m.textInput.View()),
		helpStyle.Render("Press Enter to submit • Esc to cancel"),
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
} 