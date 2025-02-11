package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
	"github.com/tedfulk/goatmeal/ui/theme"
)

// SetViewMsg is used to change the current view
type SetViewMsg struct {
	view string
}

type APIKeyMenuItem struct {
	provider string
	hasKey   bool
}

func (i APIKeyMenuItem) Title() string {
	if i.hasKey {
		return fmt.Sprintf("%s ✅", i.provider)
	}
	return i.provider
}

func (i APIKeyMenuItem) Description() string {
	if i.hasKey {
		return "API key configured"
	}
	return "API key needed"
}

func (i APIKeyMenuItem) FilterValue() string { return i.provider }

type APIKeySettings struct {
	list      list.Model
	textInput textinput.Model
	config    *config.Config
	width     int
	height    int
	selected  string
	inputting bool
}

func NewAPIKeySettings(cfg *config.Config) APIKeySettings {
	// Create list items for each provider
	items := []list.Item{
		APIKeyMenuItem{provider: "openai", hasKey: cfg.APIKeys["openai"] != ""},
		APIKeyMenuItem{provider: "anthropic", hasKey: cfg.APIKeys["anthropic"] != ""},
		APIKeyMenuItem{provider: "gemini", hasKey: cfg.APIKeys["gemini"] != ""},
		APIKeyMenuItem{provider: "deepseek", hasKey: cfg.APIKeys["deepseek"] != ""},
		APIKeyMenuItem{provider: "groq", hasKey: cfg.APIKeys["groq"] != ""},
		APIKeyMenuItem{provider: "tavily", hasKey: cfg.APIKeys["tavily"] != ""},
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

	ti := textinput.New()
	ti.Placeholder = "Enter API key"
	ti.Width = 40

	return APIKeySettings{
		list:      l,
		textInput: ti,
		config:    cfg,
	}
}

func (a APIKeySettings) Update(msg tea.Msg) (APIKeySettings, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if a.inputting {
			switch msg.Type {
			case tea.KeyEsc:
				// Return to list view
				a.inputting = false
				a.selected = ""
				a.textInput.Reset()
				return a, nil
			case tea.KeyEnter:
				if a.textInput.Value() != "" {
					// Update the config
					a.config.APIKeys[a.selected] = a.textInput.Value()
					manager, err := config.NewManager()
					if err == nil {
						manager.SetAPIKey(a.selected, a.textInput.Value())
					}

					// Update list items to reflect changes
					items := a.list.Items()
					for i, item := range items {
						if menuItem, ok := item.(APIKeyMenuItem); ok {
							if menuItem.provider == a.selected {
								items[i] = APIKeyMenuItem{
									provider: menuItem.provider,
									hasKey:   true,
								}
							}
						}
					}
					a.list.SetItems(items)
					
					// Reset input state
					a.inputting = false
					a.selected = ""
					a.textInput.Reset()
					return a, nil
				}
			}

			a.textInput, cmd = a.textInput.Update(msg)
			return a, cmd
		}

		switch msg.Type {
		case tea.KeyEnter:
			selected := a.list.SelectedItem().(APIKeyMenuItem)
			a.selected = selected.provider
			a.inputting = true
			a.textInput.Focus()
			return a, textinput.Blink
		case tea.KeyEsc:
			return a, func() tea.Msg {
				return SetViewMsg{view: "settings"}
			}
		}
	}

	if !a.inputting {
		a.list, cmd = a.list.Update(msg)
	}
	return a, cmd
}

func (a APIKeySettings) View() string {
	menuStyle := theme.BaseStyle.Menu.
		BorderForeground(theme.CurrentTheme.Primary.GetColor())

	titleStyle := theme.BaseStyle.Title.
		Foreground(theme.CurrentTheme.Primary.GetColor())

	var content string
	if a.inputting {
		inputStyle := lipgloss.NewStyle().
			Padding(0, 0, 1, 2)

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("API Key Settings"),
			"",
			fmt.Sprintf("Enter API key for %s:", a.selected),
			"",
			inputStyle.Render(a.textInput.View()),
			"",
			"Press Enter to save, Esc to cancel",
		)
	} else {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("API Key Settings"),
			a.list.View(),
		)
	}

	return lipgloss.Place(
		a.width,
		a.height,
		lipgloss.Center,
		lipgloss.Center,
		menuStyle.Render(content),
	)
}

func (a *APIKeySettings) SetSize(width, height int) {
	a.width = width
	a.height = height
	a.list.SetSize(width-4, height-12)
} 