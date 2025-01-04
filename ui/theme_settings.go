package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
	"github.com/tedfulk/goatmeal/ui/theme"
)

type ThemeMenuItem struct {
	title       string
	description string
}

func (i ThemeMenuItem) Title() string       { return i.title }
func (i ThemeMenuItem) Description() string { return i.description }
func (i ThemeMenuItem) FilterValue() string { return i.title }

type ThemeSettings struct {
	list   list.Model
	config *config.Config
	width  int
	height int
}

// ThemeChangeMsg is sent when the theme is changed
type ThemeChangeMsg struct {
	Theme theme.Theme
}

func NewThemeSettings(cfg *config.Config) ThemeSettings {
	items := []list.Item{
		ThemeMenuItem{title: "Default", description: "Default purple theme"},
		ThemeMenuItem{title: "Dracula", description: "Dark theme with vibrant colors"},
		ThemeMenuItem{title: "Nord", description: "Cool, bluish dark theme"},
		ThemeMenuItem{title: "Matrix Classic", description: "Classic green Matrix theme"},
		ThemeMenuItem{title: "Matrix Neo", description: "Modern Matrix theme with blue accents"},
		ThemeMenuItem{title: "Cyberpunk Neon", description: "Vibrant cyberpunk theme with neon colors"},
		ThemeMenuItem{title: "Cyberpunk Red", description: "Dark cyberpunk theme with red accents"},
		ThemeMenuItem{title: "Python", description: "Python's official blue and yellow colors"},
		ThemeMenuItem{title: "Monochrome", description: "Clean black and white theme"},
		ThemeMenuItem{title: "Rainbow Bright", description: "Vibrant rainbow colors"},
		ThemeMenuItem{title: "Rainbow Pastel", description: "Soft pastel rainbow colors"},
		ThemeMenuItem{title: "Barbie", description: "Pink-themed Barbie movie inspired"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Theme Settings"
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
		Foreground(theme.CurrentTheme.Primary.GetColor()).
		Align(lipgloss.Center).
		Width(46)

	return ThemeSettings{
		list:   l,
		config: cfg,
	}
}

func (t ThemeSettings) Update(msg tea.Msg) (ThemeSettings, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected := t.list.SelectedItem().(ThemeMenuItem)
			// Update the config
			t.config.Settings.Theme.Name = selected.title
			manager, err := config.NewManager()
			if err == nil {
				// Save the theme selection to config file
				manager.UpdateSettings(t.config.Settings)
			}

			// Apply the new theme immediately
			var newTheme theme.Theme
			switch selected.title {
			case "Default":
				newTheme = theme.DefaultTheme
			case "Dracula":
				newTheme = theme.DraculaTheme
			case "Nord":
				newTheme = theme.NordTheme
			case "Matrix Classic":
				newTheme = theme.MatrixClassicTheme
			case "Matrix Neo":
				newTheme = theme.MatrixNeoTheme
			case "Cyberpunk Neon":
				newTheme = theme.CyberpunkNeonTheme
			case "Cyberpunk Red":
				newTheme = theme.CyberpunkRedTheme
			case "Python":
				newTheme = theme.PythonTheme
			case "Monochrome":
				newTheme = theme.MonochromeTheme
			case "Rainbow Bright":
				newTheme = theme.RainbowBrightTheme
			case "Rainbow Pastel":
				newTheme = theme.RainbowPastelTheme
			case "Barbie":
				newTheme = theme.BarbieTheme
			default:
				newTheme = theme.DefaultTheme
			}

			// Update the current theme
			theme.CurrentTheme = newTheme

			// Return a command to notify the app about the theme change
			return t, func() tea.Msg {
				return ThemeChangeMsg{Theme: newTheme}
			}
		}
	}

	t.list, cmd = t.list.Update(msg)
	return t, cmd
}

func (t ThemeSettings) View() string {
	menuStyle := theme.BaseStyle.Menu.
		BorderForeground(theme.CurrentTheme.Primary.GetColor())

	return lipgloss.Place(
		t.width,
		t.height,
		lipgloss.Center,
		lipgloss.Center,
		menuStyle.Render(t.list.View()),
	)
}

func (t *ThemeSettings) SetSize(width, height int) {
	t.width = width
	t.height = height
	t.list.SetSize(width-4, height-12)
} 