package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
	"github.com/tedfulk/goatmeal/ui/theme"
)

type SettingsMenuItem struct {
	title       string
	description string
}

func (i SettingsMenuItem) Title() string       { return i.title }
func (i SettingsMenuItem) Description() string { return i.description }
func (i SettingsMenuItem) FilterValue() string { return i.title }

type SettingsMenu struct {
	list           list.Model
	width          int
	height         int
	currentView    string
	config         *config.Config
	themeSettings  ThemeSettings
}

func NewSettingsMenu(cfg *config.Config) SettingsMenu {
	items := []list.Item{
		SettingsMenuItem{title: "API Keys", description: "Configure API keys for providers"},
		SettingsMenuItem{title: "System Prompts", description: "Manage system prompts"},
		SettingsMenuItem{title: "Theme", description: "Change application theme"},
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
	l.Styles.Title = theme.BaseStyle.Title.
		Foreground(theme.CurrentTheme.Primary.GetColor())

	return SettingsMenu{
		list:          l,
		currentView:   "settings",
		config:        cfg,
		themeSettings: NewThemeSettings(cfg),
	}
}

func (s SettingsMenu) Update(msg tea.Msg) (SettingsMenu, tea.Cmd) {
	var cmd tea.Cmd

	switch s.currentView {
	case "theme":
		var themeCmd tea.Cmd
		s.themeSettings, themeCmd = s.themeSettings.Update(msg)
		return s, themeCmd
	default:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				selected := s.list.SelectedItem().(SettingsMenuItem)
				switch selected.title {
				case "API Keys":
					s.currentView = "apikeys"
				case "System Prompts":
					s.currentView = "systemprompts"
				case "Theme":
					s.currentView = "theme"
				case "Glamour":
					s.currentView = "glamour"
				case "Username":
					s.currentView = "username"
				}
				return s, nil
			case "esc":
				if s.currentView != "settings" {
					s.currentView = "settings"
					return s, nil
				}
			}
		}

		s.list, cmd = s.list.Update(msg)
		return s, cmd
	}
}

func (m SettingsMenu) View() string {
	switch m.currentView {
	case "theme":
		return m.themeSettings.View()
	default:
		menuStyle := theme.BaseStyle.Menu.
			BorderForeground(theme.CurrentTheme.Primary.GetColor())

		titleStyle := theme.BaseStyle.Title.
			Foreground(theme.CurrentTheme.Primary.GetColor())

		menuContent := lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("Settings"),
			m.list.View(),
		)

		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			menuStyle.Render(menuContent),
		)
	}
}

func (m *SettingsMenu) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width-4, height-12)
	m.themeSettings.SetSize(width, height)
} 