package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuItem struct {
	title       string
	description    string
}

func (i MenuItem) Title() string       { return i.title }
func (i MenuItem) Description() string { return i.description }
func (i MenuItem) FilterValue() string { return i.title }

type Menu struct {
	list     list.Model
	width    int
	height   int
	selected bool
}

func NewMenu() Menu {
	items := []list.Item{
		MenuItem{
			title:       "New Conversation",
			description:    "Start a new chat (ctrl+t)",
		},
		MenuItem{
			title:       "Conversations",
			description:    "List conversation history (ctrl+l)",
		},
		MenuItem{
			title:       "Settings",
			description:    "Configure settings (ctrl+s)",
		},
		MenuItem{
			title:       "Help",
			description:    "View keyboard shortcuts (ctrl+h)",
		},
		MenuItem{
			title:       "Quit",
			description:    "(ctrl+c, q)",
		},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = ""
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Padding(0, 0, 1, 2)

	return Menu{
		list: l,
	}
}

func (m Menu) Update(msg tea.Msg) (Menu, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Menu) View() string {
	// Create a style for the menu container
	menuStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 2).
		Width(50)

	// Create the menu content with centered title
	titleStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Width(46).
		Align(lipgloss.Center)

	menuContent := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Main Menu"),
		m.list.View(),
	)

	// Center the menu in the window
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		menuStyle.Render(menuContent),
	)
}

func (m *Menu) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width-4, height-12)
} 