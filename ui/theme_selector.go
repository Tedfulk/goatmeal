package ui

import (
	"github.com/tedfulk/goatmeal/config"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ThemeSelectorModel struct {
	themes    []string
	selected  int
	keys      menuKeyMap
	width     int
	height    int
	colors    config.ThemeColors
}

type ThemeSelectedMsg struct {
	theme string
}

func NewThemeSelector(colors config.ThemeColors) ThemeSelectorModel {
	// Get theme names from config.ThemeMap
	var themes []string
	for theme := range config.ThemeMap {
		 themes = append(themes, theme)
	}

	return ThemeSelectorModel{
		themes:   themes,
		selected: 0,
		keys:     menuKeys,
		colors:   colors,
	}
}

func (m ThemeSelectorModel) Init() tea.Cmd {
	return nil
}

func (m ThemeSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, func() tea.Msg { return ChangeViewMsg(settingsView) }
		}

		switch {
		case key.Matches(msg, m.keys.Up):
			m.selected--
			if m.selected < 0 {
				m.selected = len(m.themes) - 1
			}

		case key.Matches(msg, m.keys.Down):
			m.selected++
			if m.selected >= len(m.themes) {
				m.selected = 0
			}

		case key.Matches(msg, m.keys.Select):
			selectedTheme := m.themes[m.selected]
			return m, func() tea.Msg { return ThemeSelectedMsg{theme: selectedTheme} }
		}
	}

	return m, nil
}

func (m ThemeSelectorModel) View() string {
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

	// Build theme items
	var menuItems string
	for i, theme := range m.themes {
		menuItem := theme
		if i == m.selected {
			menuItem = "▸ " + menuItem
			menuItem = selectedStyle.Render(menuItem)
		} else {
			menuItem = "  " + menuItem
			menuItem = normalStyle.Render(menuItem)
		}
		menuItems += menuItem + "\n"
	}

	// Create the menu box
	menu := menuStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Render("Select Theme"),
			"",
			menuItems,
			"",
			normalStyle.Render("↑/↓: navigate • enter: select • esc: back • q: quit"),
		),
	)

	// Create outer container with additional padding
	containerStyle := lipgloss.NewStyle().
		Padding(4, 0).
		Width(m.width).
		Align(lipgloss.Center)

	return containerStyle.Render(menu)
} 