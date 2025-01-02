package setup

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
)

// Stage represents a step in the setup wizard
type Stage int

const (
	UsernameStage Stage = iota
	ProvidersStage
	ModelSelectionStage
)

var (
	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Padding(1, 0, 1, 0)

	boxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Width(50)
)

// Wizard represents the setup wizard
type Wizard struct {
	config         *config.Config
	stage          Stage
	username       UsernameInput
	providers      ProviderList
	modelSelection ModelSelection
	width          int
	height         int
	quitting       bool
}

// NewWizard creates a new setup wizard
func NewWizard(cfg *config.Config) *Wizard {
	w := &Wizard{
		config:         cfg,
		stage:          UsernameStage,
		username:       NewUsernameInput(),
		providers:      NewProviderList(),
		modelSelection: NewModelSelection(cfg),
	}
	return w
}

// Init initializes the wizard
func (w *Wizard) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles updates for the wizard
func (w *Wizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle window size changes
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		w.width = msg.Width
		w.height = msg.Height
	}

	switch w.stage {
	case UsernameStage:
		model, newCmd := w.username.Update(msg)
		if username, ok := model.(UsernameInput); ok {
			w.username = username
			if username.Done() {
				w.config.Settings.Username = username.Value()
				w.stage = ProvidersStage
			}
		}
		cmd = newCmd

	case ProvidersStage:
		model, newCmd := w.providers.Update(msg)
		if providers, ok := model.(ProviderList); ok {
			w.providers = providers
			if providers.Done() {
				// Save API keys
				apiKeys := providers.GetConfiguredProviders()
				if w.config.APIKeys == nil {
					w.config.APIKeys = make(map[string]string)
				}
				for provider, key := range apiKeys {
					w.config.APIKeys[provider] = key
				}
				w.stage = ModelSelectionStage
				w.modelSelection.SetProviders(apiKeys)
			}
		}
		cmd = newCmd

	case ModelSelectionStage:
		model, newCmd := w.modelSelection.Update(msg)
		if modelSelection, ok := model.(ModelSelection); ok {
			w.modelSelection = modelSelection
			if modelSelection.Done() {
				// Save selected provider and model
				provider, model := modelSelection.GetSelected()
				w.config.CurrentProvider = provider
				w.config.CurrentModel = model
				
				// Save all configuration changes
				manager, err := config.NewManager()
				if err == nil {
					manager.UpdateSettings(w.config.Settings)
					for provider, key := range w.config.APIKeys {
						manager.SetAPIKey(provider, key)
					}
					manager.SetCurrentProvider(w.config.CurrentProvider)
					manager.SetCurrentModel(w.config.CurrentModel)
					manager.SetCurrentSystemPrompt(w.config.SystemPrompts[0].Content)
				}
				
				w.quitting = true
				return w, tea.Quit
			}
		}
		cmd = newCmd
	}

	return w, cmd
}

// View renders the wizard
func (w *Wizard) View() string {
	if w.quitting {
		return "Setup complete!\n"
	}

	var content string
	switch w.stage {
	case UsernameStage:
		content = w.username.View()
	case ProvidersStage:
		content = w.providers.View()
	case ModelSelectionStage:
		content = w.modelSelection.View()
	}

	// Center all content in the terminal
	return lipgloss.Place(
		w.width,
		w.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

// Done returns whether the wizard is complete
func (w *Wizard) Done() bool {
	return w.quitting
}

// Username returns the entered username
func (w *Wizard) Username() string {
	return w.username.Value()
}

// APIKeys returns the configured API keys
func (w *Wizard) APIKeys() map[string]string {
	return w.providers.GetConfiguredProviders()
}

// GetSelectedProvider returns the selected provider and model
func (w *Wizard) GetSelectedProvider() (provider, model string) {
	return w.modelSelection.GetSelected()
} 