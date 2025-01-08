package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
	"github.com/tedfulk/goatmeal/services/providers/model_selection"
	"github.com/tedfulk/goatmeal/ui/theme"
)

type ModelProviderMenuItem struct {
	provider string
	isCurrent bool
}

func (i ModelProviderMenuItem) Title() string {
	if i.isCurrent {
		return i.provider + " ✅"
	}
	return i.provider
}

func (i ModelProviderMenuItem) Description() string { 
	if i.isCurrent {
		return "Current provider"
	}
	return "Select to view models" 
}

func (i ModelProviderMenuItem) FilterValue() string { return i.provider }

type Model struct {
	id string
}

func (i Model) Title() string       { return i.id }
func (i Model) Description() string { return "" }
func (i Model) FilterValue() string { return i.id }

type ModelSettings struct {
	providerList list.Model
	modelList    list.Model
	config       *config.Config
	width        int
	height       int
	showModels   bool
}

// Add helper method to refresh provider list
func (m *ModelSettings) refreshProviderList() {
	var items []list.Item
	for provider, key := range m.config.APIKeys {
		if key != "" && provider != "tavily" {
			items = append(items, ModelProviderMenuItem{
				provider:  provider,
				isCurrent: provider == m.config.CurrentProvider,
			})
		}
	}
	m.providerList.SetItems(items)
}

func NewModelSettings(cfg *config.Config) ModelSettings {
	// Create list items for providers with API keys
	var items []list.Item
	for provider, key := range cfg.APIKeys {
		// Skip tavily and only add providers with valid API keys
		if key != "" && provider != "tavily" {
			items = append(items, ModelProviderMenuItem{
				provider:  provider,
				isCurrent: provider == cfg.CurrentProvider,
			})
		}
	}

	providerDelegate := list.NewDefaultDelegate()
	providerList := list.New(items, providerDelegate, 0, 0)
	providerList.Title = ""
	providerList.SetShowHelp(true)
	providerList.SetFilteringEnabled(false)
	providerList.Styles.Title = theme.BaseStyle.Title.
		Foreground(theme.CurrentTheme.Primary.GetColor())

	modelList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	modelList.Title = ""
	modelList.SetShowHelp(true)
	modelList.SetFilteringEnabled(true)
	modelList.Styles.Title = theme.BaseStyle.Title.
		Foreground(theme.CurrentTheme.Primary.GetColor())

	return ModelSettings{
		providerList: providerList,
		modelList:    modelList,
		config:       cfg,
		showModels:   false,
	}
}

// Add ModelChangeMsg for notifying App about model changes
type ModelChangeMsg struct {
	Provider string
	Model    string
}

// Add a message type for model fetching results
type fetchModelsMsg struct {
	models []string
	err    error
}

// Add a command to fetch models
func fetchModels(provider, apiKey string) tea.Cmd {
	return func() tea.Msg {
		models, err := model_selection.FetchModels(provider, apiKey)
		if err != nil {
			return fetchModelsMsg{err: fmt.Errorf("error fetching models: %w", err)}
		}
		return fetchModelsMsg{models: models}
	}
}

func (m ModelSettings) Update(msg tea.Msg) (ModelSettings, tea.Cmd) {
	switch msg := msg.(type) {
	case fetchModelsMsg:
		if msg.err != nil {
			fmt.Println("error fetching models", msg.err)
			return m, nil
		}

		// Create model list items
		items := make([]list.Item, len(msg.models))
		for i, model := range msg.models {
			items[i] = Model{
				id: model,
			}
		}
		m.modelList.SetItems(items)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.showModels {
				m.showModels = false
				m.refreshProviderList() // Refresh when returning to provider view
				return m, nil
			}
			return m, func() tea.Msg {
				return SetViewMsg{view: "settings"}
			}

		case "enter":
			if m.showModels {
				if i, ok := m.modelList.SelectedItem().(Model); ok {
					// Update config with new provider and model
					manager, err := config.NewManager()
					if err == nil {
						manager.SetCurrentModel(i.id)
					}

					// Reset view to show providers
					m.showModels = false
					m.refreshProviderList() // Refresh after model change

					// Return to settings menu and notify about model change
					return m, func() tea.Msg {
						return ModelChangeMsg{
							Provider: m.config.CurrentProvider,
							Model:    i.id,
						}
					}
				}
			} else {
				if i, ok := m.providerList.SelectedItem().(ModelProviderMenuItem); ok {
					// Update config with new provider
					manager, err := config.NewManager()
					if err == nil {
						manager.SetCurrentProvider(i.provider)
						m.config.CurrentProvider = i.provider
					}

					// Show models list and fetch models
					m.showModels = true
					return m, fetchModels(i.provider, m.config.APIKeys[i.provider])
				}
			}
		}
	}

	// Update the appropriate list based on which view is showing
	if m.showModels {
		var listCmd tea.Cmd
		m.modelList, listCmd = m.modelList.Update(msg)
		return m, listCmd
	}

	var listCmd tea.Cmd
	m.providerList, listCmd = m.providerList.Update(msg)
	return m, listCmd
}

func (m ModelSettings) View() string {
	menuStyle := theme.BaseStyle.Menu.
		BorderForeground(theme.CurrentTheme.Primary.GetColor())

	titleStyle := theme.BaseStyle.Title.
		Foreground(theme.CurrentTheme.Primary.GetColor())

	var content string
	if m.showModels {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("Select Model"),
			m.modelList.View(),
		)
	} else {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("Select Provider"),
			m.providerList.View(),
		)
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		menuStyle.Render(content),
	)
}

func (m *ModelSettings) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.providerList.SetSize(width-4, height-12)
	m.modelList.SetSize(width-4, height-12)
}
