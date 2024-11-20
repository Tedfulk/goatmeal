package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Settings menu items
type SettingsMenuItem struct {
	title       string
	description string
}

type SettingsModel struct {
	items    []SettingsMenuItem
	selected int
	keys     menuKeyMap // Reuse the same key mappings as main menu
	width    int
	height   int
}

// Message type for settings actions
type SettingsAction int

const (
	EditAPIKey SettingsAction = iota
	EditTheme
	EditSystemPrompt
)

type SettingsMsg struct {
	action SettingsAction
}

func NewSettings() SettingsModel {
	items := []SettingsMenuItem{
		{title: "API Key", description: "Edit your API key"},
		{title: "Theme", description: "Change the application theme"},
		{title: "System Prompt", description: "Manage system prompts"},
	}

	return SettingsModel{
		items:    items,
		selected: 0,
		keys:     menuKeys, // Reuse main menu keys
	}
}

func (m SettingsModel) Init() tea.Cmd {
	return nil
}

func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}

		if msg.String() == "esc" {
			return m, func() tea.Msg { return ChangeViewMsg(menuView) }
		}

		switch {
		case key.Matches(msg, m.keys.Up):
			m.selected--
			if m.selected < 0 {
				m.selected = len(m.items) - 1
			}

		case key.Matches(msg, m.keys.Down):
			m.selected++
			if m.selected >= len(m.items) {
				m.selected = 0
			}

		case key.Matches(msg, m.keys.Select):
			switch m.items[m.selected].title {
			case "API Key":
				return m, func() tea.Msg { return SettingsMsg{action: EditAPIKey} }
			case "Theme":
				return m, func() tea.Msg { return SettingsMsg{action: EditTheme} }
			case "System Prompt":
				return m, func() tea.Msg { return SettingsMsg{action: EditSystemPrompt} }
			}
		}
	}

	return m, nil
}

func (m SettingsModel) View() string {
	// Create styles (reusing the same styles as main menu)
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		Padding(1, 0).
		Align(lipgloss.Center)

	menuStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(2, 4).
		Width(60).
		Align(lipgloss.Center)

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99"))

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("246"))

	descriptionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	// Build menu items
	var menuItems string
	for i, item := range m.items {
		menuItem := item.title
		if i == m.selected {
			menuItem = "▸ " + menuItem
			menuItem = selectedStyle.Render(menuItem)
		} else {
			menuItem = "  " + menuItem
			menuItem = normalStyle.Render(menuItem)
		}
		
		menuItem += " " + descriptionStyle.Render(item.description)
		menuItems += menuItem + "\n"
	}

	// Create the menu box
	menu := menuStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Render("Settings"),
			"",
			menuItems,
			"",
			descriptionStyle.Render("↑/↓: navigate • enter: select • esc: back • q: quit"),
		),
	)

	// Create outer container with additional padding
	containerStyle := lipgloss.NewStyle().
		Padding(4, 0).
		Width(m.width).
		Align(lipgloss.Center)

	return containerStyle.Render(menu)
} 