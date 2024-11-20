package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HelpModel struct {
	keys     []KeyBinding
	width    int
	height   int
}

type KeyBinding struct {
	key         string
	description string
}

func NewHelp() HelpModel {
	return HelpModel{
		keys: []KeyBinding{
			{"ctrl+n", "Start new conversation"},
			{"ctrl+l", "Show conversation list"},
			{"shift+tab", "Toggle menu"},
			{"tab", "Switch focus between chat and input"},
			{"enter", "Send message"},
			{"shift+enter", "New line in input"},
			{"ctrl+c, esc", "Quit"},
			{"↑/k", "Scroll up in chat/Navigate up in menus"},
			{"↓/j", "Scroll down in chat/Navigate down in menus"},
			{"pgup", "Page up in chat"},
			{"pgdn", "Page down in chat"},
			{"home", "Scroll to top"},
			{"end", "Scroll to bottom"},
			{"/", "Filter conversations (in list view)"},
			{"esc", "Back to previous view"},
		},
	}
}

func (m HelpModel) Init() tea.Cmd {
	return nil
}

func (m HelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "esc":
			return m, func() tea.Msg { return ChangeViewMsg(menuView) }
		}
	}
	return m, nil
}

func (m HelpModel) View() string {
	// Create styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99")).
		Padding(1, 0).
		Align(lipgloss.Center)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Width(20).
		Align(lipgloss.Left)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		PaddingLeft(2)

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Width(70).
		Padding(1, 0).
		Align(lipgloss.Center)

	menuStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Padding(2, 4).
		Width(80).
		Align(lipgloss.Center)

	// Build help content
	var menuItems string
	for _, kb := range m.keys {
		line := lipgloss.NewStyle().
			Width(70).
			Align(lipgloss.Left).
			Render(
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					keyStyle.Render(kb.key),
					descStyle.Render(kb.description),
				),
			)
		menuItems += line + "\n"
	}

	// Create the menu box
	menu := menuStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Render("Keyboard Shortcuts"),
			"",
			menuItems,
			"",
			footerStyle.Render("Press 'esc' to go back"),
		),
	)

	return menu
} 