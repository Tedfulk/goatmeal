package setup

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
	"github.com/tedfulk/goatmeal/services/providers/model_selection"
)

// ModelSelection represents the model selection stage
type ModelSelection struct {
	config          *config.Config
	providerList    list.Model
	modelList       list.Model
	providers       map[string]string
	selectedProvider string
	selectedModel   string
	loading         bool
	showModels      bool
	done            bool
	err             error
}

// Model represents an AI model
type Model struct {
	id string
}

// FilterValue implements list.Item interface
func (m Model) FilterValue() string { return m.id }

// Title implements list.Item interface
func (m Model) Title() string { return m.id }

// Description implements list.Item interface
func (m Model) Description() string { return "" }

// Provider represents a configured provider
type ConfiguredProvider struct {
	name string
	key  string
}

// FilterValue implements list.Item interface
func (p ConfiguredProvider) FilterValue() string { return p.name }

// Title implements list.Item interface
func (p ConfiguredProvider) Title() string { return p.name }

// Description implements list.Item interface
func (p ConfiguredProvider) Description() string { return "Select to view available models" }

// NewModelSelection creates a new model selection stage
func NewModelSelection(cfg *config.Config) ModelSelection {
	providerList := list.New([]list.Item{}, list.NewDefaultDelegate(), 30, 25)
	providerList.Title = "Select a Provider to View Available Models"
	providerList.SetShowHelp(true)
	providerList.SetFilteringEnabled(false)
	providerList.SetShowStatusBar(false)

	modelList := list.New([]list.Item{}, list.NewDefaultDelegate(), 30, 25)
	modelList.Title = "Select a Model"
	modelList.SetShowHelp(true)
	modelList.SetFilteringEnabled(true)
	modelList.SetShowStatusBar(false)

	return ModelSelection{
		config:       cfg,
		providerList: providerList,
		modelList:    modelList,
	}
}

// SetProviders sets the available providers
func (m *ModelSelection) SetProviders(providers map[string]string) {
	m.providers = providers
	items := make([]list.Item, 0, len(providers))
	for name, key := range providers {
		items = append(items, ConfiguredProvider{name: name, key: key})
	}
	m.providerList.SetItems(items)
}

// Init initializes the model selection
func (m ModelSelection) Init() tea.Cmd {
	return nil
}

// fetchModelsMsg represents the result of fetching models
type fetchModelsMsg struct {
	models []string
	err    error
}

// fetchModels fetches available models for the selected provider
func fetchModels(provider, apiKey string) tea.Cmd {
	return func() tea.Msg {
		// Special handling for Ollama which doesn't require an API key
		if provider == "ollama" {
			models, err := model_selection.FetchModels(provider, "")
			if err != nil {
				return fetchModelsMsg{err: fmt.Errorf("error fetching Ollama models (is Ollama running?): %w", err)}
			}
			return fetchModelsMsg{models: models}
		}

		// Normal flow for other providers
		models, err := model_selection.FetchModels(provider, apiKey)
		if err != nil {
			return fetchModelsMsg{err: fmt.Errorf("error fetching models: %w", err)}
		}
		return fetchModelsMsg{models: models}
	}
}

// Update handles updates for the model selection
func (m ModelSelection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.providerList.SetWidth(msg.Width)
		m.providerList.SetHeight(msg.Height)
		m.modelList.SetWidth(msg.Width)
		m.modelList.SetHeight(msg.Height)
		return m, nil

	case fetchModelsMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		items := make([]list.Item, len(msg.models))
		for i, model := range msg.models {
			items[i] = Model{id: model}
		}
		m.modelList.SetItems(items)
		return m, nil

	case tea.KeyMsg:
		if m.showModels {
			switch msg.String() {
			case "esc":
				m.showModels = false
				m.err = nil
				return m, nil
			case "enter":
				if i, ok := m.modelList.SelectedItem().(Model); ok {
					m.selectedModel = i.id
					m.done = true
					return m, nil
				}
			}
			m.modelList, cmd = m.modelList.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "q", "ctrl+c":
			m.done = true
			return m, tea.Quit
		case "enter":
			if i, ok := m.providerList.SelectedItem().(ConfiguredProvider); ok {
				m.selectedProvider = i.name
				m.showModels = true
				m.loading = true
				return m, fetchModels(i.name, i.key)
			}
		}
		m.providerList, cmd = m.providerList.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the model selection
func (m ModelSelection) View() string {
	if m.err != nil {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Render("Error"),
			boxStyle.Render(
				lipgloss.JoinVertical(
					lipgloss.Center,
					"Error fetching models:",
					"",
					m.err.Error(),
					"",
					"Press ESC to go back",
				),
			),
		)
	}

	if m.loading {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Render("Loading Models"),
			boxStyle.Render(
				lipgloss.JoinVertical(
					lipgloss.Center,
					"Fetching available models...",
					"Please wait...",
				),
			),
		)
	}

	if m.showModels {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Render(fmt.Sprintf("Available Models for %s", m.selectedProvider)),
			boxStyle.Render(
				lipgloss.JoinVertical(
					lipgloss.Center,
					"Select a model to use (type to filter):",
					"",
					m.modelList.View(),
				),
			),
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render("Model Selection"),
		boxStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				"Select a provider to view available models:",
				"",
				m.providerList.View(),
			),
		),
	)
}

// Done returns whether the model selection is complete
func (m ModelSelection) Done() bool {
	return m.done
}

// GetSelected returns the selected provider and model
func (m ModelSelection) GetSelected() (provider, model string) {
	return m.selectedProvider, m.selectedModel
} 