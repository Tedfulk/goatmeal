package setup

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Provider represents an AI provider
type Provider struct {
	name    string
	hasKey  bool
	apiKey  string
}

// FilterValue implements list.Item interface
func (p Provider) FilterValue() string { return p.name }

// Title implements list.Item interface
func (p Provider) Title() string { 
	if p.hasKey {
		return p.name + "✅"
	}
	return p.name
}

// Description implements list.Item interface
func (p Provider) Description() string {
	if p.name == "Continue" {
		return "Continue to next step"
	}
	if p.hasKey {
		return "API key configured"
	}
	return "API key needed"
}

// ProviderList represents the provider selection and API key input stage
type ProviderList struct {
	list          list.Model
	providers     []Provider
	textInput     textinput.Model
	selectedIndex int
	inputting     bool
	done          bool
	width         int
	height        int
}

// NewProviderList creates a new provider list
func NewProviderList() ProviderList {
	providers := []Provider{
		{name: "openai"},
		{name: "anthropic"},
		{name: "gemini"},
		{name: "deepseek"},
		{name: "groq"},
		{name: "tavily"},
	}

	l := list.New(toItems(providers), list.NewDefaultDelegate(), 30, 28)
	l.Title = "Select a Provider"
	l.SetShowHelp(false)

	ti := textinput.New()
	ti.Placeholder = "Enter API key"
	ti.Width = 40

	return ProviderList{
		list:      l,
		providers: providers,
		textInput: ti,
	}
}

// Focus focuses the component
func (m *ProviderList) Focus() {
	m.list.SetShowHelp(true)
}

// Init initializes the provider list
func (m ProviderList) Init() tea.Cmd {
	return nil
}

// Update handles updates for the provider list
func (m ProviderList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width - 4)  // Account for padding
		m.list.SetHeight(msg.Height - 8) // Account for title and padding
		return m, nil

	case tea.KeyMsg:
		if m.inputting {
			switch msg.Type {
			case tea.KeyEnter:
				if m.textInput.Value() != "" {
					// Save API key
					m.providers[m.selectedIndex].apiKey = m.textInput.Value()
					m.providers[m.selectedIndex].hasKey = true
					m.inputting = false
					m.list.SetItems(toItems(m.providers))

					// Add "Continue" option if at least one key is configured
					if !hasSkipOption(m.providers) && hasAnyKey(m.providers) {
						m.providers = append(m.providers, Provider{name: "Continue"})
						m.list.SetItems(toItems(m.providers))
					}

					return m, nil
				}
			case tea.KeyEsc:
				m.inputting = false
				return m, nil
			}

			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "ctrl+c":
			m.done = true
			return m, tea.Quit
		case "enter":
			selected := m.list.SelectedItem().(Provider)
			if selected.name == "Continue" {
				m.done = true
				return m, nil
			}
			
			m.selectedIndex = m.list.Index()
			if !m.providers[m.selectedIndex].hasKey {
				m.inputting = true
				m.textInput.Reset()
				m.textInput.Focus()
				return m, textinput.Blink
			}
		}
	}

	var listCmd tea.Cmd
	m.list, listCmd = m.list.Update(msg)
	return m, listCmd
}

// View renders the provider list
func (m ProviderList) View() string {
	var content string
	if m.inputting {
		content = lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Render("API Key Setup"),
			boxStyle.Render(
				lipgloss.JoinVertical(
					lipgloss.Center,
					"Enter API key for "+m.providers[m.selectedIndex].name+":",
					"",
					m.textInput.View(),
				),
			),
		)
	} else {
		content = lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Render("Provider Setup"),
			boxStyle.Render(m.list.View()),
		)
	}

	// Center the content in the available space
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

// Done returns whether the provider setup is complete
func (m ProviderList) Done() bool {
	return m.done
}

// GetConfiguredProviders returns the list of providers with their API keys
func (m ProviderList) GetConfiguredProviders() map[string]string {
	result := make(map[string]string)
	for _, p := range m.providers {
		if p.hasKey {
			result[p.name] = p.apiKey
		}
	}
	return result
}

// Helper functions
func toItems(providers []Provider) []list.Item {
	items := make([]list.Item, len(providers))
	for i, p := range providers {
		items[i] = p
	}
	return items
}

func hasSkipOption(providers []Provider) bool {
	for _, p := range providers {
		if p.name == "Continue" {
			return true
		}
	}
	return false
}

func hasAnyKey(providers []Provider) bool {
	for _, p := range providers {
		if p.hasKey {
			return true
		}
	}
	return false
} 