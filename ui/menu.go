package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type menuKeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Select  key.Binding
	Back    key.Binding
	Quit    key.Binding
}

var menuKeys = menuKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
}

type MenuItem struct {
	title       string
	description string
}

type MenuModel struct {
	items       []MenuItem
	selected    int
	keys        menuKeyMap
	quitting    bool
	width       int
	height      int
}

func NewMenu() MenuModel {
	items := []MenuItem{
		{title: "Settings", description: "Configure application settings"},
		{title: "Help", description: "View keyboard shortcuts and help"},
		{title: "Quit", description: "Exit the application"},
	}

	return MenuModel{
		items:    items,
		selected: 0,
		keys:     menuKeys,
	}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit

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
			case "Quit":
				m.quitting = true
				return m, tea.Quit
			case "Settings":
				// TODO: Handle settings
				return m, nil
			case "Help":
				// TODO: Handle help
				return m, nil
			}
		}
	}

	return m, nil
}

func (m MenuModel) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	// Create styles
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
		
		// Add description
		menuItem += " " + descriptionStyle.Render(item.description)
		menuItems += menuItem + "\n"
	}

	// Create the menu box
	menu := menuStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Render("Main Menu"),
			"",
			menuItems,
			"",
			descriptionStyle.Render("↑/↓: navigate • enter: select • esc: back"),
		),
	)

	// Create outer container with additional padding
	containerStyle := lipgloss.NewStyle().
		Padding(4, 0).
		Width(m.width).
		Align(lipgloss.Center)

	// Center in terminal with outer container
	return containerStyle.Render(menu)
} 