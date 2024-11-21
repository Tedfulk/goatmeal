package ui

import (
	"github.com/tedfulk/goatmeal/config"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuModel struct {
	items     []MenuItem
	selected  int
	keys      menuKeyMap
	width     int
	height    int
	colors    config.ThemeColors
	quitting  bool
	config    *config.Config
}

type MenuItem struct {
	title       string
	description string
}

type menuKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding
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
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

func NewMenu(colors config.ThemeColors, cfg *config.Config) MenuModel {
	items := []MenuItem{
		{title: "New Conversation", description: "Start a new chat (ctrl+t)"},
		{title: "Conversations", description: "List conversation history (ctrl+l)"},
		{title: "Settings", description: "Configure settings (shift+tab)"},
		{title: "Help", description: "View keyboard shortcuts"},
		{title: "Quit", description: "(ctrl+c, q)"},
	}

	return MenuModel{
		items:    items,
		selected: 0,
		keys:     menuKeys,
		colors:   colors,
		config:   cfg,
	}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

// Add new message type for view changes
type ChangeViewMsg viewState

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Add ESC handling to return to chat view
		if msg.String() == "esc" {
			return m, func() tea.Msg { return ChangeViewMsg(chatView) }
		}

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
			case "New Conversation":
				return m, func() tea.Msg { return NewChatMsg{} }
			case "Quit":
				m.quitting = true
				return m, tea.Quit
			case "Settings":
				return m, func() tea.Msg { return ChangeViewMsg(settingsView) }
			case "Help":
				return m, func() tea.Msg { return ChangeViewMsg(helpView) }
			case "Conversations":
				return m, func() tea.Msg { return ChangeViewMsg(conversationListView) }
			}
		}
	}

	return m, nil
}

func (m MenuModel) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	// Create styles with theme colors
	welcomeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(m.colors.MenuTitle)).
		Padding(0, 0, 1, 0).  // Remove bottom padding
		Align(lipgloss.Center)

	modelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.colors.MenuDescription)).
		Padding(0, 0, 1, 0).  // Add padding below the model info
		Align(lipgloss.Center)

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

	// Create the menu box without the welcome message
	menu := menuStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Render("Main Menu"),
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

	// Combine welcome message, model info, and menu box
	return containerStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			welcomeStyle.Render("Welcome, " + m.config.Username + "!"),
			modelStyle.Render("Model: " + m.config.DefaultModel),
			menu,
		),
	)
} 