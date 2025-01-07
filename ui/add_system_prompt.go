package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
	"github.com/tedfulk/goatmeal/ui/theme"
)

type AddSystemPromptsKeyMap struct {
	Back       key.Binding
	NextField  key.Binding
	PrevField  key.Binding
	Submit     key.Binding
}

var DefaultAddSystemPromptsKeyMap = AddSystemPromptsKeyMap{
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to system prompts"),
	),
	NextField: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next field"),
	),
	PrevField: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "previous field"),
	),
	Submit: key.NewBinding(
		key.WithKeys("ctrl+a"),
		key.WithHelp("ctrl+a", "save"),
	),
}

type AddSystemPromptsView struct {
	titleInput textinput.Model
	content    textarea.Model
	config     *config.Config
	width      int
	height     int
	focused    string // "title" or "content"
	keys       AddSystemPromptsKeyMap
	err        string
}

func NewAddSystemPromptsView(cfg *config.Config) *AddSystemPromptsView {
	ti := textinput.New()
	ti.Placeholder = "Enter prompt title"
	ti.Focus()
	ti.Width = 135

	ta := textarea.New()
	ta.Placeholder = "Enter system prompt content..."
	ta.ShowLineNumbers = false
	ta.SetWidth(138)
	ta.SetHeight(26)

	return &AddSystemPromptsView{
		titleInput: ti,
		content:    ta,
		config:     cfg,
		focused:    "title",
		keys:       DefaultAddSystemPromptsKeyMap,
	}
}

func (a *AddSystemPromptsView) Update(msg tea.Msg) (*AddSystemPromptsView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, a.keys.Back):
			return nil, func() tea.Msg {
				return SetViewMsg{view: "systemprompts"}
			}

		case key.Matches(msg, a.keys.Submit):
			if a.titleInput.Value() == "" {
				a.err = "Title cannot be empty"
				return a, nil
			}
			if a.content.Value() == "" {
				a.err = "Content cannot be empty"
				return a, nil
			}

			// Check for duplicate titles
			for _, prompt := range a.config.SystemPrompts {
				if prompt.Title == a.titleInput.Value() {
					a.err = "A prompt with this title already exists"
					return a, nil
				}
			}

			// Use config manager to save the prompt
			manager, err := config.NewManager()
			if err != nil {
				a.err = "Failed to create config manager: " + err.Error()
				return a, nil
			}

			// Add the system prompt using manager
			if err := manager.AddSystemPrompt(a.titleInput.Value(), a.content.Value()); err != nil {
				a.err = "Failed to save prompt: " + err.Error()
				return a, nil
			}

			// Update our local config with the new config from manager
			a.config = manager.GetConfig()

			// Return to system prompt settings view
			return nil, func() tea.Msg {
				return SetViewMsg{view: "systemprompts"}
			}

		case key.Matches(msg, a.keys.NextField):
			if a.focused == "title" {
				a.focused = "content"
				a.titleInput.Blur()
				a.content.Focus()
			}
			return a, nil

		case key.Matches(msg, a.keys.PrevField):
			if a.focused == "content" {
				a.focused = "title"
				a.content.Blur()
				a.titleInput.Focus()
			}
			return a, nil
		}
	}

	if a.focused == "title" {
		newTi, cmd := a.titleInput.Update(msg)
		a.titleInput = newTi
		cmds = append(cmds, cmd)
	} else {
		newTa, cmd := a.content.Update(msg)
		a.content = newTa
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

func (a *AddSystemPromptsView) View() string {
	titleStyle := theme.BaseStyle.Title.
		Foreground(theme.CurrentTheme.Primary.GetColor())

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.CurrentTheme.Primary.GetColor()).
		Padding(1, 1).
		Width(a.width - 4)

	if a.focused == "title" {
		inputStyle = inputStyle.BorderForeground(theme.CurrentTheme.Border.Active.GetColor())
	}

	textareaStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.CurrentTheme.Primary.GetColor()).
		Padding(1, 1).
		Width(a.width - 4).
		Height(a.height - 10)

	if a.focused == "content" {
		textareaStyle = textareaStyle.BorderForeground(theme.CurrentTheme.Border.Active.GetColor())
	}

	helpStyle := lipgloss.NewStyle().
		Foreground(theme.CurrentTheme.Secondary.GetColor())

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000"))

	help := helpStyle.Render("tab: next field • shift+tab: prev field • ctrl+a: add system prompt • esc: back")
	
	var content string
	if a.err != "" {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("Add System Prompt"),
			inputStyle.Render(a.titleInput.View()),
			textareaStyle.Render(a.content.View()),
			errorStyle.Render(a.err),
			help,
		)
	} else {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render("Add System Prompt"),
			inputStyle.Render(a.titleInput.View()),
			textareaStyle.Render(a.content.View()),
			help,
		)
	}

	return lipgloss.Place(
		a.width,
		a.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (a *AddSystemPromptsView) SetSize(width, height int) {
	a.width = width
	a.height = height
	a.content.SetWidth(width - 6)
	a.content.SetHeight(height - 12)
} 