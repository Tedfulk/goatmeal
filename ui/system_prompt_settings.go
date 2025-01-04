package ui

import (
	"os"
	"os/exec"

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
}

func NewSystemPromptSettings(cfg *config.Config) SystemPromptSettings {
	items := []list.Item{
		SystemPromptMenuItem{title: "Add System Prompt", description: "Add a new system prompt"},
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
			case "Add System Prompt":
					go s.openSystemPromptsInEditor()
			case "Delete System Prompt":
				s.showDeleteView = true
				s.deleteView = NewDeleteSystemPromptsView(s.config)
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
}

func (s *SystemPromptSettings) openSystemPromptsInEditor() {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "system-prompt-*.txt")
	if err != nil {
		return
	}
	defer os.Remove(tmpfile.Name())

	// Write template or existing content
	template := `# System Prompt Template
# Enter your system prompt below. The first line starting with "Title:" will be used as the prompt title.
# Lines starting with # will be ignored.
# Example:
# Title: My Custom Assistant
# You are a helpful assistant...

Title: 
`
	if _, err := tmpfile.Write([]byte(template)); err != nil {
		return
	}
	tmpfile.Close()

	// Open editor
	cmd := exec.Command(editor, tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return
	}

	// TODO: Implement reading and processing the content
	// This part would need to be implemented based on your config structure
} 