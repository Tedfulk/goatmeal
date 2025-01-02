package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SettingsMenuItem struct {
	title       string
	description string
}

func (i SettingsMenuItem) Title() string       { return i.title }
func (i SettingsMenuItem) Description() string { return i.description }
func (i SettingsMenuItem) FilterValue() string { return i.title }

type SettingsMenu struct {
	list        list.Model
	width       int
	height      int
	currentView string
}

func NewSettingsMenu() SettingsMenu {
	items := []list.Item{
		SettingsMenuItem{title: "API Keys", description: "Configure API keys for providers"},
		SettingsMenuItem{title: "System Prompts", description: "Manage system prompts"},
		SettingsMenuItem{title: "Theme (TODO)", description: "Change application theme"},
		SettingsMenuItem{title: "Glamour", description: "Configure markdown formatting"},
		SettingsMenuItem{title: "Username", description: "Change your username"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = ""
	l.SetShowHelp(true)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("esc", "back"),
			),
		}
	}
	l.AdditionalFullHelpKeys = l.AdditionalShortHelpKeys
	l.SetFilteringEnabled(false)
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Padding(0, 0, 1, 2)

	return SettingsMenu{
		list: l,
		currentView: "settings",
	}
}

func (s SettingsMenu) Update(msg tea.Msg) (SettingsMenu, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			selected := s.list.SelectedItem().(SettingsMenuItem)
			switch selected.title {
			case "API Keys":
				s.currentView = "apikeys"
			case "System Prompts":
				s.currentView = "systemprompts"
			case "Theme (TODO)":
				// TODO: Implement theme settings
				s.currentView = "settings"
			case "Glamour":
				s.currentView = "glamour"
			case "Username":
				s.currentView = "username"
			}
			return s, nil
		}
	}

	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (m SettingsMenu) View() string {
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
		titleStyle.Render("Settings"),
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

func (m *SettingsMenu) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width-4, height-12)
} 