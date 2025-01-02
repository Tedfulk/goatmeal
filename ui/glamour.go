package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
)

type GlamourMenuItem struct {
	title       string
	description string
	enabled     bool
}

func (i GlamourMenuItem) Title() string       { return i.title }
func (i GlamourMenuItem) Description() string { return i.description }
func (i GlamourMenuItem) FilterValue() string { return i.title }

type GlamourMenu struct {
	list   list.Model
	width  int
	height int
	config *config.Config
}

func NewGlamourMenu(cfg *config.Config) GlamourMenu {
	currentState := cfg.Settings.OutputGlamour
	items := []list.Item{
		GlamourMenuItem{
			title:       fmt.Sprintf("Enable Markdown Formatting %s", map[bool]string{true: "✅", false: ""}[currentState]),
			description: "Format AI responses with Markdown",
			enabled:     currentState,
		},
		GlamourMenuItem{
			title:       fmt.Sprintf("Disable Markdown Formatting %s", map[bool]string{true: "", false: "✅"}[currentState]),
			description: "Show AI responses as plain text",
			enabled:     !currentState,
		},
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

	return GlamourMenu{
		list:   l,
		config: cfg,
	}
}

func (m GlamourMenu) Update(msg tea.Msg) (GlamourMenu, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":
			selected := m.list.SelectedItem().(GlamourMenuItem)
			newState := selected.title == fmt.Sprintf("Enable Markdown Formatting %s", map[bool]string{true: "✅", false: ""}[m.config.Settings.OutputGlamour])
			
			// Update the config
			m.config.Settings.OutputGlamour = newState
			manager, err := config.NewManager()
			if err == nil {
				manager.UpdateSettings(m.config.Settings)
			}

			// Update the list items to reflect the new state
			items := []list.Item{
				GlamourMenuItem{
					title:       fmt.Sprintf("Enable Markdown Formatting %s", map[bool]string{true: "✅", false: ""}[newState]),
					description: "Format AI responses with Markdown",
					enabled:     newState,
				},
				GlamourMenuItem{
					title:       fmt.Sprintf("Disable Markdown Formatting %s", map[bool]string{true: "", false: "✅"}[newState]),
					description: "Show AI responses as plain text",
					enabled:     !newState,
				},
			}
			m.list.SetItems(items)
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m GlamourMenu) View() string {
	menuStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 1).
		Width(52)

	titleStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Width(48).
		Align(lipgloss.Center)

	menuContent := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("Glamour Settings"),
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

func (m *GlamourMenu) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width-4, height-12)
} 