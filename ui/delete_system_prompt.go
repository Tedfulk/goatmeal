package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
)

type DeleteSystemPromptsView struct {
	list     list.Model
	config   *config.Config
	width    int
	height   int
	selected int
}

func NewDeleteSystemPromptsView(cfg *config.Config) *DeleteSystemPromptsView {
	// Create list items from system prompts
	var items []list.Item
	for _, prompt := range cfg.SystemPrompts {
		items = append(items, SystemPromptMenuItem{title: prompt.Title})
	}

	// Setup list
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	
	l := list.New(items, delegate, 22, 30)
	l.Title = "System Prompts"
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

	return &DeleteSystemPromptsView{
		list:     l,
		config:   cfg,
		selected: 0,
	}
}

func (d *DeleteSystemPromptsView) Update(msg tea.Msg) (*DeleteSystemPromptsView, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return d, func() tea.Msg {
				return SetViewMsg{view: "systemprompts"}
			}
		case tea.KeyEnter:
			if len(d.config.SystemPrompts) > 0 {
				selected := d.list.SelectedItem().(SystemPromptMenuItem)
				// Remove the selected prompt
				var newPrompts []config.SystemPrompt
				for _, prompt := range d.config.SystemPrompts {
					if prompt.Title != selected.title {
						newPrompts = append(newPrompts, prompt)
					}
				}
				d.config.SystemPrompts = newPrompts

				// Update the config file
				manager, err := config.NewManager()
				if err == nil {
					manager.SetSystemPrompts(newPrompts)
				}

				// Update the list
				var items []list.Item
				for _, prompt := range d.config.SystemPrompts {
					items = append(items, SystemPromptMenuItem{title: prompt.Title})
				}
				d.list.SetItems(items)
			}
		}
	}

	// First update the list
	d.list, cmd = d.list.Update(msg)
	
	// Then update our selected index to match the list's current selection
	d.selected = d.list.Index()

	return d, cmd
}

func (d DeleteSystemPromptsView) View() string {
	// Left container (list)
	listStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 1).
		Width(34).
		Height(34)

	// Right container (content)
	contentStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 1).
		Width(104).
		Height(34)

	var content string
	if len(d.config.SystemPrompts) > 0 && d.selected < len(d.config.SystemPrompts) {
		content = d.config.SystemPrompts[d.selected].Content
	}

	// Join containers side by side
	containers := lipgloss.JoinHorizontal(
		lipgloss.Left,
		listStyle.Render(d.list.View()),
		contentStyle.Render(content),
	)

	return lipgloss.Place(
		d.width,
		d.height,
		lipgloss.Left,
		lipgloss.Center,
		containers,
	)
}

func (d *DeleteSystemPromptsView) SetSize(width, height int) {
	d.width = width
	d.height = height
	d.list.SetSize(30, height-16)
} 
