package ui

import (
	"github.com/tedfulk/goatmeal/config"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ModelSelectorModel struct {
	models    []string
	selected  int
	keys      menuKeyMap
	width     int
	height    int
	colors    config.ThemeColors
	loading   bool
	err       error
}

type ModelSelectedMsg struct {
	model string
}

func NewModelSelector(colors config.ThemeColors, cfg *config.Config) (ModelSelectorModel, error) {
	// Fetch available models
	models, err := config.FetchModels(cfg.APIKey)
	if err != nil {
		return ModelSelectorModel{}, err
	}

	return ModelSelectorModel{
		models:   models,
		selected: 0,
		keys:     menuKeys,
		colors:   colors,
	}, nil
}

func (m ModelSelectorModel) Init() tea.Cmd {
	return nil
}

func (m ModelSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				m.selected = len(m.models) - 1
			}

		case key.Matches(msg, m.keys.Down):
			m.selected++
			if m.selected >= len(m.models) {
				m.selected = 0
			}

		case key.Matches(msg, m.keys.Select):
			selectedModel := m.models[m.selected]
			return m, func() tea.Msg { return ModelSelectedMsg{model: selectedModel} }
		}
	}

	return m, nil
}

func (m ModelSelectorModel) View() string {
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

	// Build model items
	var menuItems string
	for i, model := range m.models {
		menuItem := model
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
			titleStyle.Render("Select Model"),
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