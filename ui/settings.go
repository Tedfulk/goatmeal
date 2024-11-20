package ui

import (
	"goatmeal/config"

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
	colors   config.ThemeColors
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

func NewSettings(colors config.ThemeColors) SettingsModel {
	items := []SettingsMenuItem{
		{title: "API Key", description: "Edit your API key"},
		{title: "Theme", description: "Change application theme"},
		{title: "System Prompt", description: "Manage system prompts"},
	}

	return SettingsModel{
			items:    items,
			selected: 0,
			keys:     menuKeys,
			colors:   colors,
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
	// Create styles with theme colors
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(m.colors.MenuTitle)).
		Padding(1, 0).
		Align(lipgloss.Center)

	menuStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.colors.MenuBorder)).
		Padding(2, 4).
		Width(60).
		Align(lipgloss.Center)

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(m.colors.MenuSelected))

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.colors.MenuNormal))

	descriptionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.colors.MenuDescription))

	// Build menu items
	var menuItems string
	for i, item := range m.items {
		if i == m.selected {
			menuItems += selectedStyle.Render("▸ "+item.title) + "\n"
			menuItems += descriptionStyle.Render("  "+item.description) + "\n\n"
		} else {
			menuItems += normalStyle.Render("  "+item.title) + "\n"
			menuItems += descriptionStyle.Render("  "+item.description) + "\n\n"
		}
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