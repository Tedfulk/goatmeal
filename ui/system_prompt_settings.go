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
)

type SystemPromptMenuItem struct {
	title string
}

func (i SystemPromptMenuItem) Title() string       { return i.title }
func (i SystemPromptMenuItem) Description() string { return "" }
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
	// Create list items for actions
	items := []list.Item{
		SystemPromptMenuItem{title: "Add System Prompt"},
		SystemPromptMenuItem{title: "Delete System Prompt"},
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
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Padding(0, 0, 1, 2)

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
		s.deleteView, deleteCmd = s.deleteView.Update(msg)
		return s, deleteCmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
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

	menuStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 1).
		Width(36).
		Height(s.height - 18)

	titleStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Width(26).
		Align(lipgloss.Center)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render("  System Prompt Settings"),
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
	s.list.SetSize(40, height-18)
	
	if s.deleteView != nil {
		s.deleteView.SetSize(width, height)
	}
}

func (s *SystemPromptSettings) openSystemPromptsInEditor() {
	// Get the default editor from environment variables
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

	// Open the config file in the editor
	cmd := exec.Command(editor, configPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return
	}
} 