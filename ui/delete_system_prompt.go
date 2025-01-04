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

// DeleteSystemPromptsKeyMap defines the key bindings for the delete system prompts view
type DeleteSystemPromptsKeyMap struct {
	Back        key.Binding
	SwitchFocus key.Binding
}

var DefaultDeleteSystemPromptsKeyMap = DeleteSystemPromptsKeyMap{
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	SwitchFocus: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch focus"),
	),
}

type DeleteSystemPromptsView struct {
	list     list.Model
	viewport viewport.Model
	config   *config.Config
	width    int
	height   int
	selected int
	focused  string // "list" or "viewport"
	keys     DeleteSystemPromptsKeyMap
}

func NewDeleteSystemPromptsView(cfg *config.Config) *DeleteSystemPromptsView {
	// Create list items from system prompts
	var items []list.Item
	for _, prompt := range cfg.SystemPrompts {
		items = append(items, SystemPromptMenuItem{
			title:       prompt.Title,
			description: prompt.Content,
		})
	}

	// Setup list
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetHeight(1)
	
	l := list.New(items, delegate, 30, 30)
	l.Title = "Delete System Prompts"
	l.SetShowHelp(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = theme.BaseStyle.Title.
		Foreground(theme.CurrentTheme.Primary.GetColor()).
		Align(lipgloss.Center).
		Width(30)

	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			DefaultDeleteSystemPromptsKeyMap.Back,
			DefaultDeleteSystemPromptsKeyMap.SwitchFocus,
		}
	}
	l.AdditionalFullHelpKeys = l.AdditionalShortHelpKeys

	// Initialize viewport
	vp := viewport.New(102, 32)

	if len(items) > 0 {
		vp.SetContent(items[0].(SystemPromptMenuItem).description)
	}

	return &DeleteSystemPromptsView{
		list:     l,
		viewport: vp,
		config:   cfg,
		focused:  "list",
		keys:     DefaultDeleteSystemPromptsKeyMap,
	}
}

func (d *DeleteSystemPromptsView) Update(msg tea.Msg) (*DeleteSystemPromptsView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle tab key for focus switching
		if key.Matches(msg, d.keys.SwitchFocus) {
			if d.focused == "list" {
				d.focused = "viewport"
			} else {
				d.focused = "list"
			}
			return d, nil
		}

		// Handle back key
		if key.Matches(msg, d.keys.Back) {
			return nil, nil
		}

		// Handle enter key for deleting prompts
		if msg.String() == "enter" && d.focused == "list" {
			selected := d.list.SelectedItem().(SystemPromptMenuItem)
			// Find and remove the selected prompt
			for i, prompt := range d.config.SystemPrompts {
				if prompt.Title == selected.title {
					d.config.SystemPrompts = append(d.config.SystemPrompts[:i], d.config.SystemPrompts[i+1:]...)
					break
				}
			}

			// Update the config file
			manager, err := config.NewManager()
			if err == nil {
				manager.SetSystemPrompts(d.config.SystemPrompts)
			}

			// Update the list items
			var items []list.Item
			for _, prompt := range d.config.SystemPrompts {
				items = append(items, SystemPromptMenuItem{
					title:       prompt.Title,
					description: prompt.Content,
				})
			}
			d.list.SetItems(items)
		}

		// Only pass key events to the focused component
		if d.focused == "list" {
			var listCmd tea.Cmd
			d.list, listCmd = d.list.Update(msg)
			cmds = append(cmds, listCmd)

			// Update viewport content when list selection changes
			if selected, ok := d.list.SelectedItem().(SystemPromptMenuItem); ok {
				d.viewport.SetContent(selected.description)
			}
		} else {
			var vpCmd tea.Cmd
			d.viewport, vpCmd = d.viewport.Update(msg)
			cmds = append(cmds, vpCmd)
		}
	}

	return d, tea.Batch(cmds...)
}

func (d *DeleteSystemPromptsView) View() string {
	// Left container (list)
	listStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.CurrentTheme.Primary.GetColor()).
		Padding(1, 1).
		Width(34).
		Height(d.height)

	// Highlight the focused container's border
	if d.focused == "list" {
		listStyle = listStyle.BorderForeground(theme.CurrentTheme.Border.Active.GetColor())
	}

	// Style the viewport based on focus
	vpStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.CurrentTheme.Primary.GetColor()).
		Padding(1, 1).
		Width(d.width - 38).  // Account for list width + gap
		Height(d.height)

	if d.focused == "viewport" {
		vpStyle = vpStyle.BorderForeground(theme.CurrentTheme.Border.Active.GetColor())
	}

	// Apply the styles
	listView := listStyle.Render(d.list.View())
	viewportView := vpStyle.Render(d.viewport.View())

	// Join the containers horizontally with a gap
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		listView,
		viewportView,
	)
}

func (d *DeleteSystemPromptsView) SetSize(width, height int) {
	d.width = width
	d.height = height

	// Set list size (account for borders and padding)
	listHeight := height - 4 // Account for borders and padding
	d.list.SetSize(30, listHeight)

	// Set viewport size (account for borders and padding)
	viewportWidth := width - 38 // Account for list width (34) and gap (4)
	d.viewport.Width = viewportWidth - 4 // Account for borders and padding
	d.viewport.Height = height - 4 // Account for borders and padding
} 
