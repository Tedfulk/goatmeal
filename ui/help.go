package ui

import (
	"github.com/tedfulk/goatmeal/config"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HelpModel struct {
	width  int
	height int
	colors config.ThemeColors
}

func NewHelp(colors config.ThemeColors) HelpModel {
	return HelpModel{
		colors: colors,
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
	// Create styles with theme colors
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(m.colors.MenuTitle)).
		Padding(1, 0).
		Align(lipgloss.Center)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.colors.MenuSelected)).
		Width(20).
		Align(lipgloss.Left)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.colors.MenuNormal)).
		PaddingLeft(2)

	menuStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.colors.MenuBorder)).
		Padding(2, 4).
		Width(80).
		Align(lipgloss.Center)

	// Build help content
	var content strings.Builder
	content.WriteString(titleStyle.Render("Keyboard Shortcuts"))
	content.WriteString("\n\n")

	shortcuts := []struct{ key, desc string }{
		// Navigation
		{"shift+tab", "Toggle menu"},
		{"tab", "Toggle focus (in chat/list view)"},
		{"esc", "Back/Quit"},
		{"q", "Quit"},
		
		// Chat Actions
		{"enter", "Send message"},
		{"shift+enter", "New line in message"},
		{"shift+:", "New conversation"},
		{"ctrl+l", "View conversations"},
		
		// Scrolling
		{"↑/k", "Scroll up"},
		{"↓/j", "Scroll down"},
		{"pgup/ctrl+b", "Page up"},
		{"pgdn/ctrl+f", "Page down"},
		{"home", "Scroll to top"},
		{"end", "Scroll to bottom"},
		
		// Menu Navigation
		{"↑/k", "Previous item"},
		{"↓/j", "Next item"},
		{"enter", "Select item"},
	}

	for _, s := range shortcuts {
		line := lipgloss.JoinHorizontal(
			lipgloss.Left,
			keyStyle.Render(s.key),
			descStyle.Render(s.desc),
		)
		content.WriteString(line + "\n")
	}

	return menuStyle.Render(content.String())
} 