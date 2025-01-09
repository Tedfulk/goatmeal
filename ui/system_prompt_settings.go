package ui

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tedfulk/goatmeal/config"
	"github.com/tedfulk/goatmeal/ui/theme"
)

type SystemPromptMenuItem struct {
	title       string
	description string
}

func (i SystemPromptMenuItem) Title() string       { return i.title }
func (i SystemPromptMenuItem) Description() string { return i.description }
func (i SystemPromptMenuItem) FilterValue() string { return i.title }

type SystemPromptSettings struct {
	list           list.Model
	config         *config.Config
	width          int
	height         int
	showDeleteView bool
	deleteView     *DeleteSystemPromptsView
	showSwitchView bool
	switchView     *SwitchSystemPromptsView
	showAddView    bool
	addView        *AddSystemPromptsView
}

func NewSystemPromptSettings(cfg *config.Config) SystemPromptSettings {
	items := []list.Item{
		SystemPromptMenuItem{title: "Open config file in editor", description: "Edit the config file directly"},
		SystemPromptMenuItem{title: "Add system prompt", description: "Add a new system prompt"},
		SystemPromptMenuItem{title: "Switch system prompt", description: "Switch to a different system prompt"},
		SystemPromptMenuItem{title: "Delete System Prompt", description: "Delete an existing system prompt"},
	}

	// Setup list
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

	return SystemPromptSettings{
		list:           l,
		config:         cfg,
		showDeleteView: false,
		deleteView:     NewDeleteSystemPromptsView(cfg),
		showSwitchView: false,
		switchView:     NewSwitchSystemPromptsView(cfg),
		showAddView:    false,
		addView:        NewAddSystemPromptsView(cfg),
	}
}

func (s SystemPromptSettings) Update(msg tea.Msg) (SystemPromptSettings, tea.Cmd) {
	var cmd tea.Cmd

	if s.showDeleteView {
		var deleteCmd tea.Cmd
		newDeleteView, deleteCmd := s.deleteView.Update(msg)
		if newDeleteView == nil {
			s.showDeleteView = false
			return s, nil
		}
		s.deleteView = newDeleteView
		return s, deleteCmd
	}

	if s.showSwitchView {
		var switchCmd tea.Cmd
		newSwitchView, switchCmd := s.switchView.Update(msg)
		
		if newSwitchView == nil {
			s.showSwitchView = false
			if switchCmd != nil {
				if cmd := switchCmd(); cmd != nil {
					if _, ok := cmd.(SystemPromptChangeMsg); ok {
						return s, tea.Batch(
							switchCmd,
							func() tea.Msg {
								return SetViewMsg{view: "settings"}
							},
						)
					}
				}
			}
			return s, nil
		}
		
		s.switchView = newSwitchView
		return s, switchCmd
	}

	if s.showAddView {
		var addCmd tea.Cmd
		newAddView, addCmd := s.addView.Update(msg)
		if newAddView == nil {
			s.showAddView = false
			return s, nil
		}
		s.addView = newAddView
		return s, addCmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return s, func() tea.Msg {
				return SetViewMsg{view: "settings"}
			}
		case "enter":
			selected := s.list.SelectedItem().(SystemPromptMenuItem)
			switch selected.title {
			case "Open config file in editor":
				go s.openSystemPromptsInEditor()
			case "Delete System Prompt":
				s.showDeleteView = true
				s.deleteView = NewDeleteSystemPromptsView(s.config)
			case "Switch system prompt":
				s.showSwitchView = true
				s.switchView = NewSwitchSystemPromptsView(s.config)
			case "Add system prompt":
				s.showAddView = true
				s.addView = NewAddSystemPromptsView(s.config)
			}
		}
	}

	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (s SystemPromptSettings) View() string {
	if s.showDeleteView {
		return s.deleteView.View()
	}

	if s.showSwitchView {
		return s.switchView.View()
	}

	if s.showAddView {
		return s.addView.View()
	}

	menuStyle := theme.BaseStyle.Menu.
		BorderForeground(theme.CurrentTheme.Primary.GetColor())

	titleStyle := theme.BaseStyle.Title.
		Foreground(theme.CurrentTheme.Primary.GetColor())

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("System Prompt Settings"),
		s.list.View(),
	)

	return lipgloss.Place(
		s.width,
		s.height,
		lipgloss.Center,
		lipgloss.Center,
		menuStyle.Render(content),
	)
}

func (s *SystemPromptSettings) SetSize(width, height int) {
	s.width = width
	s.height = height
	s.list.SetSize(width-4, height-12)
	if s.deleteView != nil {
		s.deleteView.SetSize(width, height)
	}
	if s.switchView != nil {
		s.switchView.SetSize(width, height)
	}
	if s.addView != nil {
		s.addView.SetSize(width, height)
	}
}

func (s *SystemPromptSettings) openSystemPromptsInEditor() {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		// Try common editors in order of preference
		if _, err := exec.LookPath("nvim"); err == nil {
			editor = "nvim"
		} else if _, err := exec.LookPath("nano"); err == nil {
			editor = "nano"
		} else {
			editor = "vim"
		}
	}

	// Get the config file path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	configPath := filepath.Join(homeDir, ".config", "goatmeal", "config.yaml")

	// Open editor
	cmd := exec.Command(editor, configPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return
	}

	// Reload the config after editing
	if newConfig, err := config.Load(); err == nil {
		s.config = newConfig
	}
} 