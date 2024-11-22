package ui

import (
	"strings"

	"github.com/tedfulk/goatmeal/config"

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

	// Create fixed-width columns for better alignment
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.colors.MenuSelected)).
		Width(20).
		Align(lipgloss.Right)  // Right align keys

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.colors.MenuNormal)).
		Width(40).             // Fixed width for descriptions
		Align(lipgloss.Left)   // Left align descriptions
		//PaddingLeft(2)        // Remove padding since we'll handle spacing in JoinHorizontal

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
		{"ctrl+l", "List conversations"},
		{"ctrl+t", "New conversation"},
		{";", "Go to theme selector"},
		{"#", "Go to image input"},

		// Scrolling
		{"↑/k", "Scroll up"},
		{"↓/j", "Scroll down"},
		
		// Menu Navigation
		{"↑/k", "Previous item"},
		{"↓/j", "Next item"},
		{"enter", "Select item"},
	}

	for _, s := range shortcuts {
		// Join the key and description with 2 spaces between them
		line := lipgloss.JoinHorizontal(
			lipgloss.Center,  // Center align the joined elements
			keyStyle.Render(s.key),
			"  ",  // Add explicit spacing between columns
			descStyle.Render(s.desc),
		)
		content.WriteString(line + "\n")
	}

	return menuStyle.Render(content.String())
} 