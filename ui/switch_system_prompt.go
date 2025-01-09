package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
	"github.com/tedfulk/goatmeal/ui/theme"
)

type SwitchSystemPromptsKeyMap struct {
	Back        key.Binding
	SwitchFocus key.Binding
}

var DefaultSwitchSystemPromptsKeyMap = SwitchSystemPromptsKeyMap{
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	SwitchFocus: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch focus"),
	),
}

type SwitchSystemPromptsView struct {
	list     list.Model
	viewport viewport.Model
	config   *config.Config
	width    int
	height   int
	selected int
	focused  string
	keys     SwitchSystemPromptsKeyMap
}

type SystemPromptChangeMsg struct {
	NewPrompt string
}

func NewSwitchSystemPromptsView(cfg *config.Config) *SwitchSystemPromptsView {
	// Reload the config to get any newly added prompts
	if newConfig, err := config.Load(); err == nil {
		cfg = newConfig
	}

	var items []list.Item
	for _, prompt := range cfg.SystemPrompts {
		title := prompt.Title
		if prompt.Content == cfg.CurrentSystemPrompt {
			title = prompt.Title + " ✅"
		}
		items = append(items, SystemPromptMenuItem{
			title:       title,
			description: prompt.Content,
		})
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetHeight(1)
	
	l := list.New(items, delegate, 30, 30)
	l.Title = "Switch System Prompts"
	l.SetShowHelp(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = theme.BaseStyle.Title.
		Foreground(theme.CurrentTheme.Primary.GetColor()).
		Align(lipgloss.Center).
		Width(30)

	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			DefaultSwitchSystemPromptsKeyMap.Back,
			DefaultSwitchSystemPromptsKeyMap.SwitchFocus,
		}
	}
	l.AdditionalFullHelpKeys = l.AdditionalShortHelpKeys

	vp := viewport.New(102, 32)
	if len(items) > 0 {
		vp.SetContent(items[0].(SystemPromptMenuItem).description)
	}

	return &SwitchSystemPromptsView{
		list:     l,
		viewport: vp,
		config:   cfg,
		focused:  "list",
		keys:     DefaultSwitchSystemPromptsKeyMap,
	}
}

func (s *SwitchSystemPromptsView) Update(msg tea.Msg) (*SwitchSystemPromptsView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, s.keys.SwitchFocus) {
			if s.focused == "list" {
				s.focused = "viewport"
			} else {
				s.focused = "list"
			}
			return s, nil
		}

		if key.Matches(msg, s.keys.Back) {
			return nil, nil
		}

		if msg.String() == "enter" && s.focused == "list" {
			selected := s.list.SelectedItem().(SystemPromptMenuItem)

			// Update the current system prompt
			manager, err := config.NewManager()
			if err == nil {
				manager.SetCurrentSystemPrompt(selected.description)
				s.config.CurrentSystemPrompt = selected.description

				// Update list items to show new checkmark
				var items []list.Item
				for _, prompt := range s.config.SystemPrompts {
					itemTitle := prompt.Title
					if prompt.Content == selected.description {
						itemTitle = prompt.Title + " ✅"
					}
					items = append(items, SystemPromptMenuItem{
						title:       itemTitle,
						description: prompt.Content,
					})
				}
				s.list.SetItems(items)
			}
			return s, func() tea.Msg {
				return SystemPromptChangeMsg{NewPrompt: selected.description}
			}
		}

		if s.focused == "list" {
			var listCmd tea.Cmd
			s.list, listCmd = s.list.Update(msg)
					cmds = append(cmds, listCmd)

			if selected, ok := s.list.SelectedItem().(SystemPromptMenuItem); ok {
				s.viewport.SetContent(selected.description)
			}
		} else {
			var vpCmd tea.Cmd
			s.viewport, vpCmd = s.viewport.Update(msg)
			cmds = append(cmds, vpCmd)
		}
	}

	return s, tea.Batch(cmds...)
}

func (s *SwitchSystemPromptsView) View() string {
	listStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.CurrentTheme.Primary.GetColor()).
		Padding(1, 1).
		Width(34).
		Height(s.height)

	if s.focused == "list" {
		listStyle = listStyle.BorderForeground(theme.CurrentTheme.Border.Active.GetColor())
	}

	vpStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.CurrentTheme.Primary.GetColor()).
		Padding(1, 1).
		Width(s.width - 38).
		Height(s.height)

	if s.focused == "viewport" {
		vpStyle = vpStyle.BorderForeground(theme.CurrentTheme.Border.Active.GetColor())
	}

	listView := listStyle.Render(s.list.View())
	viewportView := vpStyle.Render(s.viewport.View())

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		listView,
		viewportView,
	)
}

func (s *SwitchSystemPromptsView) SetSize(width, height int) {
	s.width = width
	s.height = height
	listHeight := height - 4
	s.list.SetSize(30, listHeight)
	viewportWidth := width - 38
	s.viewport.Width = viewportWidth - 4
	s.viewport.Height = height - 4
} 