package ui

import (
	"fmt"
	"goatmeal/config"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SystemPromptAction represents different actions we can take
type SystemPromptAction int

const (
	AddPrompt SystemPromptAction = iota
	EditPrompt
	SelectPrompt
	Back
)

// Message type for system prompt actions
type SystemPromptMsg struct {
	action SystemPromptAction
	prompt string
}

type editorFinishedMsg struct{ err error }

type SystemPromptModel struct {
	items    []SystemPromptItem
	selected int
	keys     menuKeyMap
	width    int
	height   int
	colors   config.ThemeColors
	err      error
}

type SystemPromptItem struct {
	title       string
	description string
	action      SystemPromptAction
}

// Add new message type for prompt selection
type SystemPromptSelectedMsg struct {
	prompt string
}

func NewSystemPromptMenu(config *config.Config) SystemPromptModel {
	return SystemPromptModel{
		items: []SystemPromptItem{
			{title: "Add New Prompt", action: AddPrompt},
			{title: "Edit Prompts", action: EditPrompt},
			{title: "Select Prompt", action: SelectPrompt},
			{title: "Back", action: Back},
		},
		selected: 0,
		keys:     menuKeys,
		colors:   config.GetThemeColors(),
		err:      nil,
	}
}

func (m SystemPromptModel) Init() tea.Cmd {
	return nil
}

func openConfigInEditor() tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nvim"
	}

	// Get config file path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return func() tea.Msg {
			return editorFinishedMsg{err: fmt.Errorf("error getting home directory: %w", err)}
		}
	}
	configPath := filepath.Join(homeDir, ".goatmeal", "config.yaml")

	// Create command to open config file
	c := exec.Command(editor, configPath)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

func (m SystemPromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case editorFinishedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		return m, func() tea.Msg { return ChangeViewMsg(settingsView) }

	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}

		if msg.String() == "esc" {
			return m, func() tea.Msg { return ChangeViewMsg(settingsView) }
		}

		switch {
		case key.Matches(msg, m.keys.Up):
			m.selected--
			if m.selected < 0 {
				m.selected = len(m.items) - 1
			}

		case key.Matches(msg, m.keys.Down):
			m.selected++
			if m.selected >= len(m.items) {
				m.selected = 0
			}

		case key.Matches(msg, m.keys.Select):
			selectedItem := m.items[m.selected]
			switch selectedItem.action {
			case AddPrompt:
				return m, openConfigInEditor()
			case EditPrompt:
				return m, openConfigInEditor()
			case SelectPrompt:
				// Show prompt selection view
				return m, func() tea.Msg {
					return ChangeViewMsg(promptSelectionView)
				}
			}
		}
	}

	return m, nil
}

func (m SystemPromptModel) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(m.colors.MenuTitle)).
		Padding(1, 0).
		Align(lipgloss.Center)

	menuStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.colors.MenuBorder)).
		Padding(2, 4).
		Width(80).
		Align(lipgloss.Center)

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(m.colors.MenuSelected))

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.colors.MenuNormal))

	// Build menu items
	var menuItems string
	for i, item := range m.items {
		menuItem := item.title
		if i == m.selected {
			menuItem = "▸ " + menuItem
			menuItem = selectedStyle.Render(menuItem)
		} else {
			menuItem = "  " + menuItem
			menuItem = normalStyle.Render(menuItem)
		}
		
		if item.description != "" {
			menuItem += " " + normalStyle.Render(item.description)
		}
		menuItems += menuItem + "\n"
	}

	// Create the menu box
	menu := menuStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			titleStyle.Render("System Prompt Management"),
			"",
			menuItems,
			"",
			normalStyle.Render("↑/↓: navigate • enter: select • esc: back • q: quit"),
		),
	)

	return menu
}

// Helper function to truncate long strings
func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
} 